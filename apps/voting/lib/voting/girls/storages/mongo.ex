defmodule Voting.Girls.Storages.Mongo do
  @moduledoc false
  alias Voting.Girls.Girl
  @behaviour Voting.Girls.Storage

  @collection "girls"
  @duplication_code 11000
  @max_random_get_attempt 5

  @process_name :mongo_girls
  def process_name, do: @process_name

  def get_top(number) do
    Mongo.find(
      @process_name,
      @collection,
      %{},
      sort: %{
        rating: -1
      },
      limit: number
    )
    |> transform_girls()
  end

  def get_girl(username) do
    Mongo.find_one(@process_name, @collection, %{username: username})
    |> transform_girl()
  end

  def update_girl(girl) do
    Mongo.update_one(
      @process_name,
      @collection,
      %{username: girl.username},
      %{
        "$set" => %{
          rating: girl.rating,
          photo: girl.photo
        }
      }
    )
    girl
  end

  def add_girl(girl) do
    insert_result = Mongo.insert_one(@process_name, @collection, girl)
    case insert_result do
      {:ok, _} -> :ok
      {:error, %Mongo.Error{code: @duplication_code}} -> :error
    end
  end

  def get_random_pair() do
    attempt = 0
    get_random_pair(attempt)
  end

  defp get_random_pair(attempt = @max_random_get_attempt) do
    raise RuntimeError, message: "Can't get two distinct rows, attempts limit is reached"
  end
  defp get_random_pair(attempt) do
    girls = Mongo.aggregate(
              @process_name,
              @collection,
              [
                %{
                  "$sample" => %{
                    size: 2
                  }
                }
              ]
            )
            |> transform_girls()
    [girl_one, girl_two] = girls
    if girl_one.username != girl_two.username do
      {girl_one, girl_two}
    else
      get_random_pair(attempt + 1)
    end
  end

  defp transform_girls(rows) do
    Enum.map(rows, &transform_girl/1)
  end

  defp transform_girl(nil), do: nil
  defp transform_girl(row) do
    %Girl{
      username: row["username"],
      added_at: row["added_at"],
      photo: row["photo"],
      rating: row["rating"]
    }
  end

end
