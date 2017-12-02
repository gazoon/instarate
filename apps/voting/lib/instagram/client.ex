defmodule Instagram.Client do
  @moduledoc false
  alias Instagram.Media

  @callback parse_username(profile_uri :: String.t) :: String.t
  @callback parse_media_code(media_uri :: String.t) :: String.t
  @callback get_media_info(media_code :: String.t) :: {:ok, Media.t} | {:error, String.t}
  @callback build_profile_url(username :: String.t) :: String.t
end
