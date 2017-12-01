defmodule TGBot.Messages.Callback do
  @moduledoc false
  defstruct callback_id: nil, user_id: nil, chat_id: nil, payload: nil

  def type, do: :callback

end
