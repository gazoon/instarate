defmodule Voting.Competitors.Model do

  alias Voting.Competitors.Model, as: Competitor

  @initial_rating 1500
  @storage Application.get_env(:voting, __MODULE__)[:storage]

  @type t :: %Competitor{
               username: String.t,
               competition: String.t,
               rating: integer,
               matches: integer,
               wins: integer,
               loses: integer
             }

  defstruct username: nil,
            competition: nil,
            rating: @initial_rating,
            matches: 0,
            wins: 0,
            loses: 0

  @spec new(String.t, String.t) :: Competitor.t
  def new(competition, username) do
    %Competitor{competition: competition, username: username}
  end

  @spec get_position(Competitor.t) :: integer
  def get_position(girl) do
    @storage.get_higher_ratings_number(girl.competition, girl.rating) + 1
  end
end
