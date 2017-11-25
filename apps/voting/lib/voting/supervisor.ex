defmodule Voting.Supervisor do
  @moduledoc false



  use Supervisor

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end

  def init(_) do
    mongo_options = [name: :mongo] ++ Application.get_env(:voting, :mongodb)
    children = [
      # add pool
      {Mongo, mongo_options},
    ]

    Supervisor.init(children, strategy: :one_for_one)
  end
end