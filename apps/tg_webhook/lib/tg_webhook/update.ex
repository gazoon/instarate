defmodule TGWebhook.Update do
  alias Utils.Messages.Text, as: TextMessage
  alias Utils.Messages.Callback, as: CallbackMessage
  require Logger

  @config Application.get_env(:tg_webhook, __MODULE__)
  @queue @config[:queue]

  def process(queue_name, update) do
    Logger.info("Receive update #{inspect update}")
    bot_message = cond do
      update["message"] != nil && update["message"]["text"] != nil ->
        message = update["message"]
        reply_to = message["reply_to_message"]
        reply_to_data = if reply_to, do: convert_message_to_data(reply_to), else: nil
        message_data = convert_message_to_data(message)
        message_data = Map.put(message_data, :reply_to, reply_to_data)
        %{
          type: TextMessage.type,
          data: message_data
        }
      update["callback_query"] != nil && update["callback_query"]["message"]["chat"] ->
        callback = update["callback_query"]
        %{
          type: CallbackMessage.type,
          data: %{
            callback_id: callback["id"],
            user: convert_user_to_data(callback["from"]),
            parent_msg_id: callback["message"]["message_id"],
            payload: callback["data"],
            chat_id: callback["message"]["chat"]["id"],
            is_group_chat: callback["message"]["chat"]["type"] != "private",
          }
        }
      true -> nil
    end
    if bot_message do
      chat_id = bot_message.data.chat_id
      Logger.info("Put message chat=#{chat_id} queue=#{queue_name} msg=#{inspect bot_message}")
      @queue.put(chat_id, bot_message, queue_name: queue_name)
      update
    else
      update
    end
  end

  defp concatenate_name(user_from) do
    user_from["first_name"] <> (user_from["last_name"] || "")
  end

  defp convert_user_to_data(user_from) do
    %{
      id: user_from["id"],
      name: concatenate_name(user_from),
      username: user_from["username"] || ""
    }
  end

  defp convert_message_to_data(message) do
    %{
      user: convert_user_to_data(message["from"]),
      text: message["text"] || "",
      chat_id: message["chat"]["id"],
      is_group_chat: message["chat"]["type"] != "private"
    }
  end
end
