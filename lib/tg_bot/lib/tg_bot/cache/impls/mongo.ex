defmodule TGBot.Cache.Impls.Mongo do
  alias TGBot.Cache.Behaviour, as: Cache
  @behaviour Cache

  @duplication_code 11_000
  @collection "insta_cache"
  @process_name :mongo_cache

  @spec child_spec :: tuple
  def child_spec do
    options = [name: @process_name, pool: DBConnection.Poolboy] ++
              Application.get_env(:tg_bot, :mongo_cache)
    Utils.set_child_id(Mongo.child_spec(options), {Mongo, :cache})
  end

  @spec get(String.t) :: Cache.value | nil
  def get(key) do
    row = Mongo.find_one(
      @process_name,
      @collection,
      %{key: key},
      pool: DBConnection.Poolboy
    )
    case row do
      nil -> nil
      %{"_single_value" => value} when map_size(row) == 3 -> value
      _ -> Map.pop(row, "key")
    end
  end

  @spec set(String.t, Cache.value) :: :ok
  def set(key, value) do
    doc = if is_map(value),
             do: Map.put(value, :key, key),
             else: %{
               key: key,
               _single_value: value
             }
    case Mongo.insert_one(@process_name, @collection, doc, pool: DBConnection.Poolboy) do
      {:error, error = %Mongo.Error{code: code}} when code != @duplication_code -> raise error
      _ -> :ok
    end
  end
end

