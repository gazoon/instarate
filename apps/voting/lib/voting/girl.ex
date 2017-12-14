defmodule Voting.Girl do

  alias Voting.Girl
  alias Instagram.Client, as: InstagramClient
  alias Voting.Competitors.Model, as: Competitor
  alias Voting.InstagramProfiles.Model, as: Profile
  @storage Application.get_env(:voting, __MODULE__)[:storage]

  @type t :: %Girl{
               username: String.t,
               photo: String.t,
               competition: String.t,
               rating: integer,
               matches: integer,
               wins: integer,
               loses: integer
             }


  @spec combine(Competitor.t, Profile.t) :: Girl.t
  def combine(competitor, profile) do
    %Girl{
      username: profile.username,
      photo: profile.photo,
      competition: competitor.competition,
      rating: competitor.rating,
      matches: competitor.matches,
      wins: competitor.wins,
      loses: competitor.loses,
    }
  end

  @spec get_profile_url(Girl.t) :: String.t
  def get_profile_url(girl) do
    InstagramClient.build_profile_url(girl.username)
  end

  @spec get_position(Girl.t) :: integer
  def get_position(girl) do
    @storage.get_higher_ratings_number(girl.competition, girl.rating) + 1
  end

end
