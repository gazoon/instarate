defmodule Voting.Girls.Girl do

  alias Voting.Girls.Girl
  alias Instagram.Client, as: InstagramClient
  @initial_rating 1500
  @storage Application.get_env(:voting, __MODULE__)[:storage]

  @type t :: %Girl{
               username: String.t,
               photo: String.t,
               added_at: integer,
               rating: integer,
               matches: integer,
               wins: integer,
               loses: integer
             }

  defstruct username: nil,
            photo: nil,
            added_at: nil,
            rating: @initial_rating,
            matches: 0,
            wins: 0,
            loses: 0

  @spec new(String.t, String.t) :: Girl.t
  def new(username, photo) do
    current_time = Utils.current_timestamp()
    %Girl{username: username, photo: photo, added_at: current_time}
  end

  @spec get_profile_url(Girl.t) :: String.t
  def get_profile_url(girl) do
    InstagramClient.build_profile_url(girl.username)
  end

  @spec get_position(Girl.t) :: integer
  def get_position(girl) do
    @storage.get_higher_ratings_number(girl.rating) + 1
  end

end
