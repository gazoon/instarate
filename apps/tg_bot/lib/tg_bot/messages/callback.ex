defmodule TGBot.Messages.Callback do
  @moduledoc false
  alias TGBot.Messages.Callback
  defstruct callback_id: nil, user_id: nil, chat_id: nil, is_group_chat: true, payload: nil

  def type, do: :callback

  def from_data(message_data) do
    struct(Callback, message_data)
  end

end
