defmodule Scheduler.Storage do
  alias Utils.Messages.Task
  @type t :: module
  @callback get_available_task() :: Task.t | nil
end
