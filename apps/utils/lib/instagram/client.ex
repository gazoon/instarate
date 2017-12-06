defmodule Instagram.Client do

  alias Instagram.Media

  @type t :: module
  @client Application.get_env(:utils, __MODULE__)

  @callback parse_username(profile_uri :: String.t) :: String.t
  @callback parse_media_code(media_uri :: String.t) :: String.t
  @callback get_media_info(media_code :: String.t) :: {:ok, Media.t} | {:error, String.t}
  @callback build_profile_url(username :: String.t) :: String.t

  @spec parse_username(String.t) :: String.t
  def parse_username(profile_uri) do
    @client.parse_username(profile_uri)
  end

  @spec parse_media_code(String.t) :: String.t
  def parse_media_code(media_uri) do
    @client.parse_media_code(media_uri)
  end

  @spec get_media_info(String.t) :: {:ok, Media.t} | {:error, String.t}
  def get_media_info(media_code) do
    @client.get_media_info(media_code)
  end

  @spec build_profile_url(String.t) :: String.t
  def build_profile_url(username) do
    @client.build_profile_url(username)
  end
end
