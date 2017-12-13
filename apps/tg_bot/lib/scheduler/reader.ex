defmodule Scheduler.Reader do
  alias TGBot.Messages.Task

  @config Application.get_env(:tg_bot, __MODULE__)
  @storage @config[:tasks_storage]
  @queue @config[:queue]

  use Utils.Reader, otp_app: :tg_bot

  @spec fetch :: Task.t
  def fetch do
    @storage.get_available_task()
  end

  @spec process(Task.t) :: any
  def process(task) do
    Logger.info("Send task to the queue #{inspect task}")
    message = %{type: Task.type, data: Map.from_struct(task)}
    @queue.put(task.chat_id, message)
  end

end
