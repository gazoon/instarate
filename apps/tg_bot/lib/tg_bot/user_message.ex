defprotocol TGBot.UserMessage do
  def user(message)
  def chat_id(message)
  def is_group_chat(message)
end

alias  TGBot.UserMessage
alias  TGBot.Messages.Text, as: TextMessage
alias  TGBot.Messages.Callback, as: CallbackMessage

defimpl UserMessage, for: TextMessage do
  def user(message), do: message.user
  def chat_id(message), do: message.chat_id
  def is_group_chat(message), do: message.is_group_chat
end

defimpl UserMessage, for: CallbackMessage do
  def user(message), do: message.user
  def chat_id(message), do: message.chat_id
  def is_group_chat(message), do: message.is_group_chat
end
