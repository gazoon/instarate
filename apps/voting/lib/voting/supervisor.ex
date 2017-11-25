defmodule Voting.Supervisor do
  @moduledoc false

  alias Voting.Girls.Storages.Mongo, as: GirlsMongoStorage
  alias Voting.Voters.Storages.Mongo, as: VotersMongoStorage

  use Supervisor

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end
  def set_child_id(spec, child_id) do
    spec
    |> Tuple.delete_at(0)
    |> Tuple.insert_at(0, child_id)
  end


  def init(_) do
    mongo_girls_options = [name: GirlsMongoStorage.process_name] ++
                          Application.get_env(:voting, :mongo_girls)
    mongo_voters_options = [name: VotersMongoStorage.process_name] ++
                           Application.get_env(:voting, :mongo_voters)
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
      set_child_id(Mongo.child_spec(mongo_girls_options), {Mongo, :girls}),
      set_child_id(Mongo.child_spec(mongo_voters_options), {Mongo, :voters}),
    ]

    Supervisor.init(children, strategy: :one_for_one)

  end
end