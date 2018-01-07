defmodule Voting.Application do

  use Application

  require Logger

  def start(_type, _args) do
    Logger.info "Started application voting"

    Voting.Supervisor.start_link([])
  end

end
