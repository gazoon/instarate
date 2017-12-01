defmodule TGBot.Messages.Text do
  @moduledoc false
  defstruct text: "", chat_id: nil, user_id: nil, message_id: nil
  def type, do: :text
end
