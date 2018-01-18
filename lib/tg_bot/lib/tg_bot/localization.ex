defmodule TGBot.Localization do
  @english_lang "en"
  @russian_lang "ru"

  def english_lang, do: @english_lang
  def russian_lang, do: @russian_lang
  @config Application.get_env(:tg_bot, __MODULE__)
  @disable_translation? @config[:disable_translation]
  use Gettext, otp_app: :tg_bot
  alias TGBot.Chats.Chat


  @spec  get_translation(Chat.t, String.t, Keyword.t) :: String.t
  def get_translation(chat, msgid, bindings \\ []) do
    translate_to(chat.language, msgid, bindings)
  end

  @spec get_default_translation(String.t, Keyword.t) :: String.t
  def get_default_translation(msgid, bindings \\ []) do
    translate_to(@english_lang, msgid, bindings)
  end

  @spec  translate_to(String.t, String.t, Keyword.t) :: String.t
  defp translate_to(language, msgid, bindings) do
    if @disable_translation? do
      msgid
    else
      if language == "en" && msgid == "place_in_competition" do
        place = Keyword.get(bindings, :place)
        case rem(place, 10) do
          1 -> "#{place}st place: "
          2 -> "#{place}nd place: "
          3 -> "#{place}rd place: "
          _ -> "#{place}th place: "
        end
      else
        Gettext.with_locale __MODULE__, language, fn ->
          Gettext.dgettext(__MODULE__, "messages", msgid, bindings)
        end
      end
    end
  end
end
