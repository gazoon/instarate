defmodule TGBot.Messengers.NadiaLib do
  @behaviour TGBot.Messenger

  alias Nadia.Model.{InlineKeyboardButton, ReplyKeyboardHide, InlineKeyboardMarkup}
  alias Nadia.Model.ReplyKeyboardMarkup
  alias Nadia.Model.KeyboardButton, as: ReplyKeyboardButton

  @spec send_text(integer, String.t, Keyword.t) :: integer
  def send_text(chat_id, text, opts \\ []) do
    opts = transform_opts(opts)
    send_message(chat_id, text, opts)
  end

  @spec send_markdown(integer, String.t, Keyword.t) :: integer
  def send_markdown(chat_id, text, opts \\ []) do
    opts = transform_opts(opts)
    send_message(
      chat_id,
      text,
      opts ++
      [
        parse_mode: "Markdown",
        disable_web_page_preview: true
      ]
    )
  end

  @spec get_chat_members_number(integer) :: integer
  def get_chat_members_number(chat_id) do
    case Nadia.get_chat_members_count(chat_id) do
      {:ok, number} -> number
      {:error, error} -> raise error
    end
  end

  @spec send_photo(integer, binary, Keyword.t) :: integer
  def send_photo(chat_id, photo, opts \\ []) do
    opts = transform_opts(opts)
    case Nadia.send_photo(chat_id, photo, opts) do
      {:error, error} -> raise error
      {:ok, msg} -> msg.message_id
    end
  end

  @spec send_notification(String.t, String.t) :: any
  def send_notification(callback_id, text) do
    answer_callback(callback_id, text: text)
  end

  @spec answer_callback(String.t, Keyword.t) :: any
  def answer_callback(callback_id, opts \\ []) do
    case Nadia.answer_callback_query(callback_id, opts) do
      {:error, error} -> raise error
      _ -> nil
    end
  end

  @spec delete_attached_keyboard(integer, integer) :: integer
  def delete_attached_keyboard(chat_id, message_id) do
    case Nadia.API.request("editMessageReplyMarkup", chat_id: chat_id, message_id: message_id) do
      {:error, error} -> raise error
      {:ok, msg} -> msg.message_id
    end
  end

  @spec send_photos(integer, [map()]) :: any
  def send_photos(chat_id, photos) do
    media = Enum.map(
      photos,
      fn (photo) -> %{type: "photo", media: photo.url, caption: photo.caption} end
    )
    media_encoded = Poison.encode!(media)
    try do
      case Nadia.API.request("sendMediaGroup", chat_id: chat_id, media: media_encoded) do
        {:error, error} -> raise error
        _ -> nil
      end
    rescue
      FunctionClauseError -> nil
    end
  end

  @spec transform_opts(Keyword.t) :: Keyword.t
  defp transform_opts(opts) do
    {keyboard, opts} = Keyword.pop(opts, :keyboard)
    inline_markup = transform_keyboard(keyboard)
    {static_keyboard, opts} = Keyword.pop(opts, :static_keyboard)
    reply_markup = transform_static_keyboard(static_keyboard)
    cond do
      inline_markup -> Keyword.put(opts, :reply_markup, inline_markup)
      reply_markup ->
        reply_markup = if Keyword.get(opts, :one_time_keyboard, false),
                          do: %ReplyKeyboardMarkup{reply_markup | one_time_keyboard: true},
                          else: reply_markup
        Keyword.put(opts, :reply_markup, reply_markup)
      true -> opts
    end
  end

  defp transform_static_keyboard(nil), do: nil
  defp transform_static_keyboard(:remove), do: %ReplyKeyboardHide{}
  @spec transform_static_keyboard([[map()]]) :: ReplyKeyboardMarkup.t
  defp transform_static_keyboard(keyboard_data) do
    keyboard = Enum.map(
      keyboard_data,
      fn keyboard_line ->
        Enum.map(keyboard_line, fn key_text -> %ReplyKeyboardButton{text: key_text} end)
      end
    )
    %ReplyKeyboardMarkup{keyboard: keyboard, resize_keyboard: true}
  end

  defp transform_keyboard(nil), do: nil
  @spec transform_keyboard([[map()]]) :: InlineKeyboardMarkup.t
  defp transform_keyboard(keyboard_data) do
    keyboard = Enum.map(
      keyboard_data,
      fn keyboard_line ->
        Enum.map(
          keyboard_line,
          fn item_data ->
            %InlineKeyboardButton{
              text: item_data.text,
              callback_data: item_data.payload
            }
          end
        )
      end
    )
    %InlineKeyboardMarkup{inline_keyboard: keyboard}
  end

  @spec send_message(integer, String.t, Keyword.t) :: integer
  defp send_message(chat_id, text, opts \\ []) do
    case Nadia.send_message(chat_id, text, opts) do
      {:error, error} -> raise error
      {:ok, msg} -> msg.message_id
    end
  end
end

