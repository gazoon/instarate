defmodule Scheduler.Impls.Mongo do
  alias Scheduler.{Storage, Scheduler}
  @behaviour Storage
  @behaviour Scheduler
  @process_name :mongo_scheduler
  @duplication_code 11_000
  @collection "tasks"
  alias Utils.Messages.Task

  @spec child_spec :: tuple
  def child_spec do
    options = [name: @process_name, pool: DBConnection.Poolboy] ++
              Application.get_env(:tg_bot, :mongo_scheduler)
    Utils.set_child_id(Mongo.child_spec(options), {Mongo, :scheduler})
  end

  @spec create_task(Task.t) :: {:ok, Task.t} | {:error, String.t}
  def create_task(task) do
    insert_result = Mongo.insert_one(@process_name, @collection, task, pool: DBConnection.Poolboy)
    case insert_result do
      {:ok, _} ->
        {:ok, task}
      {:error, %Mongo.Error{code: @duplication_code}} ->
        {:error, "Task #{task.name} for chat #{task.chat_id} already exists"}
      {:error, error} -> raise error
    end
  end

  @spec create_or_replace_task(Task.t) :: Task.t
  def create_or_replace_task(task) do
    task_data = Map.from_struct(task)
    Mongo.replace_one!(
      @process_name,
      @collection,
      %{chat_id: task.chat_id, name: task.name},
      task_data,
      upsert: true,
      pool: DBConnection.Poolboy
    )
    task
  end

  @spec delete_task(integer, atom) :: :ok
  def delete_task(chat_id, name) do
    Mongo.delete_one!(
      @process_name,
      @collection,
      %{chat_id: chat_id, name: name},
      pool: DBConnection.Poolboy
    )
    :ok
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
           },
           pool: DBConnection.Poolboy
         ) do
      {:ok, row} -> transform_task(row)
      {:error, error} -> raise  error
    end
  end

  @spec transform_task(map()) :: Task.t
  defp transform_task(row)

  defp transform_task(nil), do: nil

  defp transform_task(row) do
    Task.from_data(row)
  end
end
