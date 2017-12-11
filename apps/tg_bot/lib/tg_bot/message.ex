alias  TGBot.Message
alias  TGBot.Messages.Text, as: TextMessage
alias  TGBot.Messages.Callback, as: CallbackMessage
alias TGBot.Messages.Task

defprotocol TGBot.Message do
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

defmodule TGBot.MessageBuilder do
  @type t :: module
  @callback from_data(data :: map()) :: Message.t
end

