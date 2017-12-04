defmodule TGBot.Messages.Text do

  alias TGBot.Messages.Text, as: TextMessage
  @type t :: %TextMessage{
               text: String.t,
               text_lowercase: String.t,
               chat_id: integer,
               user_id: integer,
               is_group_chat: boolean,
               message_id: integer
             }
  defstruct text: "",
            text_lowercase: "",
            chat_id: nil,
            user_id: nil,
            is_group_chat: true,
            message_id: nil
  def type, do: :text

  @spec from_data(map()) :: TextMessage.t
  def from_data(message_data) do
    message = struct(TextMessage, message_data)
    %TextMessage{message | text_lowercase: String.downcase(message_data.text)}
  end

  @spec appeal_to_bot?(TextMessage.t) :: boolean
  def appeal_to_bot?(message) do
    bot_name = String.downcase(Application.get_env(:tg_bot, :bot_name))
    String.contains?(message.text_lowercase, bot_name)
  end

  @spec get_command_arg(TextMessage.t) :: String.t | nil
  def get_command_arg(message) do
    List.last(String.split(message.text, " "))
  end
end
