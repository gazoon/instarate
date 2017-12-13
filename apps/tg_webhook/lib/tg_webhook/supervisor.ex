defmodule TGWebhook.Supervisor do

  use Supervisor

  alias TGWebhook.Poller

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end

  def init(_) do
    children = [
      Poller,
    ]

    Supervisor.init(children, strategy: :one_for_one)

  end
end
