defmodule Instagram.Client do
  @moduledoc false

  @callback get_media_owner(media_code :: String.t) :: {:ok, String.t} | {:error, String.t}

  @callback is_photo?(media_code :: String.t) :: boolean
end
