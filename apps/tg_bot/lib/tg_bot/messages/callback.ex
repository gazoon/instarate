defmodule TGBot.Messages.Callback do

  alias TGBot.Messages.Callback
  @type t :: %Callback{
               callback_id: String.t,
               chat_id: integer,
               user_id: integer,
               parent_msg_id: integer,
               is_group_chat: boolean,
               payload: map
             }
  defstruct callback_id: nil,
            user_id: nil,
            chat_id: nil,
            parent_msg_id: nil,
            is_group_chat: true,
            payload: nil

  def type, do: :callback
  @callback_name_separator ":"

  @spec from_data(map()) :: Callback.t
  def from_data(message_data) do
    struct(Callback, message_data)
  end

  @spec get_name(Callback.t) :: String.t
  def get_name(callback) do
    callback
    |> split_payload()
    |> List.first()
  end

  @spec get_args(Callback.t) :: String.t
  def get_args(callback) do
    callback
    |> split_payload()
    |> List.last()
  end

  @spec build_payload(String.t, String.t) :: String.t
  def build_payload(callback_name, args) do
    callback_name <> @callback_name_separator <> args
  end

  @spec split_payload(Callback.t) :: [String.t]
  defp split_payload(callback) do
    callback.payload
    |> String.split(@callback_name_separator)
  end

end
