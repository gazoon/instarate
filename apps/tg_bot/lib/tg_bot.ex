defmodule TGBot do
  @moduledoc false
  require Logger
  alias TGBot.Messages.Text, as: TextMessage
  alias TGBot.Messages.Callback, as: CallbackMessage
  alias Voting.Girls.Girl
  alias TGBot.UserMessage
  alias TGBot.Messenger
  @start_cmd "start"
  @add_girl_cmd "addgirl"
  @get_top_cmd "showtop"
  @get_girl_info_cmd "girlinfo"
  @voter_prefix  "tg:"
  @top_page_size 3

  @spec on_message(map()) :: any
  def on_message(message_container) do
    message_type = message_container.type
    message_data = message_container.data
    case message_type do
      :text ->
        message = TextMessage.from_data(message_data)
        on_text_message(message)
      :callback ->
        message = CallbackMessage.from_data(message_data)
        on_callback(message)
      _ -> Logger.error("Unsupported message type: #{message_type}")
    end
  end

  @spec on_text_message(TextMessage.t) :: any
  defp on_text_message(message) do
    if TextMessage.appeal_to_bot?(message) || !message.is_group_chat do
      process_text_message(message)
    else
      Logger.info("Skip message #{inspect message} it's not an appeal")
    end
  end

  @spec process_text_message(TextMessage.t) :: any
  defp process_text_message(message) do
    commands = [
      {@start_cmd, &handle_start_cmd/1},
      {@add_girl_cmd, &handle_add_girl_cmd/1},
      {@get_top_cmd, &handle_get_top_cmd/1},
      {@get_girl_info_cmd, &handle_get_girl_info_cmd/1},
    ]
    message_text = message.text_lowercase
    command = commands
              |> Enum.find(fn ({cmd_name, _}) -> String.contains?(message_text, cmd_name) end)
    case command do
      {command_name, handler} -> Logger.info("Handle #{command_name} command")
                                 handler.(message)
      nil -> Logger.info("Message #{message_text} doesn't contain commands")
    end

    IO.inspect(message)
  end

  @spec build_voter_id(UserMessage.t) :: any
  defp build_voter_id(user_message) do
    user_id = UserMessage.user_id(user_message)
    @voter_prefix <> Integer.to_string(user_id)
  end

  @spec handle_start_cmd(TextMessage.t) :: any
  defp handle_start_cmd(message) do
    voter_id = build_voter_id(message)
    {girl_one, girl_two} = Voting.get_next_pair(voter_id)
    IO.inspect(girl_one)
    IO.inspect(girl_two)
  end

  @spec handle_add_girl_cmd(TextMessage.t) :: any
  defp handle_add_girl_cmd(message) do
    photo_link = TextMessage.get_command_arg(message)
    case Voting.add_girl(photo_link) do
      {:ok, girl} ->
        profile_url = Girl.get_profile_url(girl)
        text = "Girl [#{girl.username}](#{profile_url}) has been successfully added!"
        Messenger.send_markdown(message.chat_id, text)
      {:error, error_msg} -> Nadia.send_message(message.chat_id, error_msg)
    end
  end

  @spec handle_get_top_cmd(TextMessage.t) :: any
  defp handle_get_top_cmd(message) do
    Voting.get_top(@top_page_size)
    |> Enum.with_index(1)
    |> Enum.each(
         fn ({girl, i}) ->
           Messenger.send_text(message.chat_id, "#{i}th place:")
           Messenger.send_photo(message.chat_id, girl.photo, caption: Girl.get_profile_url(girl))
         end
       )
  end

  @spec handle_get_girl_info_cmd(TextMessage.t) :: any
  defp handle_get_girl_info_cmd(message) do
    girl_link = TextMessage.get_command_arg(message)
    case Voting.get_girl(girl_link) do
      {:ok, girl} -> display_girl_info(message.chat_id, girl)
      {:error, error_msg} -> Nadia.send_message(message.chat_id, error_msg)
    end
  end

  @spec display_girl_info(integer, Girl.t) :: any
  defp display_girl_info(chat_id, girl) do
    profile_url = Girl.get_profile_url(girl)
    Messenger.send_markdown(chat_id, "[#{girl.username}](#{profile_url})")
    Messenger.send_photo(chat_id, girl.photo)
    girl_position = Girl.get_position(girl)
    text = """
    Position in the competition: #{girl_position}
    Number of wins: #{girl.wins}
    Number of loses: #{girl.loses}
    """
    Messenger.send_text(chat_id, text)
  end

  defp on_callback(message) do

    IO.inspect("caaaaaaaaaaaaalllllback2222")
    IO.inspect(message)
  end

end
