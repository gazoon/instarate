defmodule TGWebhook.Supervisor do

  use Supervisor

  alias TGWebhook.Poller
  alias Utils.Queue.Impls.Mongo, as: MongoQueue

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end

  def init(_) do
    children = [
      MongoQueue.child_spec(),

      Poller,
    ]

    Supervisor.init(children, strategy: :one_for_one)

  end
end
