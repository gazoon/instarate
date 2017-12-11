defmodule Scheduler.Impls.Mongo do
  alias Scheduler.{Storage, Scheduler}
  @behaviour Storage
  @behaviour Scheduler
  @process_name :mongo_scheduler
  @duplication_code 11_000
  @collection "tasks"
  alias TGBot.Messages.Task

  @spec child_spec :: tuple
  def child_spec do
    options = [name: @process_name] ++ Application.get_env(:tg_bot, :mongo_scheduler)
    Utils.set_child_id(Mongo.child_spec(options), {Mongo, :scheduler})
  end

  @spec create_task(Task.t) :: {:ok, Task.t} | {:error, String.t}
  def create_task(task) do
    task_data = if task.unique_mark, do: task, else: Map.delete(task, :unique_mark)
    insert_result = Mongo.insert_one(@process_name, @collection, task_data)
    case insert_result do
      {:ok, _} ->
        {:ok, task}
      {:error, %Mongo.Error{code: @duplication_code}} ->
        {:error, "Task with mark #{task.unique_mark} for chat #{task.chat_id} already exists"}
      {:error, error} -> raise error
    end

  end

  @spec get_available_task :: Task.t
  def get_available_task do
    current_time = Utils.timestamp_milliseconds()
    case Mongo.find_one_and_delete(
           @process_name,
           @collection,
           %{
             do_at: %{
               "$lte": current_time
             }
           },
           sort: %{
             do_at: 1
           }
         ) do
      {:ok, row} -> transform_task(row)
      {:error, error} -> raise  error
    end
  end

  defp transform_task(nil), do: nil
  @spec transform_task(map()) :: Task.t
  defp transform_task(row) do
    Task.from_data(row)
  end
end
