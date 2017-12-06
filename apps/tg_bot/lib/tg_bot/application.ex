defmodule TGBot.Application do

  use Application

  require Logger

  def start(_type, _args) do
    Logger.info "Started application tg_bot"

    TGBot.Supervisor.start_link([])
  end

end
