defmodule Voting.InstagramProfiles.Model do
  alias Voting.InstagramProfiles.Model, as: Profile
  alias Instagram.Client, as: InstagramClient

  @type t :: %Profile{
               username: String.t,
               photo: String.t,
               added_at: integer,
               followers: integer
             }

  defstruct username: nil,
            photo: nil,
            added_at: nil,
            followers: nil

  @spec new(String.t, String.t, integer) :: Profile.t
  def new(username, photo, followers) do
    current_time = Utils.timestamp()
    %Profile{username: username, photo: photo, followers: followers, added_at: current_time}
  end

  @spec get_profile_url(Profile.t) :: String.t
  def get_profile_url(girl) do
    InstagramClient.build_profile_url(girl.username)
  end
end


