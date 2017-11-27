defmodule Voting.Application do
  @moduledoc false

  use Application

  require Logger

  def start(_type, _args) do
    Logger.info "Started application"

    Voting.Supervisor.start_link([])
  end

end
