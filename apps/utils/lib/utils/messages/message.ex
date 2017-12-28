alias  Utils.Messages.Message
alias  Utils.Messages.Text, as: TextMessage
alias  Utils.Messages.Callback, as: CallbackMessage
alias Utils.Messages.Task

defprotocol Utils.Messages.Message do
  def chat_id(message)
end

defimpl Message, for: TextMessage do
  def chat_id(message), do: message.chat_id
end

defimpl Message, for: CallbackMessage do
  def chat_id(message), do: message.chat_id
end

defimpl Message, for: Task do
  def chat_id(message), do: message.chat_id
end

