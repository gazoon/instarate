defmodule TGBot do

  require Logger
  alias TGBot.Messages.Text, as: TextMessage
  alias TGBot.Messages.Callback, as: Callback
  alias TGBot.Messages.Task, as: TaskMessage
  alias TGBot.Messages.User, as: MessageUser
  alias TGBot.{Message, Pictures}
  alias TGBot.Chats.Chat
  alias Voting.Girl

  @start_cmd "start"
  @add_girl_cmd "addGirl"
  @get_top_cmd "showTop"
  @next_top_cmd "Next girl"
  @get_girl_info_cmd "girlInfo"
  @help_cmd "help"
  @left_vote_cmd "left"
  @right_vote_cmd "right"
  @global_competition_cmd "globalCompetition"
  @celebrities_competition_cmd "celebritiesCompetition"
  @normal_competition_cmd "normalCompetition"
  @enable_daily_activation_cmd "enableActivation"
  @disable_daily_activation_cmd "disableActivation"
  @set_voting_timeout_cmd "votingTimeout"
  @delete_girls_cmd "deleteGirls"

  @min_voting_timeout 5

  @vote_callback "vt"
  @get_top_callback "top"

  @session_duration 1_200_000 # 20 minutes in milliseconds

  @session_duration_seconds div(@session_duration, 1000)
  @usernames_separator "|"

  @next_pair_task :send_next_pair
  @daily_activation_task :daily_activation

  @config Application.get_env(:tg_bot, __MODULE__)
  @chats_storage @config[:chats_storage]
  @messenger @config[:messenger]
  @scheduler @config[:scheduler]
  @admins @config[:admins]

  @spec on_message(map()) :: any
  def on_message(message_container) do
    message_type = String.to_atom(message_container["type"])
    message_data = message_container["data"]
    message_info = case message_type do
      :text -> {TextMessage, &on_text_message/2}
      :callback -> {Callback, &on_callback/2}
      :task -> {TaskMessage, &on_task/2}
      _ -> nil
    end
    case message_info do
      {message_cls, handler_func} ->
        message = message_cls.from_data(message_data)
        process_message(message, handler_func)
        Logger.info("Finish message processing")
      nil -> Logger.error("Unsupported message type: #{message_type}")
    end
  end

  @spec initialize_context(Message.t) :: any
  defp initialize_context(message) do
    request_id = UUID.uuid4()
    Logger.metadata([request_id: request_id, chat_id: message.chat_id])
  end

  @spec process_message(Message.t, ((Message.t, Chat.t) -> Chat.t)) :: any
  defp process_message(message, handler) do
    initialize_context(message)
    {chat, is_new} = get_chat(message)
    chat_after_processing = handler.(message, chat)
    if chat_after_processing != chat || is_new do
      Logger.info("Save updated chat info")
      @chats_storage.save(chat_after_processing)
    end
  end

  @spec get_chat(Message.t) :: {Chat.t, boolean}
  defp get_chat(message) do
    chat_id = Message.chat_id(message)
    case @chats_storage.get(chat_id) do
      nil ->
        members_count = @messenger.get_chat_members_number(chat_id) - 1
        {Chat.new(chat_id, members_count), true}
      chat -> {chat, false}
    end
  end

  @spec on_task(TaskMessage.t, Chat.t) :: Chat.t
  defp on_task(task, chat)  do
    handler = case task.name do
      @next_pair_task -> &handle_next_pair_task/2
      @daily_activation_task -> &handle_daily_activation_task/2
      _ -> nil
    end
    if handler do
      Logger.info("Handle #{task.name} task")
      handler.(task, chat)
    else
      Logger.info("Received unknown #{task.name} task")
      chat
    end
  end

  @spec handle_next_pair_task(TaskMessage.t, Chat.t) :: Chat.t
  defp handle_next_pair_task(task, chat) do
    task_match_message_id = task.args.last_match_message_id
    actual_match_message_id = chat.last_match.message_id
    if task_match_message_id == actual_match_message_id do
      send_next_girls_pair(chat)
    else
      Logger.info(
        "Skip next pair task, not invalid match message id: #{task_match_message_id}
         actual one: #{actual_match_message_id}"
      )
      chat
    end
  end

  @spec handle_daily_activation_task(TaskMessage.t, Chat.t) :: Chat.t
  defp handle_daily_activation_task(_task, chat) do
    send_next_girls_pair(chat, message_before: "Hi! don't you mind to vote?")
  end

  @spec on_text_message(TextMessage.t, Chat.t) :: Chat.t
  defp on_text_message(message, chat) do
    if TextMessage.appeal_to_bot?(message) || TextMessage.reply_to_bot?(message)
       || !message.is_group_chat do
      process_text_message(message, chat)
    else
      Logger.info("Skip message #{inspect message} it's not an appeal, reply to the bot or private")
      chat
    end
  end

  @spec process_text_message(TextMessage.t, Chat.t) :: Chat.t
  defp process_text_message(message, chat) do
    Logger.info("Process text message #{inspect message}")
    commands = [
      {@start_cmd, &handle_start_cmd/2},
      {@add_girl_cmd, &handle_add_girl_cmd/2},
      {@get_top_cmd, &handle_get_top_cmd/2},
      {@next_top_cmd, &handle_next_top_cmd/2},
      {@get_girl_info_cmd, &handle_get_girl_info_cmd/2},
      {@help_cmd, &handle_help_cmd/2},
      {@left_vote_cmd, &handle_left_vote_cmd/2},
      {@right_vote_cmd, &handle_right_vote_cmd/2},
      {@global_competition_cmd, &handle_global_competition_cmd/2},
      {@celebrities_competition_cmd, &handle_celebrities_competition_cmd/2},
      {@normal_competition_cmd, &handle_normal_competition_cmd/2},
      {@enable_daily_activation_cmd, &handle_enable_activation_cmd/2},
      {@disable_daily_activation_cmd, &handle_disable_activation_cmd/2},
      {@set_voting_timeout_cmd, &handle_set_voting_timeout_cmd/2},
      {@delete_girls_cmd, &handle_delete_girls_cmd/2},
    ]
    message_text = message.text_lowercase
    command = commands
              |> Enum.find(
                   fn ({cmd_name, _}) -> String.contains?(message_text, String.downcase(cmd_name))
                   end
                 )
    case command do
      {command_name, handler} -> Logger.info("Handle #{command_name} command")
                                 handler.(message, chat)
      nil ->
        Logger.info("Message #{message_text} doesn't contain commands, handle as regular message")
        handle_regular_message(message, chat)
    end
  end

  @spec handle_regular_message(TextMessage.t, Chat.t) :: Chat.t
  defp handle_regular_message(message, chat) do
    photo_links = String.split(message.text, "\n")
    functions = for photo_link <- photo_links do
      fn -> Voting.add_girl(photo_link) end
    end
    Utils.parallelize_tasks(functions)
    @messenger.send_text(chat.chat_id, "All girls were processed")
    chat
  end

  @spec build_voter_id(MessageUser.t) :: String.t
  defp build_voter_id(user) do
    "tg_user:" <> Integer.to_string(user.id)
  end

  @spec build_voters_group_id(integer) :: String.t
  defp build_voters_group_id(chat_id) do
    "tg_chat:" <> Integer.to_string(chat_id)
  end

  @spec handle_start_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_start_cmd(_message, chat) do
    try_to_send_pair(chat)
  end

  @spec send_next_girls_pair(Chat.t, Keyword.t) :: Chat.t
  defp send_next_girls_pair(chat, opts \\ []) do
    voters_group_id = build_voters_group_id(chat.chat_id)
    {girl_one, girl_two} = Voting.get_next_pair(chat.competition, voters_group_id)
    girl_one_url = Girl.get_profile_url(girl_one)
    girl_two_url = Girl.get_profile_url(girl_two)

    match_photo = Pictures.concatenate(girl_one.photo, girl_two.photo)

    #    keyboard = [
    #      [
    #        %{
    #          text: "Left",
    #          payload: Callback.build_payload(
    #            @vote_callback,
    #            girl_one.username <> @usernames_separator <> girl_two.username
    #          )
    #        },
    #        %{
    #          text: "Right",
    #          payload: Callback.build_payload(
    #            @vote_callback,
    #            girl_two.username <> @usernames_separator <> girl_one.username
    #          )
    #        },
    #      ]
    #    ]
    keyboard = [["Left", "Right"]]
    caption_text = "#{girl_one_url} vs #{girl_two_url}"
    message_before = Keyword.get(opts, :message_before)
    if message_before, do: @messenger.send_text(chat.chat_id, message_before)
    message_id = try do
      @messenger.send_photo(
        chat.chat_id,
        match_photo,
        caption: caption_text,
        static_keyboard: keyboard,
        one_time_keyboard: true,
      )
    after
      File.rm!(match_photo)
    end
    schedule_daily_activation(chat)
    current_match = Chat.Match.new(message_id, girl_one.username, girl_two.username)
    %Chat{chat | last_match: current_match}
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

  @spec handle_delete_girls_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_delete_girls_cmd(message, chat) do
    girl_uris = TextMessage.get_command_args(message)
    if Enum.member?(@admins, message.user.id) do
      Voting.delete_girls(girl_uris)
      @messenger.send_text(message.chat_id, "Girls were deleted")
    else
      Logger.warn("Non-admin user #{message.user.id} tried to delete girls")
    end
    chat
  end

  @spec handle_add_girl_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_add_girl_cmd(message, chat) do
    photo_link = TextMessage.get_command_arg(message)
    if photo_link do
      case Voting.add_girl(photo_link) do
        {:ok, girl} ->
          profile_url = Girl.get_profile_url(girl)
          text = "Girl [#{girl.username}](#{profile_url}) has been successfully added!"
          @messenger.send_markdown(message.chat_id, text)

        {:error, error_msg} -> @messenger.send_text(message.chat_id, error_msg)
      end
    else
      @messenger.send_text(
        message.chat_id,
        "Please send me a link\naddGirl@InstaRateBot <photo_link>",
        disable_web_page_preview: true
      )
    end
    chat
  end

  @spec handle_get_top_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_get_top_cmd(message, chat) do
    optional_start_position = TextMessage.get_command_arg(message)
    offset = case Integer.parse(optional_start_position) do
      {start_position, ""} when start_position > 0 -> start_position - 1
      _ -> 0
    end
    send_girl_from_top(chat, offset)
  end

  @spec handle_enable_activation_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_enable_activation_cmd(_message, chat) do
    @messenger.send_text(chat.chat_id, "Daily notifications enabled")
    %Chat{chat | self_activation_allowed: true}
  end

  @spec handle_disable_activation_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_disable_activation_cmd(_message, chat) do
    @messenger.send_text(chat.chat_id, "Daily notifications disabled")
    @scheduler.delete_task(chat.chat_id, @daily_activation_task)
    %Chat{chat | self_activation_allowed: false}
  end

  @spec handle_next_top_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_next_top_cmd(_message, chat) do
    next_offset = chat.current_top_offset + 1
    send_girl_from_top(chat, next_offset)
  end

  @spec handle_set_voting_timeout_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_set_voting_timeout_cmd(message, chat) do
    arg = TextMessage.get_command_arg(message)
    case Integer.parse(arg) do
      {timeout, ""} when @min_voting_timeout <= timeout and timeout < @session_duration_seconds ->
        @messenger.send_text(chat.chat_id, "Now voting timeout is #{timeout} seconds")
        %Chat{chat | voting_timeout: timeout}
      {_, ""} ->
        @messenger.send_text(
          chat.chat_id,
          "Timeout must be more than #{@min_voting_timeout - 1} seconds and less than #{
            div @session_duration_seconds, 60
          } minutes"
        )
        chat
      _ ->
        @messenger.send_text(chat.chat_id, "Please, enter valid number")
        chat
    end
  end

  @spec handle_help_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_help_cmd(message, chat) do
    text = """
    Hi there! My mission is to find the most attractive girls on Instagram!
    Just select which of two girls looks better and vote by pressing a button below.
    I support following commands:

    /start@InstaRateBot - Get the next girls pair to compare.

    /showTop@InstaRateBot - Show the top girls in the competition. You can also pass a position to start showing from, i.e pass 10 to start from the 10th girl in the competition and skip the first 9.

    /celebritiesCompetition@InstaRateBot - You will vote and see only girls with 500k+ followers.

    /normalCompetition@InstaRateBot - You will vote and see girls who have less than 500k followers

    /globalCompetition@InstaRateBot - You will vote for all girls, it's a default option.

    /enableActivation@InstaRateBot - Let me daily send you a new girls pair, just in case you forgot about me.

    /disableActivation@InstaRateBot - Disable daily new match sending. By default it's enabled.

    /help@InstaRateBot - Show this message.

    And also another type of commands, with an additional input:

    addGirl@InstaRateBot <link to one of her photos on instagram> - Add girl to the competition, You can add any girl, just paste a link to her instagram photo. The girl must have a public account. For example:
    addGirl@InstaRateBot <photo_link>

    girlInfo@InstaRateBot <girl username or link to the profile> - Get info about particular girl statistic in the competition. For example:
    girlInfo@InstaRateBot <username>
    girlInfo@InstaRateBot <profile_link>

    votingTimeout@InstaRateBot <timeout> - Set the time, I will wait before sendind a new girls pair. In seconds, minimum - 5 seconds.
    votingTimeout@InstaRateBot 10

    Some notes:
    In group chats you have to mention me (with the '@' sign) in the message, if you want to command me. In the private chat no mention needed.
    """
    @messenger.send_text(message.chat_id, text, disable_web_page_preview: true)
    chat
  end

  @spec handle_global_competition_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_global_competition_cmd(_message, chat) do
    @messenger.send_text(chat.chat_id, "Now you see all girls", static_keyboard: :remove)
    %Chat{chat | competition: Voting.global_competition()}
  end

  @spec handle_celebrities_competition_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_celebrities_competition_cmd(_message, chat) do
    @messenger.send_text(
      chat.chat_id,
      "Now you see only celebrity-level girls",
      static_keyboard: :remove
    )
    %Chat{chat | competition: Voting.celebrities_competition()}
  end

  @spec handle_normal_competition_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_normal_competition_cmd(_message, chat) do
    @messenger.send_text(chat.chat_id, "Now you see only ordinary girls", static_keyboard: :remove)
    %Chat{chat | competition: Voting.normal_competition()}
  end

  @spec handle_get_girl_info_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_get_girl_info_cmd(message, chat) do
    girl_link = TextMessage.get_command_arg(message)
    if girl_link do
      case Voting.get_girl(chat.competition, girl_link) do
        {:ok, girl} -> display_girl_info(message.chat_id, girl)
        {:error, error_msg} -> @messenger.send_text(message.chat_id, error_msg)
      end
    else
      @messenger.send_text(
        message.chat_id,
        "Please send me a username\ngirlInfo@InstaRateBot <username>",
        disable_web_page_preview: true
      )
    end
    chat
  end

  @spec display_girl_info(integer, Girl.t) :: any
  defp display_girl_info(chat_id, girl) do
    profile_url = Girl.get_profile_url(girl)
    @messenger.send_markdown(chat_id, "[#{girl.username}](#{profile_url})")
    @messenger.send_photo(chat_id, girl.photo)
    girl_position = Girl.get_position(girl)
    text = """
    Position in the competition: #{girl_position}
    Number of wins: #{girl.wins}
    Number of loses: #{girl.loses}
    """
    @messenger.send_text(chat_id, text)
  end

  @spec on_callback(Callback.t, Chat.t) :: Chat.t
  defp on_callback(message, chat) do
    Logger.info("Process callback #{inspect message}")
    callback_name = Callback.get_name(message)
    handler = case callback_name do
      @vote_callback -> &handle_vote_callback/2
      @get_top_callback -> &handle_get_top_callback/2
      _ -> nil
    end
    if handler do
      Logger.info("Handle #{callback_name} callback")
      handler.(message, chat)
    else
      Logger.warn("Unknown callback: #{callback_name}")
      chat
    end
  end

  @spec try_to_send_pair(Chat.t) :: Chat.t
  defp try_to_send_pair(chat) do
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

  @spec handle_left_vote_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_left_vote_cmd(message, chat)  do
    if chat.last_match do
      winer_username = chat.last_match.left_girl
      loser_username = chat.last_match.right_girl
      process_vote_message(message, chat, winer_username, loser_username)
    else
      chat
    end
  end

  @spec handle_right_vote_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_right_vote_cmd(message, chat)  do
    if chat.last_match do
      loser_username = chat.last_match.left_girl
      winer_username = chat.last_match.right_girl
      process_vote_message(message, chat, winer_username, loser_username)
    else
      chat
    end
  end

  @spec process_vote_message(TextMessage.t, Chat.t, String.t, String.t) :: Chat.t
  defp process_vote_message(message, chat, winner_username, loser_username) do
    voters_group_id = build_voters_group_id(message.chat_id)
    voter_id = build_voter_id(message.user)
    case Voting.vote(
           chat.competition,
           voters_group_id,
           voter_id,
           winner_username,
           loser_username
         ) do
      :ok -> try_to_send_pair(chat)
      {:error, error} ->
        Logger.warn("Can't vote by message: #{error}")
        chat
    end
  end

  @spec handle_vote_callback(Callback.t, Chat.t) :: Chat.t
  defp handle_vote_callback(message, chat) do
    callback_args = Callback.get_args(message)
    [winner_username, loser_username] = String.split(callback_args, @usernames_separator)
    voters_group_id = build_voters_group_id(message.chat_id)
    voter_id = build_voter_id(message.user)
    case Voting.vote(
           chat.competition,
           voters_group_id,
           voter_id,
           winner_username,
           loser_username
         ) do
      :ok ->
        @messenger.send_notification(message.callback_id, "Vote for #{winner_username}")
        if chat.last_match.message_id == message.parent_msg_id do
          try_to_send_pair(chat)
        else
          chat
        end
      {:error, error} ->
        @messenger.send_notification(message.callback_id, "You already voted.")
        Logger.warn("Can't vote by callback: #{error}")
        chat
    end
  end

  @spec handle_get_top_callback(Callback.t, Chat.t) :: Chat.t
  defp handle_get_top_callback(_message, chat) do
    #    @messenger.delete_attached_keyboard(message.chat_id, message.parent_msg_id)
    #    callback_args = Callback.get_args(message)
    #    girl_offset = case Integer.parse(callback_args) do
    #      {offset, ""} -> offset
    #      _ -> raise "Non-int arg for get top callback: #{callback_args}"
    #    end
    #    if chat.current_top_offset == girl_offset - 1 do
    #      chat = send_girl_from_top(chat, girl_offset)
    #      @messenger.answer_callback(message.callback_id)
    #      chat
    #    else
    #      @messenger.send_notification(
    #        message.callback_id,
    #        "Please, continue from the most recent girl."
    #      )
    #      chat
    #    end
    chat
  end

  @spec send_girl_from_top(Chat.t, integer) :: Chat.t
  defp send_girl_from_top(chat, girl_offset) do
    case Voting.get_top(chat.competition, 2, offset: girl_offset) do
      [current_girl | next_girls] ->
        keyboard = if length(next_girls) != 0, do: [[@next_top_cmd]], else: :remove
        #        keyboard = if length(next_girls) != 0 do
        #          next_girl_offset = Integer.to_string(girl_offset + 1)
        #          [[%{text: "Next", payload: Callback.build_payload(@get_top_callback, next_girl_offset)}]]
        #        else
        #          nil
        #        end
        @messenger.send_photo(
          chat.chat_id,
          current_girl.photo,
          caption: "#{girl_offset + 1}th place: " <> Girl.get_profile_url(current_girl),
          #          keyboard: keyboard,
          static_keyboard: keyboard,
        )
        %Chat{chat | current_top_offset: girl_offset}
      [] ->
        Logger.warn("Girl offset #{girl_offset} more than number of girls in the competition")
        girls_number = Voting.get_girls_number(chat.competition)
        @messenger.send_text(
          chat.chat_id,
          "Sorry, but there are only #{girls_number} girls in the competition"
        )
        chat
    end
  end
end
