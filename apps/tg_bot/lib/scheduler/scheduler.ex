defmodule Scheduler.Scheduler do
  alias TGBot.Messages.Task

  @type t :: module
  @callback create_task(task :: Task.t) :: {:ok, Task.t} | {:error, String.t}
end

