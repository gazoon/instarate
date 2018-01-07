defmodule Scheduler.Scheduler do
  alias Utils.Messages.Task

  @type t :: module
  @callback create_task(task :: Task.t) :: {:ok, Task.t} | {:error, String.t}
  @callback create_or_replace_task(task :: Task.t) :: Task.t
  @callback delete_task(chat_id :: integer, name :: atom) :: :ok
end

