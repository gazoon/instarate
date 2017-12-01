defmodule TGWebhook.Router do
  @moduledoc false

  use Plug.Router


  plug :match
  plug :dispatch

  get "/" do
    IO.puts("fff")
    conn
    |> send_resp(200, "Plug!")
  end

end
