defmodule TGBot.Supervisor do

  alias TGBot.Chats.Storages.Mongo, as: ChatsMongoStorage
  alias Scheduler.Impls.Mongo, as: MongoScheduler
  use Supervisor

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end

  def init(_) do
    children = [
      ChatsMongoStorage.child_spec(),
      MongoScheduler.child_spec(),
    ]
    children = children ++ Scheduler.Reader.children_spec()

    Supervisor.init(children, strategy: :one_for_one)

  end
end

