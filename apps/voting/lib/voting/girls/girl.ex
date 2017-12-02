defmodule Voting.Girls.Girl do
  @moduledoc false
  alias Voting.Girls.Girl
  @initial_rating 1500
  @instagram_client Application.get_env(:voting, :instagram_client)
  @storage Application.get_env(:voting, :girls_storage)

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
    current_time = DateTime.utc_now()
                   |> DateTime.to_unix()
    %Girl{username: username, photo: photo, added_at: current_time}
  end

  @spec get_profile_url(Girl.t) :: String.t
  def get_profile_url(girl) do
    @instagram_client.build_profile_url(girl.username)
  end

  @spec get_position(Girl.t) :: integer
  def get_position(girl) do
    @storage.get_higher_ratings_number(girl.rating) + 1
  end

end
