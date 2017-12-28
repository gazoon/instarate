defmodule TGBot.Supervisor do

  alias TGBot.Chats.Storages.Mongo, as: ChatsMongoStorage
  alias Scheduler.Impls.Mongo, as: MongoScheduler
  alias Utils.Queue.Impls.Mongo, as: MongoQueue
  alias TGBot.Cache.Impls.Mongo, as: MongoCache
  use Supervisor

  def start_link(arg) do
    Supervisor.start_link(__MODULE__, arg)
  end

  def init(_) do
    children = [
                 ChatsMongoStorage.child_spec(),
                 MongoScheduler.child_spec(),
                 MongoQueue.child_spec(),
                 MongoCache.child_spec(),
                 Utils.tasks_supervisor_spec(),
               ]
               |> Kernel.++(Scheduler.Reader.children_spec())
               |> Kernel.++(TGBot.QueueReader.children_spec())

    Supervisor.init(children, strategy: :one_for_one)

  end
end

