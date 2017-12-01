defmodule TGBot do
  @moduledoc false
  alias TGBot.Messages.Text, as: TextMessage
  alias TGBot.Messages.Callback, as: CallbackMessage

  def on_message(message) do
    case message.type do
      :text -> on_text_message(message.data)
      :callback -> on_callback(message.data)
    end
  end
  def on_text_message(message) do
    reply_markup = %Nadia.Model.InlineKeyboardMarkup{
      inline_keyboard: [
        [
          %Nadia.Model.InlineKeyboardButton{text: "first", callback_data: "first data", url: ""},
          %Nadia.Model.InlineKeyboardButton{text: "second", callback_data: "second data", url: ""},
        ],
        [
          %Nadia.Model.InlineKeyboardButton{text: "third", callback_data: "third data", url: ""},
          %Nadia.Model.InlineKeyboardButton{text: "fourth", callback_data: "fourth data", url: ""},
        ]
      ]
    }
    x = Nadia.send_photo(
      message.chat_id,
      "https://scontent-arn2-1.cdninstagram.com/t51.2885-15/e35/23498332_181024335786384_452312429100007424_n.jpg",
      reply_markup: reply_markup
    )
  end

  def on_callback(message) do
    IO.inspect("caaaaaaaaaaaaalllllback")
    IO.inspect(message)
  end

end
