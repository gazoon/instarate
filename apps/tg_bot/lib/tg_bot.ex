defmodule TGBot do

  require Logger
  alias TGBot.Messages.Text, as: TextMessage
  alias TGBot.Messages.Callback, as: Callback
  alias TGBot.Messages.User, as: MessageUser
  alias TGBot.{Message, Pictures}
  alias TGBot.Chats.Chat
  alias Voting.Girls.Girl

  @start_cmd "start"
  @add_girl_cmd "addgirl"
  @get_top_cmd "showtop"
  @next_top_cmd "Next girl"
  @get_girl_info_cmd "girlinfo"
  @help_cmd "help"

  @vote_callback "vt"
  @get_top_callback "top"

  @usernames_separator "|"

  @config Application.get_env(:tg_bot, __MODULE__)
  @chats_storage @config[:chats_storage]
  @messenger @config[:messenger]

  @spec on_message(map()) :: any
  def on_message(message_container) do
    message_type = message_container.type
    message_data = message_container.data
    message_info = case message_type do
      :text ->
        {TextMessage, &on_text_message/2}

      :callback ->
        {Callback, &on_callback/2}
      _ -> nil
    end
    case message_info do
      {message_cls, handler_func} ->
        message = message_cls.from_data(message_data)
        process_message(message, handler_func)
      nil -> Logger.error("Unsupported message type: #{message_type}")
    end
  end

  @spec process_message(Message.t, ((Message.t, Chat.t) -> Chat.t)) :: any
  def process_message(message, handler) do
    chat = get_chat(message)
    chat_after_processing = handler.(message, chat)
    if chat_after_processing != chat do
      Logger.info("Save updated chat info")
      @chats_storage.save(chat_after_processing)
    end
  end

  @spec get_chat(Message.t) :: Chat.t
  defp get_chat(message) do
    chat_id = Message.chat_id(message)
    case @chats_storage.get(chat_id) do
      nil -> Chat.new(chat_id)
      chat -> chat
    end
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
        Logger.info("Message #{message_text} doesn't contain commands")
        chat
    end
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
    send_next_girls_pair(chat)
  end

  @spec send_next_girls_pair(Chat.t) :: Chat.t
  defp send_next_girls_pair(chat) do
    voters_group_id = build_voters_group_id(chat.chat_id)
    {girl_one, girl_two} = Voting.get_next_pair(voters_group_id)
    girl_one_url = Girl.get_profile_url(girl_one)
    girl_two_url = Girl.get_profile_url(girl_two)

    match_photo = Pictures.concatenate(girl_one.photo, girl_two.photo)

    keyboard = [
      [
        %{
          text: "Left",
          payload: Callback.build_payload(
            @vote_callback,
            girl_one.username <> @usernames_separator <> girl_two.username
          )
        },
        %{
          text: "Right",
          payload: Callback.build_payload(
            @vote_callback,
            girl_two.username <> @usernames_separator <> girl_one.username
          )
        },
      ]
    ]
    caption_text = "#{girl_one_url} vs #{girl_two_url}"
    message_id = try do
      @messenger.send_photo(
        chat.chat_id,
        match_photo,
        keyboard: keyboard,
        caption: caption_text
      )
    after
      File.rm!(match_photo)
    end
    current_match = Chat.Match.new(message_id, girl_one.username, girl_two.username)
    %Chat{chat | last_match: current_match}
  end

  @spec handle_add_girl_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_add_girl_cmd(message, chat) do
    photo_link = TextMessage.get_command_arg(message)
    case Voting.add_girl(photo_link) do
      {:ok, girl} ->
        profile_url = Girl.get_profile_url(girl)
        text = "Girl [#{girl.username}](#{profile_url}) has been successfully added!"
        @messenger.send_markdown(message.chat_id, text)
      {:error, error_msg} -> @messenger.send_text(message.chat_id, error_msg)
    end
    chat
  end
  #  @spec handle_get_top_cmd(TextMessage.t) :: any
  #  defp handle_get_top_cmd(message) do
  #    Voting.get_top(@top_page_size)
  #    |> Enum.with_index(1)
  #    |> Enum.each(
  #         fn ({girl, i}) ->
  #           @messenger.send_text(message.chat_id, "#{i}th place:")
  #           @messenger.send_photo(message.chat_id, girl.photo, caption: Girl.get_profile_url(girl))
  #         end
  #       )
  #  end

  #  @spec handle_get_top_cmd(TextMessage.t) :: any
  #  defp handle_get_top_cmd(message) do
  #    0..5
  #    |> Enum.each(
  #         fn (part) ->
  #           start_index = part * @top_page_size
  #           photos = Voting.get_top(@top_page_size, offset: start_index)
  #                    |> Enum.map(
  #                         fn (girl) -> %{url: girl.photo, caption: Girl.get_profile_url(girl)} end
  #                       )
  #           @messenger.send_photos(message.chat_id, photos)
  #           :timer.sleep(5000)
  #         end
  #       )
  #  end

  @spec handle_get_top_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_get_top_cmd(message, chat) do
    optional_start_position = TextMessage.get_command_arg(message)
    offset = case Integer.parse(optional_start_position) do
      {start_position, ""} when start_position > 0 -> start_position - 1
      _ -> 0
    end
    send_girl_from_top(chat, offset)
  end

  @spec handle_next_top_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_next_top_cmd(_message, chat) do
    next_offset = chat.current_top_offset + 1
    send_girl_from_top(chat, next_offset)
  end

  @spec handle_help_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_help_cmd(message, chat) do
    text = """
    Hi there! My mission is to find the most attractive girls on Instagram!
    Just select which of two girls looks better and vote by pressing a button below.
    I support following commands:

    /start - Get the next girls pair to compare.

    /showtop - Show the top girls in the competition. You can also pass a position to start showing from, i.e pass 10 to start from the 10th girl in the competition and skip the first 9.

    /help - Show this message.

    And also another type of commands, with an additional input:

    addgirl <link to one of her photos on instagram> - Add girl to the competition, You can add any girl, just paste a link of her instagram photo. The girl must have a public account. For example:
    addgirl https://www.instagram.com/p/BcPqz6sFMbb/

    girlinfo <girl username or link to the profile> - Get info about particular girl statistic in the competition. For example:
    girlinfo https://www.instagram.com/svetabily/
    girlinfo svetabily

    Some notes:
    In group chats you have to mention me (with the '@' sign) in the message, if you want to command me. In the private chat no mention needed.
    """
    @messenger.send_text(message.chat_id, text, disable_web_page_preview: true)
    chat
  end

  @spec handle_get_girl_info_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_get_girl_info_cmd(message, chat) do
    girl_link = TextMessage.get_command_arg(message)
    case Voting.get_girl(girl_link) do
      {:ok, girl} -> display_girl_info(message.chat_id, girl)
      {:error, error_msg} -> @messenger.send_text(message.chat_id, error_msg)
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
    callbacks = %{
      @vote_callback => &handle_vote_callback/2,
      @get_top_callback => &handle_get_top_callback/2,
    }
    callback_name = Callback.get_name(message)
    handler = callbacks[callback_name]
    if handler do
      Logger.info("handle #{callback_name} callback")
      handler.(message, chat)
    else
      Logger.warn("Unknown callback: #{callback_name}")
      chat
    end
  end

  @spec handle_vote_callback(Callback.t, Chat.t) :: Chat.t
  defp handle_vote_callback(message, chat) do
    callback_args = Callback.get_args(message)
    [winner_username, loser_username] = String.split(callback_args, @usernames_separator)
    voters_group_id = build_voters_group_id(message.chat_id)
    voter_id = build_voter_id(message.user)
    case Voting.vote(voters_group_id, voter_id, winner_username, loser_username) do
      :ok ->
        @messenger.send_notification(message.callback_id, "Vote for #{winner_username}")
        if chat.last_match.message_id == message.parent_msg_id do
          send_next_girls_pair(chat)
        else
          chat
        end
      {:error, error} ->
        @messenger.send_notification(message.callback_id, "You already voted.")
        Logger.warn("Can't vote: #{error}")
        chat
    end
  end

  @spec handle_get_top_callback(Callback.t, Chat.t) :: Chat.t
  defp handle_get_top_callback(message, chat) do
    #    @messenger.delete_attached_keyboard(message.chat_id, message.parent_msg_id)
    callback_args = Callback.get_args(message)
    girl_offset = case Integer.parse(callback_args) do
      {offset, ""} -> offset
      _ -> raise "Non-int arg for get top callback: #{callback_args}"
    end
    if chat.current_top_offset == girl_offset - 1 do
      chat = send_girl_from_top(chat, girl_offset)
      @messenger.answer_callback(message.callback_id)
      chat
    else
      @messenger.send_notification(
        message.callback_id,
        "Please, continue from the most recent girl."
      )
      chat
    end
  end

  @spec send_girl_from_top(Chat.t, integer) :: Chat.t
  defp send_girl_from_top(chat, girl_offset) do
    case Voting.get_top(2, offset: girl_offset) do
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
        girls_number = Voting.get_girls_number()
        @messenger.send_text(
          chat.chat_id,
          "Sorry, but there are only #{girls_number} girls in the competition"
        )
        chat
    end
  end
end
