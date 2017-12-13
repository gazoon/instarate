defmodule TGBot.Supervisor do

  alias TGBot.Chats.Storages.Mongo, as: ChatsMongoStorage
  alias Scheduler.Impls.Mongo, as: MongoScheduler
  alias TGBot.Queue.Impls.Mongo, as: MongoQueue
  use Supervisor

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end

  def init(_) do
    children = [
                 ChatsMongoStorage.child_spec(),
                 MongoScheduler.child_spec(),
                 MongoQueue.child_spec(),
                 {Task.Supervisor, name: :message_workers_supervisor}
               ]
               |> Kernel.++(Scheduler.Reader.children_spec())
               |> Kernel.++(TGBot.Queue.Reader.children_spec())

    Supervisor.init(children, strategy: :one_for_one)

  end
end

