defmodule TGBot.Chats.Storages.Mongo do
  @behaviour TGBot.Chats.Storage
  alias TGBot.Chats.Chat

  @collection "girls"
  @process_name :mongo_girls
  def process_name, do: @process_name

  @spec get(integer) :: Chat.t | nil
  def get(chat_id) do
    row = Mongo.find_one(@process_name, @collection, %{chat_id: chat_id})
    transform_chat(row)
  end

  @spec save(Chat.t) :: Chat.t
  def save(chat) do
    Mongo.replace_one(@process_name, @collection, %{chat_id: chat.chat_id}, chat, upsert: true)
    chat
  end

  defp transform_chat(nil), do: nil
  @spec transform_chat(map()) :: Chat.t
  defp transform_chat(row) do
    last_match_data = row["last_match"]
    %Chat{
      chat_id: row["chat_id"],
      current_top_offset: row["current_top_offset"] || 0,
      last_match: %Chat.Match{
        message_id: last_match_data["message_id"],
        left_girl: last_match_data["left_girl"],
        right_girl: last_match_data["right_girl"],
        shown_at: last_match_data["shown_at"],
      },
    }
  end
end
