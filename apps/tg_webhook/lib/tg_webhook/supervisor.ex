defmodule TGWebhook.Supervisor do
  @moduledoc false


  use Supervisor
  alias TGWebhook.Poller

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end


  def init(_) do
    children = [
      Poller
      #      {
      #        Plug.Adapters.Cowboy,
      #        scheme: :http,
      #        plug: TGWebhook.Router,
      #        options: [
      #          port: 8080
      #        ]
      #      },
    ]

    Supervisor.init(children, strategy: :one_for_one)

  end
end