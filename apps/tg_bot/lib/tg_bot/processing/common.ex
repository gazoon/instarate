defmodule TGBot.Processing.Common do

  require Logger
  alias Utils.Messages.Task, as: TaskMessage
  alias Utils.Messages.User, as: MessageUser
  alias TGBot.{MatchPhotoCache, Localization}
  alias Voting.Girl
  alias TGBot.Chats.Chat
  alias Voting.InstagramProfiles.Model, as: InstagramProfile
  import Localization, only: [get_translation: 3, get_translation: 2]
  use Utils.Meter
  @session_duration 1_200_000 # 20 minutes in milliseconds

  @next_pair_task :send_next_pair
  @daily_activation_task :daily_activation
  def next_pair_task, do: @next_pair_task
  def daily_activation_task, do: @daily_activation_task

  @config Application.get_env(:tg_bot, __MODULE__)
  @messenger @config[:messenger]
  @scheduler @config[:scheduler]
  @photos_cache @config[:photos_cache]
  @pictures @config[:pictures_concatenator]

  @spec build_voter_id(MessageUser.t) :: String.t
  def build_voter_id(user) do
    "tg_user:" <> Integer.to_string(user.id)
  end

  @spec build_voters_group_id(integer) :: String.t
  def build_voters_group_id(chat_id) do
    "tg_chat:" <> Integer.to_string(chat_id)
  end

  @spec send_next_girls_pair(Chat.t, Keyword.t) :: Chat.t
  def send_next_girls_pair(chat, opts \\ []) do
    voters_group_id = build_voters_group_id(chat.chat_id)
    case Voting.get_next_pair(chat.competition, voters_group_id) do
      {girl_one, girl_two} -> send_girls_pair(chat, girl_one, girl_two, opts)
      :error ->
        @messenger.send_text(
          chat.chat_id,
          get_translation(chat, "no_more_girls_in_competition"),
          static_keyboard: :remove
        )
        chat
    end
  end

  @spec display_girl_info(Chat.t, Girl.t) :: any
  def display_girl_info(chat, girl) do
    profile_url = Girl.get_profile_url(girl)
    @messenger.send_markdown(chat.chat_id, "[#{girl.username}](#{profile_url})")
    send_single_girl_photo(chat.chat_id, girl)
    @messenger.send_text(
      chat.chat_id,
      get_translation(
        chat,
        "girl_statistics",
        position: Girl.get_position(girl),
        wins: girl.wins,
        loses: girl.loses
      )
    )
  end

  @spec try_to_send_next_pair(Chat.t) :: Chat.t
  def try_to_send_next_pair(chat) do
    if chat.last_match do
      time_to_show = 1000 * chat.voting_timeout + chat.last_match.shown_at
      if time_to_show > Utils.timestamp_milliseconds() do
        task_args = %{last_match_message_id: chat.last_match.message_id}
        task = TaskMessage.new(
          chat.chat_id,
          time_to_show,
          @next_pair_task,
          args: task_args
        )
        case @scheduler.create_task(task) do
          {:ok, _} -> Logger.info("Schedule send next pair task")
          _ -> nil
        end
        chat
      else
        send_next_girls_pair(chat)
      end
    else
      send_next_girls_pair(chat)
    end
  end


  @spec add_single_girl(Chat.t, String.t) :: any
  def add_single_girl(chat, photo_link) do
    case Voting.add_girl(photo_link) do
      {:ok, girl} ->
        profile_url = InstagramProfile.get_profile_url(girl)
        text = get_translation(
          chat,
          "girl_added",
          username: girl.username,
          profile_url: profile_url
        )
        @messenger.send_markdown(chat.chat_id, text)

      {:error, error_msg} -> @messenger.send_text(chat.chat_id, error_msg)
    end
  end


  @spec send_girl_from_top(Chat.t, integer) :: Chat.t
  def send_girl_from_top(chat, girl_offset) do
    case Voting.get_top(chat.competition, 2, offset: girl_offset) do
      [current_girl | next_girls] ->
        keyboard = if length(next_girls) != 0, do: [["Next girl"]], else: :remove
        #        keyboard = if length(next_girls) != 0 do
        #          next_girl_offset = Integer.to_string(girl_offset + 1)
        #          [[%{text: "Next", payload: Callback.build_payload(@get_top_callback, next_girl_offset)}]]
        #        else
        #          nil
        #        end
        caption = get_translation(chat, "place_in_competition", place: girl_offset + 1) <>
                  Girl.get_profile_url(current_girl)
        send_single_girl_photo(
          chat.chat_id,
          current_girl,
          caption: caption,
          #          keyboard: keyboard,
          static_keyboard: keyboard,
        )
        %Chat{chat | current_top_offset: girl_offset}
      [] ->
        girls_number = Voting.get_girls_number(chat.competition)
        Logger.warn(
          "Girl offset #{girl_offset} more than number of girls in the competition: #{girls_number}"
        )
        @messenger.send_text(
          chat.chat_id,
          get_translation(chat, "no_more_girls_in_top"),
          static_keyboard: :remove
        )
        chat
    end
  end

  @spec send_girls_pair(Chat.t, Girl.t, Girl.t, Keyword.t) :: Chat.t
  defp send_girls_pair(chat, girl_one, girl_two, opts) do
    tg_file_id = MatchPhotoCache.get(girl_one.photo, girl_two.photo)
    {left_girl, right_girl, tg_file_id} = if tg_file_id do
      {girl_one, girl_two, tg_file_id}
    else
      tg_file_id = MatchPhotoCache.get(girl_two.photo, girl_one.photo)
      {girl_two, girl_one, tg_file_id}
    end
    left_girl_url = Girl.get_profile_url(left_girl)
    right_girl_url = Girl.get_profile_url(right_girl)
    keyboard = [["Left", "Right"]]
    #    keyboard = [
    #      [
    #        %{
    #          text: "Left",
    #          payload: Callback.build_payload(
    #            @vote_callback,
    #            left_girl.username <> @usernames_separator <> right_girl.username
    #          )
    #        },
    #        %{
    #          text: "Right",
    #          payload: Callback.build_payload(
    #            @vote_callback,
    #            right_girl.username <> @usernames_separator <> left_girl.username
    #          )
    #        },
    #      ]
    #    ]
    caption_text = "#{left_girl_url} vs #{right_girl_url}"
    message_before = Keyword.get(opts, :message_before)
    if message_before, do: @messenger.send_text(chat.chat_id, message_before)

    message_id = if tg_file_id do
      Logger.info("Use cached match photo #{tg_file_id}")
      {message_id, _} = @messenger.send_photo(
        chat.chat_id,
        tg_file_id,
        caption: caption_text,
        static_keyboard: keyboard,
        one_time_keyboard: true
      )
      message_id
    else
      left_girl_photo_url = Girl.get_photo_url(left_girl)
      right_girl_photo_url = Girl.get_photo_url(right_girl)
      Logger.info("Concatenate #{left_girl_photo_url} and #{right_girl_photo_url}")
      match_photo = measure metric_name: "concatenate_photos" do
        @pictures.concatenate(left_girl_photo_url, right_girl_photo_url)
      end
      {message_id, tg_file_id} = @messenger.send_photo(
          chat.chat_id,
          match_photo,
          caption: caption_text,
          static_keyboard: keyboard,
          one_time_keyboard: true,
        binary_data: true
      )
      MatchPhotoCache.set(left_girl.photo, right_girl.photo, tg_file_id)
      message_id
    end

    schedule_daily_activation(chat)
    current_match = Chat.Match.new(message_id, left_girl.username, right_girl.username)
    %Chat{chat | last_match: current_match}
  end

  @spec send_single_girl_photo(integer, Girl.t, Keyword.t) :: any
  defp send_single_girl_photo(chat_id, girl, opts \\ []) do
    tg_file_id = @photos_cache.get(girl.photo)
    {photo, is_new} = if tg_file_id, do: {tg_file_id, false}, else: {Girl.get_photo_url(girl), true}
    {_, tg_file_id} = @messenger.send_photo(chat_id, photo, opts)
    if is_new do
      Logger.info("Save new tg file id for photo #{girl.photo}")
      @photos_cache.set(girl.photo, tg_file_id)
    end
  end

  @spec schedule_daily_activation(Chat.t) :: any
  defp schedule_daily_activation(chat) do
    if chat.self_activation_allowed do
      time_to_activate = Utils.timestamp_milliseconds() + 24 * 60 * 60 * 1000 - @session_duration
      task = TaskMessage.new(chat.chat_id, time_to_activate, @daily_activation_task)
      @scheduler.create_or_replace_task(task)
      Logger.info("Schedule next day activation")
    else
      Logger.info("Self activation disabled for the chat")
    end
  end
end
