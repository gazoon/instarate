defmodule Voting.Girls.Storages.Mongo do
  @moduledoc false
  alias Voting.Girls.Girl
  @behaviour Voting.Girls.Storage

  @collection "girls"
  @duplication_code 11000
  @max_random_get_attempt 5

  @process_name :mongo_girls
  def process_name, do: @process_name

  @spec get_top(integer) :: [Girl.t]
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

  @spec get_girl(String.t) :: Girl.t | nil
  def get_girl(username) do
    Mongo.find_one(@process_name, @collection, %{username: username})
    |> transform_girl()
  end

  @spec update_girl(Girl.t) :: Girl.t
  def update_girl(girl) do
    Mongo.update_one(
      @process_name,
      @collection,
      %{username: girl.username},
      %{
        "$set" => %{
          rating: girl.rating,
          photo: girl.photo,
          matches: girl.matches,
          wins: girl.wins,
          loses: girl.loses,
        }
      }
    )
    girl
  end

  @spec add_girl(Girl.t) :: {:ok, Girl.t} | {:error, String.t}
  def add_girl(girl) do
    insert_result = Mongo.insert_one(@process_name, @collection, girl)
    case insert_result do
      {:ok, _} ->
        {:ok, girl}
      {:error, %Mongo.Error{code: @duplication_code}} ->
        {:error, "Girl #{girl.username} already added"}
    end
  end

  @spec get_random_pair() :: {Girl.t, Girl.t}
  def get_random_pair() do
    attempt = 0
    get_random_pair(attempt)
  end

  defp get_random_pair(attempt = @max_random_get_attempt) do
    raise RuntimeError, message: "Can't get two distinct rows, attempts limit is reached"
  end
  @spec get_random_pair(integer) :: {Girl.t, Girl.t}
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

  @spec transform_girls(Enum.t) :: [Girl.t]
  defp transform_girls(rows) do
    Enum.map(rows, &transform_girl/1)
  end

  defp transform_girl(nil), do: nil

  @spec transform_girl(map()) :: Girl.t
  defp transform_girl(row) do
    %Girl{
      username: row["username"],
      added_at: row["added_at"],
      photo: row["photo"],
      rating: row["rating"],
      matches: row["matches"],
      wins: row["wins"],
      loses: row["loses"],
    }
  end

end
