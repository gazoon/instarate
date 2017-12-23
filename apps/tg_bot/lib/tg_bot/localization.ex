defmodule TGBot.Localization do
  @english_lang "en"
  @russian_lang "ru"

  def english_lang, do: @english_lang
  def russian_lang, do: @russian_lang
  @config Application.get_env(:tg_bot, __MODULE__)
  @disable_translation? @config[:disable_translation]
  use Gettext, otp_app: :tg_bot
  alias TGBot.Chats.Chat


  @spec  get_translation(Chat.t, String.t, map()) :: String.t
  def get_translation(chat, msgid, bindings \\ %{}) do
    if @disable_translation? do
      msgid
    else
      Gettext.with_locale __MODULE__, chat.language, fn ->
        Gettext.dgettext(__MODULE__, "messages", msgid, bindings)
      end
    end
  end
end
