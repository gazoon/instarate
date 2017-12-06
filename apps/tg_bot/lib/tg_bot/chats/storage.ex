defmodule TGBot.Chats.Storage do
  alias TGBot.Chats.Chat
  @type t :: module

  @callback get(chat_id :: integer) :: Chat.t | nil
  @callback save(chat :: Chat.t) :: Chat.t
end
