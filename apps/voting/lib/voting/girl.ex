defmodule Voting.Girl do

  alias Voting.Girl
  alias Voting.Competitors.Model, as: Competitor
  alias Voting.InstagramProfiles.Model, as: Profile

  @type t :: %Girl{
               username: String.t,
               photo: String.t,
               photo_code: String.t,
               added_at: integer,
               followers: integer,
               competition: String.t,
               rating: integer,
               matches: integer,
               wins: integer,
               loses: integer
             }

  defstruct username: nil,
            photo: nil,
            photo_code: nil,
            added_at: nil,
            followers: nil,
            competition: nil,
            rating: nil,
            matches: 0,
            wins: 0,
            loses: 0


  @spec combine(Competitor.t, Profile.t) :: Girl.t
  def combine(competitor, profile) do
    %Girl{
      username: profile.username,
      photo: profile.photo,
      photo_code: profile.photo_code,
      followers: profile.followers,
      added_at: profile.added_at,
      competition: competitor.competition,
      rating: competitor.rating,
      matches: competitor.matches,
      wins: competitor.wins,
      loses: competitor.loses,
    }
  end

  @spec get_profile_url(Girl.t) :: String.t
  def get_profile_url(girl) do
    girl
    |> to_profile
    |> Profile.get_profile_url()
  end

  @spec get_position(Girl.t) :: integer
  def get_position(girl) do
    girl
    |> to_competitor
    |> Profile.get_profile_url()
    Competitor.get_position(girl)
  end

  @spec get_profile_url(Girl.t) :: String.t
  def get_photo_url(girl) do
    girl
    |> to_profile
    |> Profile.get_photo_url()
  end

  @spec to_profile(Girl.t) :: Profile.t
  def to_profile(girl) do
    %Profile{
      username: girl.username,
      photo: girl.photo,
      photo_code: girl.photo_code,
      followers: girl.followers,
      added_at: girl.added_at
    }
  end

  @spec to_competitor(Girl.t) :: Competitor.t
  def to_competitor(girl) do
    %Competitor{
      username: girl.username,
      competition: girl.competition,
      rating: girl.rating,
      matches: girl.matches,
      wins: girl.wins,
      loses: girl.loses
    }
  end
end
