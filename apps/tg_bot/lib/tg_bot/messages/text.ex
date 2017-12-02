defmodule TGBot.Messages.Text do
  @moduledoc false
  alias TGBot.Messages.Text, as: TextMessage
  defstruct text: "",
            text_lowercase: "",
            chat_id: nil,
            user_id: nil,
            is_group_chat: true,
            message_id: nil
  def type, do: :text

  def from_data(message_data) do
    message = struct(TextMessage, message_data)
    %TextMessage{message | text_lowercase: String.downcase(message_data.text)}
  end
  def appeal_to_bot?(message) do
    bot_name = String.downcase(Application.get_env(:tg_bot, :bot_name))
    String.contains?(message.text_lowercase, bot_name)
  end
end
