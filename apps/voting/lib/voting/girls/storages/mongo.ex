defmodule Voting.Girls.Storages.Mongo do

  alias Voting.Girls.Girl
  @behaviour Voting.Girls.Storage

  @collection "girls"
  @duplication_code 11_000
  @max_random_get_attempt 5

  @process_name :mongo_girls

  @spec child_spec :: tuple
  def child_spec do
    options = [name: @process_name] ++ Application.get_env(:voting, :mongo_girls)
    Utils.set_child_id(Mongo.child_spec(options), {Mongo, :girls})
  end

  @spec get_top(integer, integer) :: [Girl.t]
  def get_top(number, offset) do
    transform_girls(
      Mongo.find(
      @process_name,
      @collection,
      %{},
      sort: %{
        rating: -1
      },
      limit: number,
      skip: offset
      )
    )
  end

  @spec get_girls_number :: integer
  def get_girls_number do
    Mongo.count!(@process_name, @collection, %{})
  end

  @spec get_girl(String.t) :: {:ok, Girl.t} | {:error, String.t}
  def get_girl(username) do
    row = Mongo.find_one(@process_name, @collection, %{username: username})
    if row do
      {:ok, transform_girl(row)}
    else
      {:error, "Girl #{username} not found"}
    end
  end

  @spec get_higher_ratings_number(integer) :: integer
  def get_higher_ratings_number(rating) do
    ratings = Mongo.distinct!(
      @process_name,
      @collection,
      "rating",
      %{
        rating: %{
          "$gt" => rating
        }
      }
    )
    length(ratings)
  end

  @spec update_girl(Girl.t) :: Girl.t
  def update_girl(girl) do
    Mongo.update_one!(
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
      {:error, error} -> raise error
    end
  end

  @spec get_random_pair :: {Girl.t, Girl.t}
  def get_random_pair do
    attempt = 0
    get_random_pair(attempt)
  end

  defp get_random_pair(_attempt = @max_random_get_attempt) do
    raise "Can't get two distinct rows, attempts limit is reached"
  end
  @spec get_random_pair(integer) :: {Girl.t, Girl.t}
  defp get_random_pair(attempt) do
    girls = transform_girls(
      Mongo.aggregate(
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
    )

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

  @spec transform_girl(map() | Mongo.Cursor.t) :: Girl.t
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
