defmodule Voting.Supervisor do

  alias Voting.Competitors.Storages.Mongo, as: CompetitorsMongoStorage
  alias Voting.InstagramProfiles.Storages.Mongo, as: ProfilesMongoStorage
  alias Voting.Voters.Storages.Mongo, as: VotersMongoStorage

  use Supervisor

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end

  def init(_) do
    children = [
      {
        Plug.Adapters.Cowboy,
        scheme: :http,
        plug: Voting.Router,
        options: [
          port: 8080
        ]
      },
      CompetitorsMongoStorage.child_spec(),
      ProfilesMongoStorage.child_spec(),
      VotersMongoStorage.child_spec(),
    ]

    Supervisor.init(children, strategy: :one_for_one)

  end
end
