defmodule Voting.InstagramProfiles.Storage do

  alias Voting.InstagramProfiles.Model, as: Profile
  @callback add(profile :: Profile) :: {:ok, Profile.t | :error, String.t}
  @callback get(username :: String.t) :: Profile.t
  @callback get_multiple(usernames :: [String.t]) :: [Profile.t]
end
