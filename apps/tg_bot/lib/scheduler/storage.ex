defmodule Scheduler.Storage do
  alias TGBot.Messages.Task
  @type t :: module
  @callback get_available_task() :: Task.t | nil
end
