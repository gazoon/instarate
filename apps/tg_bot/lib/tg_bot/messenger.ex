defmodule TGBot.Messenger do
  @callback send_text(chat_id :: integer, text :: String.t, opts :: Keyword.t) :: integer
  @callback send_markdown(chat_id :: integer, text :: String.t, opts :: Keyword.t) :: integer
  @callback get_chat_members_number(chat_id :: integer) :: integer
  @callback send_photo(chat_id :: integer, photo :: binary, opts :: Keyword.t)
            :: {integer, String.t}
  @callback send_notification(callback_id :: String.t, text :: String.t) :: any
  @callback answer_callback(callback_id :: String.t) :: any
  @callback delete_attached_keyboard(chat_id :: integer, message_id :: integer) :: integer
  @callback send_photos(chat_id :: integer, photos :: [map()]) :: any
end
