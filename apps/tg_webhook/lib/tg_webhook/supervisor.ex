defmodule TGWebhook.Supervisor do
  @moduledoc false


  use Supervisor
  alias TGWebhook.Poller

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end


  def init(_) do
    children = [
      Poller,
      {Task.Supervisor, name: :messages_supervisor}
    ]

    Supervisor.init(children, strategy: :one_for_one)

  end
end