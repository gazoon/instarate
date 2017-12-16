defmodule Scheduler.Scheduler do
  alias TGBot.Messages.Task

  @type t :: module
  @callback create_task(task :: Task.t) :: {:ok, Task.t} | {:error, String.t}
  @callback create_or_replace_task(task :: Task.t) :: Task.t
  @callback delete_task(chat_id :: integer, name :: String.t) :: :ok
end

