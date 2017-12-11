defmodule Scheduler.Reader do
  alias TGBot.Messages.Task
  @config Application.get_env(:tg_bot, __MODULE__)
  @storage @config[:tasks_storage]
  use Utils.Reader, workers_number: 5

  @spec fetch :: Task.t
  def fetch do
    @storage.get_available_task()
  end

  @spec process(Task.t) :: any
  def process(task) do
    task_data = task_to_raw_data(task)
    message = %{
      type: Task.type,
      data: task_data
    }
    TGBot.on_message(message)
  end

  defp task_to_raw_data(task) do
    Poison.encode!(task)
    |> Poison.decode!(as: %{})
  end

end
