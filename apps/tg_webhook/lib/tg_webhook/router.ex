defmodule TGWebhook.Router do

  use Plug.Router

  plug :match
  plug :dispatch

  get "/" do
    #    IO.inspect(Voting.get_next_pair("228"))

    #    IO.inspect(Voting.vote("228", "svetabily", "32"))
    #    IO.inspect(Voting.get_top(3))
    #    IO.inspect(x)
    conn
    |> send_resp(200, "Plug!")
  end

end

