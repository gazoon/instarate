defmodule TGBot do
  @moduledoc false
  require Logger
  alias TGBot.Messages.Text, as: TextMessage
  alias TGBot.Messages.Callback, as: CallbackMessage
  @start_cmd "start"
  @add_girl_cmd "addgirl"
  @get_top_cmd "showtop"
  @get_girl_info_cmd "girlinfo"

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

  defp on_text_message(message) do
    if TextMessage.appeal_to_bot?(message) || !message.is_group_chat do
      process_text_message(message)
    else
      Logger.info("Skip message #{inspect message} it's not an appeal")
    end
  end

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

  defp handle_start_cmd(message) do

  end

  defp handle_add_girl_cmd(message) do

  end

  defp handle_get_top_cmd(message) do

  end

  defp handle_get_girl_info_cmd(message) do

  end

  defp on_callback(message) do

    IO.inspect("caaaaaaaaaaaaalllllback2222")
    IO.inspect(message)
  end

end
