defmodule Instagram.Client do
  @moduledoc false
  alias Instagram.Media

  @callback get_media_info(media_code :: String.t) :: {:ok, Media.t} | {:error, String.t}
end
