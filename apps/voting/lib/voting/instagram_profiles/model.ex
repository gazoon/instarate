defmodule Voting.InstagramProfiles.Model do
  alias Voting.InstagramProfiles.Model, as: Profile
  alias Instagram.Client, as: InstagramClient
  @config Application.get_env(:voting, __MODULE__)
  @photos_storage @config[:photos_storage]

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
            followers: nil

  @spec new(String.t, String.t, String.t, integer) :: Profile.t
  def new(username, photo, photo_code, followers) do
    current_time = Utils.timestamp()
    %Profile{
      username: username,
      photo: photo,
      photo_code: photo_code,
      followers: followers,
      added_at: current_time
    }
  end

  @spec upload_photo(Profile.t, String.t) :: String.t
  def upload_photo(girl, original_photo_url) do
    @photos_storage.upload(girl.photo, original_photo_url)
  end

  @spec get_profile_url(Profile.t) :: String.t
  def get_profile_url(girl) do
    InstagramClient.build_profile_url(girl.username)
  end

  @spec get_photo_url(Profile.t) :: String.t
  def get_photo_url(girl) do
    @photos_storage.build_url(girl.photo)
  end
end


