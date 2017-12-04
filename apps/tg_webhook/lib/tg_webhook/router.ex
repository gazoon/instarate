defmodule TGWebhook.Router do
  use Plug.Router

  plug :match
  plug :dispatch

  get "/" do
    IO.puts("fff")
    conn
    |> send_resp(200, "Plug!")
  end

end
