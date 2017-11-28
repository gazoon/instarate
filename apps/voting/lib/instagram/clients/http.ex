defmodule Instagram.Clients.Http do
  @moduledoc false
  alias Instagram.Media

  @behaviour Instagram.Client

  @api_url "https://www.instagram.com/p/"
  @magic_suffix "/?__a=1"

  @spec get_media_info(String.t) :: {:ok, Media.t} | {:error, String.t}
  def get_media_info(media_code) do
    case request_media(media_code) do
      {:ok, media_resp} ->
        media_data = retrieve_media_data(media_resp)
        media_info = %Media{
          owner: retrieve_username(media_data),
          url: retrieve_display_url(media_data),
          is_photo: retrieve_is_photo(media_data)
        }
        {:ok, media_info}
      error -> error
    end
  end

  @spec retrieve_username(map()) :: String.t
  defp retrieve_username(media_data) do
    username = media_data["owner"]["username"]
    if username do
      username
    else
      raise "Media doesn't contain owner info"
    end
  end

  @spec retrieve_display_url(map()) :: String.t
  defp retrieve_display_url(media_data) do
    display_url = media_data["display_url"]
    if display_url do
      display_url
    else
      raise "Media doesn't contain display url"
    end
  end

  @spec retrieve_is_photo(map()) :: boolean
  defp retrieve_is_photo(media_data) do
    !media_data["is_video"]
  end

  @spec retrieve_media_data(map()) :: map()
  defp retrieve_media_data(media_response), do: media_response["graphql"]["shortcode_media"]

  @spec request_media(String.t) :: {:ok, map()} | {:error, String.t}
  defp request_media(media_code) do
    media_url = @api_url <> media_code <> @magic_suffix
    resp = HTTPoison.get!(media_url)
    if resp.status_code == 404 do
      {:error, "Media #{media_code} not found"}
    else
      case Poison.decode(resp.body, as: %{}) do
        {:ok, data} -> {:ok, data}
        _ -> raise  "Got invalid json url #{media_url}, response code #{resp.status_code}"
      end
    end
  end
end
