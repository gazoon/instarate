defmodule Voting.Fakes.Competitors do
  alias Voting.Competitors.Model, as: Competitor
  @behaviour Voting.Competitors.Storage


  @spec update_girl(Competitor.t) :: Competitor.t
  def update_girl(girl), do: girl

  @spec add_girl(Competitor.t) :: Competitor.t
  def add_girl(girl) do
    send(self(), {:competitor_added, girl})
  end

  @spec get_random_pair(String.t) :: {Competitor.t, Competitor.t}
  def get_random_pair(competition) do
    {
      %Competitor{username: "first_girl", competition: competition},
      %Competitor{username: "second_girl", competition: competition},
    }
  end

  # not used
  @spec get_girl(String.t, String.t) :: {:ok, Competitor.t} | {:error, String.t}
  def get_girl(competition, username) do
    {:ok, %Competitor{username: username, competition: competition, rating: 1500}}
  end

  # not used
  @spec get_top(String.t, integer, integer) :: [Competitor.t]
  def get_top(_competition, _number, _offset), do: []

  # not used
  @spec delete_girls([String.t]) :: :ok
  def delete_girls(_usernames), do: :ok

  # not used
  @spec get_girls_number(String.t) :: integer
  def get_girls_number(_competition), do: :ok

  # not used
  @spec get_higher_ratings_number(String.t, integer) :: integer
  def get_higher_ratings_number(_competition, _rating), do: 0
end
