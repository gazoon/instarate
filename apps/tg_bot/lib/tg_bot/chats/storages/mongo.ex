defmodule TGBot.Chats.Storages.Mongo do
  @behaviour TGBot.Chats.Storage
  alias TGBot.Chats.Chat

  @collection "insta_chats"
  @process_name :mongo_chats

  @spec child_spec :: tuple
  def child_spec do
    options = [name: @process_name] ++ Application.get_env(:tg_bot, :mongo_chats)
    Utils.set_child_id(Mongo.child_spec(options), {Mongo, :chats})
  end

  @spec get(integer) :: Chat.t | nil
  def get(chat_id) do
    row = Mongo.find_one(@process_name, @collection, %{chat_id: chat_id})
    transform_chat(row)
  end

  @spec save(Chat.t) :: Chat.t
  def save(chat) do
    chat_data = Map.from_struct(chat)
    chat_data = if chat.last_match,
                   do: %{chat_data | last_match: Map.from_struct(chat.last_match)},
                   else: chat_data
    Mongo.replace_one!(
      @process_name,
      @collection,
      %{chat_id: chat.chat_id},
      chat_data,
      upsert: true
    )
    chat
  end

  defp transform_chat(nil), do: nil
  @spec transform_chat(map()) :: Chat.t
  defp transform_chat(row) do
    last_match_data = row["last_match"]
    last_match = if last_match_data do
      %Chat.Match{
        message_id: last_match_data["message_id"],
        left_girl: last_match_data["left_girl"],
        right_girl: last_match_data["right_girl"],
        shown_at: last_match_data["shown_at"],
      }
    else
      nil
    end
    %Chat{
      chat_id: row["chat_id"],
      current_top_offset: row["current_top_offset"] || 0,
      last_match: last_match,
      competition: row["competition"],
    }
  end
end
