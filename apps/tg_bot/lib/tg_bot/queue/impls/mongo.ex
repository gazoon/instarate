defmodule TGBot.Queue.Impls.Mongo do
  alias TGBot.Queue.{Consumer, Producer}
  @behaviour Consumer
  @behaviour Producer
  @process_name :mongo_queue
  @config Application.get_env(:tg_bot, :mongo_queue)
  @collection @config[:collection]

  @spec child_spec :: tuple
  def child_spec do
    options = [name: @process_name] ++ Keyword.pop(@config, :collection)
    Utils.set_child_id(Mongo.child_spec(options), {Mongo, :queue})
  end

  @spec put(integer, any, Keyword.t) :: any
  def put(chat_id, message, opts \\ []) do
    message_id = UUID.uuid4()
    Mongo.update_one(%{"chat_id" => chat_id})
  end

  @spec get_next :: {any, String.t}
  def get_next do

  end

  @spec finish_processing(String.t) :: any
  def finish_processing(processing_id)do

  end
end
