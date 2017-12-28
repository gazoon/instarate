defmodule Utils.Messages.Text do
  @behaviour Utils.Messages.Builder

  alias Utils.Messages.Text, as: TextMessage
  alias Utils.Messages.User, as: MessageUser
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

  @spec reply_to_bot?(TextMessage.t, String.t) :: boolean
  def reply_to_bot?(message, bot_username) do
    reply_to = message.reply_to
    if reply_to, do: MessageUser.is_bot?(reply_to.user, bot_username), else: false
  end

  @spec appeal_to_bot?(TextMessage.t, String.t) :: boolean
  def appeal_to_bot?(message, bot_name) do
    bot_name = String.downcase(bot_name)
    String.contains?(message.text_lowercase, bot_name)
  end

  @spec get_command_arg(TextMessage.t) :: String.t | nil
  def get_command_arg(message) do
    List.last(get_command_args(message))
  end

  @spec get_command_args(TextMessage.t) :: [String.t]
  def get_command_args(message) do
    tokens = String.split(message.text, " ")
    [_ | args] = tokens
    args
  end
end
