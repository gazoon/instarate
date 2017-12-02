defmodule TGWebhook.Poller do
  @moduledoc false
  use GenServer
  require Logger
  alias TGBot.Messages.Text, as: TextMessage
  alias TGBot.Messages.Callback, as: CallbackMessage

  def start_link(opts) do
    GenServer.start_link(__MODULE__, 0, opts)
  end

  def handle_cast(:update, offset) do
    Logger.info("get updates")
    new_offset = Nadia.get_updates([offset: offset, timeout: 60])
                 |> process_updates
    :timer.sleep(1000)
    next_cast()
    Logger.info("get updates2")
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

  defp next_cast() do
    GenServer.cast(self(), :update)
  end

  defp process_update(update) do

    IO.inspect(update)
    bot_message = cond do
      update.message != nil ->
        message = update.message
        %{
          type: TextMessage.type,
          data: %{
            user_id: message.from.id,
            chat_id: message.chat.id,
            text: message.text,
            is_group_chat: update.message.chat.type != "private"
          }
        }
      update.callback_query != nil && update.callback_query.message.chat ->
        callback = update.callback_query
        %{
          type: CallbackMessage.type,
          data: %{
            callback_id: callback.id,
            user_id: callback.from.id,
            chat_id: callback.message.chat.id,
            is_group_chat: callback.message.chat.type != "private",
            payload: callback.data
          }
        }
      true -> nil
    end
    if bot_message do
      TGBot.on_message(bot_message)
      update
    else
      Logger.error("unsupported message format")
      update
    end
  end


  def init(state) do
    next_cast()
    {:ok, state}
  end

end
