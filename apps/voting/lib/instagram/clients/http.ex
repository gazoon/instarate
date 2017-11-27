defmodule Instagram.Clients.Http do
  @moduledoc false
  @behaviour Instagram.Client

  @api_url "https://www.instagram.com/p/"
  @magic_suffix "?__a=1"

  @spec get_media_owner(String.t) :: {:ok, String.t} | {:error, String.t}
  def get_media_owner(media_code) do
    case request_media(media_code) do
      {:ok, media_resp} -> {:ok, retrieve_username(media_resp)}
      error -> error
    end
  end

  @spec is_photo?(String.t) :: boolean
  def is_photo?(media_code) do
    case request_media(media_code) do
      {:ok, media_resp} -> !retrieve_media_data(media_resp)["is_video"]
      _ -> false
    end
  end

  @spec build_media_url(String.t) :: String.t
  def build_media_url(media_code) do
    @api_url <> media_code <> "/"
  end

  @spec retrieve_username(map()) :: String.t
  defp retrieve_username(media_resp) do
    username = retrieve_media_data(media_resp)["owner"]["username"]
    if username do
      username
    else
      raise "Media doesn't contain owner info"
    end
  end

  @spec retrieve_media_data(map()) :: map()
  defp retrieve_media_data(media_response), do: media_response["graphql"]["shortcode_media"]

  @spec request_media(String.t) :: {:ok, map()} | {:error, String.t}
  defp request_media(media_code) do
    media_url = build_media_url(media_code) <> @magic_suffix
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
