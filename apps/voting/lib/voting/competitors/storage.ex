defmodule Voting.Competitors.Storage do

  alias Voting.Competitors.Model, as: Competitor
  @type t :: module

  @callback get_top(competition :: String.t, number :: integer, offset :: integer) :: [Competitor.t]
  @callback get_random_pair(competition :: String.t) :: {Competitor.t, Competitor.t}
  @callback get_girl(competition :: String.t, username :: Competitor.t)
            :: {:ok, Competitor.t} | {:error, Competitor.t}
  @callback get_girls_number(competition :: String.t) :: integer
  @callback get_higher_ratings_number(competition :: String.t, rating :: integer) :: integer
  @callback update_girl(girl :: Competitor.t) :: Competitor.t
  @callback add_girl(girl :: Competitor.t) :: Competitor.t
end
