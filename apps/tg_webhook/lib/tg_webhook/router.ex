defmodule TGWebhook.Router do

  use Plug.Router

  @bot_to_queue Application.get_env(:tg_webhook, :known_bots)


  plug :match
  plug Plug.Parsers,
       parsers: [:json],
       pass: ["application/json"],
       json_decoder: Poison
  plug :dispatch
  use Plug.ErrorHandler
  alias TGWebhook.Update

  post "/update/:bot_token" do
    queue_name = Map.get(@bot_to_queue, bot_token)
    request_data = conn.body_params
    if queue_name && request_data != %{} do
      Update.process(queue_name, request_data)
      send_resp(conn, 200, "OK")
    else
      send_resp(conn, 404, "Not found")
    end
  end

  match _ do
    send_resp(conn, 404, "Not found")
  end


  defp handle_errors(conn, %{kind: _kind, reason: _reason, stack: _stack}) do
    send_resp(conn, conn.status, "Something went wrong")
  end
end

