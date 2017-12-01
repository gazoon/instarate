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
    {:noreply, new_offset + 1, 100}
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
        %TGBot.Message{
          type: TextMessage.type,
          data: %TextMessage{
            user_id: message.from.id,
            chat_id: message.chat.id,
            text: message.text
          }
        }
      update.callback_query != nil && update.callback_query.message.chat ->
        callback = update.callback_query
        %TGBot.Message{
          type: CallbackMessage.type,
          data: %CallbackMessage{
            callback_id: callback.id,
            user_id: callback.from.id,
            chat_id: callback.message.chat.id,
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
