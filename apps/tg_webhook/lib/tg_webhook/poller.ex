defmodule TGWebhook.Poller do

  #  use GenServer
  #  require Logger
  #  alias TGWebhook.Update
  #
  #  def start_link(opts) do
  #    GenServer.start_link(__MODULE__, 0, opts)
  #  end
  #
  #  def handle_info(_msg, offset) do
  #    Logger.info("get updates")
  #    new_offset = Nadia.get_updates([offset: offset, timeout: 60])
  #                 |> process_updates
  #    next_cast()
  #    if is_integer(new_offset) do
  #      {:noreply, new_offset + 1}
  #    else
  #      {:noreply, offset + 1}
  #    end
  #  end
  #
  #  defp process_updates({:ok, results}) do
  #    results
  #    |> Enum.map(
  #         fn %{update_id: id} = update ->
  #           update
  #           |> Update.process()
  #           id
  #         end
  #       )
  #    |> List.last
  #  end
  #
  #  defp next_cast do
  #    send(self(), :next)
  #  end
  #
  #  def init(state) do
  #    Process.flag(:trap_exit, true)
  #    next_cast()
  #    {:ok, state}
  #  end

end
