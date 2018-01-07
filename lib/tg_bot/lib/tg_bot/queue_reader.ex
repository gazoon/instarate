defmodule TGBot.QueueReader do

  @config Application.get_env(:tg_bot, __MODULE__)
  @queue @config[:queue]
  use Utils.Reader, otp_app: :tg_bot

  @spec fetch :: {any, String.t}
  def fetch do
    @queue.get_next()
  end

  @spec process({any, String.t}) :: any
  def process({message, processing_id}) do
    Task.start(
      fn ->
        try do
          TGBot.on_message(message)
        after
          @queue.finish_processing(processing_id)
        end
      end
    )
  end
end

