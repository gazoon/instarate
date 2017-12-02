defmodule TGBot.Messages.Callback do
  @moduledoc false
  alias TGBot.Messages.Callback
  @type t :: %Callback{
               callback_id: String.t,
               chat_id: integer,
               user_id: integer,
               is_group_chat: boolean,
               payload: map
             }
  defstruct callback_id: nil, user_id: nil, chat_id: nil, is_group_chat: true, payload: nil

  def type, do: :callback

  @spec from_data(map()) :: Callback.t
  def from_data(message_data) do
    struct(Callback, message_data)
  end

end
