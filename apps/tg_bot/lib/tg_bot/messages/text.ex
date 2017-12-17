defmodule TGBot.Messages.Text do
  @behaviour TGBot.MessageBuilder

  alias TGBot.Messages.Text, as: TextMessage
  alias TGBot.Messages.User, as: MessageUser
  @type t :: %TextMessage{
               text: String.t,
               text_lowercase: String.t,
               chat_id: integer,
               user: MessageUser.t,
               is_group_chat: boolean,
               message_id: integer,
               reply_to: TextMessage.t | nil
             }
  defstruct text: "",
            text_lowercase: "",
            chat_id: nil,
            user: nil,
            is_group_chat: true,
            message_id: nil,
            reply_to: nil

  def type, do: :text

  @spec from_data(map()) :: TextMessage.t
  def from_data(message_data) do
    message_data = Utils.keys_to_atoms(message_data)
    {reply_to_data, message_data} = Map.pop(message_data, :reply_to)
    {user_data, message_data} = Map.pop(message_data, :user)
    user = MessageUser.from_data(user_data)
    reply_to = if reply_to_data, do: from_data(reply_to_data), else: nil
    message = struct(TextMessage, message_data)
    %TextMessage{
      message |
      text_lowercase: String.downcase(message_data.text),
      user: user,
      reply_to: reply_to
    }
  end

  @spec reply_to_bot?(TextMessage.t) :: boolean
  def reply_to_bot?(message) do
    reply_to = message.reply_to
    if reply_to, do: MessageUser.is_bot?(reply_to.user), else: false
  end

  @spec appeal_to_bot?(TextMessage.t) :: boolean
  def appeal_to_bot?(message) do
    bot_name = String.downcase(Application.get_env(:tg_bot, :bot_name))
    String.contains?(message.text_lowercase, bot_name)
  end

  @spec get_command_arg(TextMessage.t) :: String.t | nil
  def get_command_arg(message) do
    tokens = String.split(message.text, " ")
    [_ | args] = tokens
    List.last(args)
  end
end
