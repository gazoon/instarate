defmodule TGWebhook.Application do
  @moduledoc false

  use Application

  require Logger

  def start(_type, _args) do
    Logger.info "Started application tg_webhook"

    TGWebhook.Supervisor.start_link([])
  end

end
