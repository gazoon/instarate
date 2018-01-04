defmodule TGBot.Pictures.Concatenators.HttpService do
  require Logger
  alias TGBot.Pictures.Concatenator
  use Utils.Meter

  @version "v1"
  @behaviour Concatenator
  @spec version :: String.t
  def version, do: @version
  @service_url Application.get_env(:tg_bot, :concatenation_service)

  @spec concatenate(String.t, String.t) :: String.t
  def concatenate(left_picture_url, right_picture_url) do
    case HTTPoison.post!(
           @service_url <> "/concatenate",
           {:form, [{"left_picture", left_picture_url}, {"right_picture", right_picture_url}]}
         ) do
      %HTTPoison.Response{body: body, status_code: 200} -> body
      %HTTPoison.Response{body: body, status_code: status_code} ->
        raise "Concatenation request failed #{status_code}: #{body}"
    end
  end

end
