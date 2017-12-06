defmodule Voting.Supervisor do

  alias Voting.Girls.Storages.Mongo, as: GirlsMongoStorage
  alias Voting.Voters.Storages.Mongo, as: VotersMongoStorage

  use Supervisor

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end

  def init(_) do
    children = [
      # TODO: add pool
      {
        Plug.Adapters.Cowboy,
        scheme: :http,
        plug: Voting.Router,
        options: [
          port: 8080
        ]
      },
      GirlsMongoStorage.child_spec(),
      VotersMongoStorage.child_spec(),
    ]

    Supervisor.init(children, strategy: :one_for_one)

  end
end
