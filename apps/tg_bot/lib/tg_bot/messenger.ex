defmodule TGBot.Messenger do
  @moduledoc false

  @spec send_text(integer, String.t) :: integer
  def send_text(chat_id, text) do
    send_message(chat_id, text)
  end

  @spec send_markdown(integer, String.t) :: integer
  def send_markdown(chat_id, text) do
    send_message(chat_id, text, parse_mode: "Markdown", disable_web_page_preview: true)
  end

  @spec send_photo(integer, binary, [{atom, any}]) :: integer
  def send_photo(chat_id, photo, opts \\ []) do
    case Nadia.send_photo(chat_id, photo, opts) do
      {:error, error} -> raise error
      {:ok, msg} -> msg.message_id
    end
  end

  @spec send_message(integer, String.t, [{atom, any}]) :: integer
  defp send_message(chat_id, text, opts \\ []) do
    case Nadia.send_message(chat_id, text, opts) do
      {:error, error} -> raise error
      {:ok, msg} -> msg.message_id
    end
  end
end
