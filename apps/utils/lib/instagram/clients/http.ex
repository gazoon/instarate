defmodule Instagram.Clients.Http do

  alias Instagram.Media

  @behaviour Instagram.Client

  @api_url "https://www.instagram.com/"
  @media_path "p/"
  @magic_suffix "/?__a=1"
  @http_pool_name :instagram_api

  @spec get_media_info(String.t) :: {:ok, Media.t} | {:error, String.t}
  def get_media_info(media_code) do
    case request_media(media_code) do
      {:ok, media_resp} ->
        media_data = retrieve_media_data(media_resp)
        media_info = %Media{
          owner: retrieve_username_from_data(media_data),
          url: retrieve_display_url_from_data(media_data),
          is_photo: retrieve_is_photo_from_data(media_data)
        }
        {:ok, media_info}
      error -> error
    end
  end

  @spec get_followers_number(String.t) :: integer
  def get_followers_number(username) do
    profile_url = @api_url <> username <> @magic_suffix
    resp = HTTPoison.get!(
      profile_url,
      [],
      hackney: [
        pool: @http_pool_name
      ]
    )
    if resp.status_code == 404 do
      raise "Instagram profile #{username} not found"
    end

    data = case Poison.decode(resp.body, as: %{}) do
      {:ok, data} -> data
      _ -> raise  "Profile got invalid json, url #{profile_url}, response code #{resp.status_code}"
    end
    followers = data["user"]["followed_by"]["count"]
    if is_integer(followers),
       do: followers,
       else: raise "Profile response doesn't contain followers"
  end

  @spec build_profile_url(String.t) :: String.t
  def build_profile_url(username) do
    @api_url <> username <> "/"
  end

  @spec parse_username(String.t) :: String.t
  def parse_username(profile_uri) do
    get_last_path_part(profile_uri)
  end

  @spec parse_media_code(String.t) :: String.t
  def parse_media_code(media_uri) do
    get_last_path_part(media_uri)
  end

  @spec get_last_path_part(String.t) :: String.t
  defp get_last_path_part(url) do
    URI.parse(url).path
    |> Path.split()
    |> List.last()
  end

  @spec retrieve_username_from_data(map()) :: String.t
  defp retrieve_username_from_data(media_data) do
    username = media_data["owner"]["username"]
    if username do
      username
    else
      raise "Media doesn't contain owner info"
    end
  end

  @spec retrieve_display_url_from_data(map()) :: String.t
  defp retrieve_display_url_from_data(media_data) do
    display_url = List.first(media_data["display_resources"])["src"]
    if display_url do
      display_url
    else
      raise "Media doesn't contain display url"
    end
  end

  @spec retrieve_is_photo_from_data(map()) :: boolean
  defp retrieve_is_photo_from_data(media_data) do
    if media_data["is_video"] do
      false
    else
      !List.first(media_data["edge_sidecar_to_children"]["edges"] || [])["node"]["is_video"]
    end
  end

  @spec retrieve_media_data(map()) :: map()
  defp retrieve_media_data(media_response), do: media_response["graphql"]["shortcode_media"]

  @spec request_media(String.t) :: {:ok, map()} | {:error, String.t}
  defp request_media(media_code) do
    media_url = @api_url <> @media_path <> media_code <> @magic_suffix
    resp = HTTPoison.get!(
      media_url,
      [],
      hackney: [
        pool: @http_pool_name
      ]
    )
    if resp.status_code == 404 do
      {:error, "It's a private account or the media doesn't exist"}
    else
      case Poison.decode(resp.body, as: %{}) do
        {:ok, data} -> {:ok, data}
        _ -> raise  "Media got invalid json url #{media_url}, response code #{resp.status_code}"
      end
    end
  end
end
