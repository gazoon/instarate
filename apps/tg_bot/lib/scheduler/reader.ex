defmodule Scheduler.Reader do
  alias TGBot.Messages.Task
  require Logger
  @config Application.get_env(:tg_bot, __MODULE__)
  @storage @config[:tasks_storage]
  @fetch_delay 100
  @workers_number 1
  use GenServer

  def init(state) do
    Process.flag(:trap_exit, true)
    next_fetch()
    {:ok, state}
  end

  defp next_fetch do
    send(self(), :fetch)
  end

  def start_link(opts) do
    GenServer.start_link(__MODULE__, nil, opts)
  end

  def handle_info(_msg, state) do
    data = fetch()
    if data do
      Logger.info("Fetched #{inspect data} start processing")
      process(data)
    else
      Process.sleep(@fetch_delay)
    end
    next_fetch()
    {:noreply, state}
  end

  defp fetch do
    @storage.get_available_task()
  end

  defp task_to_raw_data(task) do
    Poison.encode!(task)
    |> Poison.decode!(as: %{})
  end

  defp process(task) do
    task_data = task_to_raw_data(task)
    message = %{
      type: Task.type,
      data: task_data
    }
    TGBot.on_message(message)
  end

  def children_spec do
    0..@workers_number - 1
    |> Enum.map(
         fn (i) ->
           process_name = String.to_atom("#{__MODULE__}.#{i}")
           Supervisor.child_spec({__MODULE__, name: process_name}, id: {__MODULE__, i})
         end
       )
  end
end
