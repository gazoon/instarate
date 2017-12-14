defmodule Voting.Competitors.Storages.Mongo do

  alias Voting.Competitors.Model, as: Competitor
  @behaviour Voting.Competitors.Storage

  @collection "girls"
  @max_random_get_attempt 5

  @process_name :mongo_girls

  @spec child_spec :: tuple
  def child_spec do
    options = [name: @process_name] ++ Application.get_env(:voting, :mongo_girls)
    Utils.set_child_id(Mongo.child_spec(options), {Mongo, :girls})
  end

  @spec get_top(String.t, integer, integer) :: [Competitor.t]
  def get_top(competition, number, offset) do
    transform_girls(
      Mongo.find(
        @process_name,
        @collection,
        %{competition: competition},
        sort: %{
        rating: -1
        },
        limit: number,
        skip: offset
      )
    )
  end

  @spec get_girls_number(String.t) :: integer
  def get_girls_number(competition) do
    Mongo.count!(@process_name, @collection, %{competition: competition})
  end

  @spec get_girl(String.t, String.t) :: {:ok, Competitor.t} | {:error, String.t}
  def get_girl(competition, username) do
    row = Mongo.find_one(
      @process_name,
      @collection,
      %{competition: competition, username: username}
    )
    if row do
      {:ok, transform_girl(row)}
    else
      {:error, "Girl #{username} not found"}
    end
  end

  @spec get_higher_ratings_number(String.t, integer) :: integer
  def get_higher_ratings_number(competition, rating) do
    ratings = Mongo.distinct!(
      @process_name,
      @collection,
      "rating",
      %{
        competition: competition,
        rating: %{
          "$gt" => rating
        }
      }
    )
    length(ratings)
  end

  @spec update_girl(Competitor.t) :: Competitor.t
  def update_girl(girl) do
    Mongo.update_one!(
      @process_name,
      @collection,
      %{competition: girl.competition, username: girl.username},
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

  @spec add_girl(Competitor.t) :: Competitor.t
  def add_girl(girl) do
    insert_result = Mongo.insert_one(@process_name, @collection, girl)
    case insert_result do
      {:ok, _} -> girl
      {:error, error} -> raise error
    end
  end

  @spec get_random_pair(String.t) :: {Competitor.t, Competitor.t}
  def get_random_pair(competition) do
    attempt = 0
    get_random_pair(competition, attempt)
  end

  defp get_random_pair(_competition, _attempt = @max_random_get_attempt) do
    raise "Can't get two distinct rows, attempts limit is reached"
  end
  @spec get_random_pair(String.t, integer) :: {Competitor.t, Competitor.t}
  defp get_random_pair(competition, attempt) do
    girls = transform_girls(
      Mongo.aggregate(
        @process_name,
        @collection,
        [
          %{
            "$match" => %{
              competition: competition
            }
          },
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

  @spec transform_girls(Enum.t) :: [Competitor.t]
  defp transform_girls(rows) do
    Enum.map(rows, &transform_girl/1)
  end

  @spec transform_girl(map() | Mongo.Cursor.t) :: Competitor.t
  defp transform_girl(row) do
    %Competitor{
      username: row["username"],
      competition: row["competition"],
      rating: row["rating"],
      matches: row["matches"],
      wins: row["wins"],
      loses: row["loses"],
    }
  end
end
