defmodule TGWebhook.Poller do

  use GenServer, shutdown: 10_000
  require Logger
  alias TGBot.Messages.Text, as: TextMessage
  alias TGBot.Messages.Callback, as: CallbackMessage

  def start_link(opts) do
    GenServer.start_link(__MODULE__, 0, opts)
  end

  def handle_info(_msg, offset) do
    Logger.info("get updates")
    new_offset = Nadia.get_updates([offset: offset, timeout: 60])
                 |> process_updates
    next_cast()
    if is_integer(new_offset) do
      {:noreply, new_offset + 1}
    else
      {:noreply, offset + 1}
    end
  end

  defp process_updates({:ok, results}) do
    results
    |> Enum.map(
         fn %{update_id: id} = update ->
           update
           |> process_update
           id
         end
       )
    |> List.last
  end

  defp next_cast do
    send(self(), :next)
  end

  defp concatenate_name(user_from) do
    user_from.first_name <> (user_from.last_name || "")
  end

  defp convert_user_to_data(user_from) do
    %{
      id: user_from.id,
      name: concatenate_name(user_from),
      username: user_from.username || ""
    }
  end

  defp convert_message_to_data(message) do
    %{
      user: convert_user_to_data(message.from),
      text: message.text || "",
      chat_id: message.chat.id,
      is_group_chat: message.chat.type != "private"
    }
  end

  defp process_update(update) do

    Logger.info("Receive update #{inspect update}")
    bot_message = cond do
      update.message != nil && update.message.text != nil ->
        message = update.message
        reply_to = message.reply_to_message
        reply_to_data = if reply_to, do: convert_message_to_data(reply_to), else: nil
        message_data = convert_message_to_data(message)
        message_data = Map.put(message_data, :reply_to, reply_to_data)
        %{
          type: TextMessage.type,
          data: message_data
        }
      update.callback_query != nil && update.callback_query.message.chat ->
        callback = update.callback_query
        %{
          type: CallbackMessage.type,
          data: %{
            callback_id: callback.id,
            user: convert_user_to_data(callback.from),
            parent_msg_id: callback.message.message_id,
            payload: callback.data,
            chat_id: callback.message.chat.id,
            is_group_chat: callback.message.chat.type != "private",
          }
        }
      true -> nil
    end
    if bot_message do
      Task.Supervisor.start_child(
        :messages_supervisor,
        fn -> TGBot.on_message(bot_message) end
      )
      update
    else
      Logger.error("unsupported message format")
      update
    end
  end

  def init(state) do
    Process.flag(:trap_exit, true)
    next_cast()
    {:ok, state}
  end

end
