defmodule TGBot.Chats.Storage do
  alias TGBot.Chats.Chat

  @callback get(chat_id :: integer) :: Chat.t | nil
  @callback save(chat :: Chat.t) :: Chat.t
end
