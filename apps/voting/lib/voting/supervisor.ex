defmodule Voting.Supervisor do
  @moduledoc false



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
    # TODO: try to use :name_girls for config field too,
    mongo_girls_options = [name: :mongo_girls] ++ Application.get_env(:voting, :mongo_girls)
    mongo_voters_options = [name: :mongo_voters] ++ Application.get_env(:voting, :mongo_voters)
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