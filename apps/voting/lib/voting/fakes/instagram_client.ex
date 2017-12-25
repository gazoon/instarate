defmodule Voting.Fakes.InstagramClient do
  @behaviour Instagram.Client
  alias Instagram.Media

  @spec get_media_info(String.t) :: {:ok, Media.t} | {:error, String.t}
  def get_media_info(media_code) do
    case media_code do
      "not_photo_media" ->
        {:ok, %Media{is_photo: false}}
      "non_existent_media" ->
        {:error, "not found"}
      media_code ->
        {:ok, %Media{owner: "#{media_code}_owner", url: "#{media_code}_url", is_photo: true}}
    end
  end

  @spec get_followers_number(String.t) :: integer
  def get_followers_number(username) do
    cond do
      String.starts_with?(username, "celebrity") -> 10_000_000
      String.starts_with?(username, "model") -> 100_000
      true -> 100
    end
  end

  @spec parse_media_code(String.t) :: String.t
  def parse_media_code(media_uri), do: media_uri

  # not used
  @spec build_profile_url(String.t) :: String.t
  def build_profile_url(_username), do: ""

  # not used
  @spec parse_username(String.t) :: String.t
  def parse_username(_profile_uri), do: ""
end
