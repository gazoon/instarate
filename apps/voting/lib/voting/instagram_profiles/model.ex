defmodule Voting.InstagramProfiles.Model do
  alias Voting.InstagramProfiles.Model, as: Profile
  alias Instagram.Client, as: InstagramClient

  @type t :: %Profile{
               username: String.t,
               photo: String.t,
               photo_code: String.t,
               added_at: integer,
               followers: integer
             }

  defstruct username: nil,
            photo: nil,
            photo_code: nil,
            added_at: nil,
            followers: nil,
            unreachable: nil

  @spec new(String.t, String.t, String.t, integer) :: Profile.t
  def new(username, photo, photo_code, followers) do
    current_time = Utils.timestamp()
    %Profile{
      username: username,
      photo: photo,
      photo_code: photo_code,
      followers: followers,
      added_at: current_time,
      unreachable: false
    }
  end

  @spec get_profile_url(Profile.t) :: String.t
  def get_profile_url(girl) do
    InstagramClient.build_profile_url(girl.username)
  end
end


