defmodule TGBot.MatchPhotoCache do
  alias TGBot.Pictures
  @config Application.get_env(:tg_bot, __MODULE__)
  @cache @config[:cache]

  def get(left_photo, right_photo) do
    key = build_key(left_photo, right_photo)
    @cache.get(key)
  end

  def set(left_photo, right_photo, tg_file_id) do
    key = build_key(left_photo, right_photo)
    @cache.set(key, tg_file_id)
  end

  defp build_key(left_photo, right_photo) do
    left_photo <> " | " <> right_photo <> Pictures.version()
  end
end
