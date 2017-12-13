defmodule TGBot.Queue.Reader do

  @config Application.get_env(:tg_bot, __MODULE__)
  @queue @config[:queue]
  use Utils.Reader, otp_app: :tg_bot

  @spec fetch :: {any, String.t}
  def fetch do
    @queue.get_next()
  end

  @spec process({any, String.t}) :: any
  def process({message, processing_id}) do
    Task.Supervisor.start_child(
      :message_workers_supervisor,
      fn ->
        TGBot.on_message(message)
        @queue.finish_processing(processing_id)
      end
    )
  end
end

