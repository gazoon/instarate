defmodule TGBot.Messenger do
  @moduledoc false
  alias Nadia.Model.InlineKeyboardButton
  alias Nadia.Model.InlineKeyboardMarkup

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

  @spec send_photo(integer, binary, Keyword.t) :: integer
  def send_photo(chat_id, photo, opts \\ []) do
    opts = transform_opts(opts)
    case Nadia.send_photo(chat_id, photo, opts) do
      {:error, error} -> raise error
      {:ok, msg} -> msg.message_id
    end
  end

  @spec transform_opts(Keyword.t) :: Keyword.t
  defp transform_opts(opts) do
    {keyboard, opts} = Keyword.pop(opts, :keyboard)
    reply_markup = transform_keyboard(keyboard)
    if reply_markup, do: Keyword.put(opts, :reply_markup, reply_markup), else: opts
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
              callback_data: Poison.encode!(item_data.payload),
              url: "",
              switch_inline_query: ""
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
