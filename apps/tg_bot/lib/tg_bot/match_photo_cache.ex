defmodule TGBot.MatchPhotoCache do
  @config Application.get_env(:tg_bot, __MODULE__)
  alias TGBot.Cache.Behaviour, as: Cache
  @cache @config[:cache]
  @pictures @config[:pictures_concatenator]

  @spec get(String.t, String.t) :: Cache.value | nil
  def get(left_photo, right_photo) do
    key = build_key(left_photo, right_photo)
    @cache.get(key)
  end

  @spec set(String.t, String.t, String.t) :: :ok
  def set(left_photo, right_photo, tg_file_id) do
    key = build_key(left_photo, right_photo)
    @cache.set(key, tg_file_id)
  end

  @spec build_key(String.t, String.t) :: String.t
  defp build_key(left_photo, right_photo) do
    left_photo <> " | " <> right_photo <> @pictures.version()
  end
end
