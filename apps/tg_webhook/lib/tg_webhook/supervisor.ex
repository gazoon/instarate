defmodule TGWebhook.Supervisor do

  use Supervisor

  alias Utils.Queue.Impls.Mongo, as: MongoQueue
  require Logger

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end

  def init(_) do
    port = Application.get_env(:tg_webhook, :port)
    children = [
      MongoQueue.child_spec(),
      {
        Plug.Adapters.Cowboy,
        scheme: :http,
        plug: TGWebhook.Router,
        options: [
          port: port
        ]
      },
    ]

    Logger.info "Start server on #{port} port"

    Supervisor.init(children, strategy: :one_for_one)

  end
end
