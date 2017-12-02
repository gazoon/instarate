defmodule Voting.Router do
  @moduledoc false

  use Plug.Router
  alias Voting.Voters.Storages.Mongo, as: VotersStorage
  alias Voting.Girls.Storages.Mongo, as: GirlsStorage
  alias Voting.Girls.Girl
  alias Instagram.Clients.Http, as: InstagramClient
  alias Voting.EloRating


  plug :match
  plug :dispatch

  get "/" do
    #    IO.inspect(Voting.get_next_pair("228"))

    #    IO.inspect(Voting.vote("228", "svetabily", "32"))
    #    IO.inspect(Voting.get_top(3))
    IO.inspect Voting.add_girl("Bap6TqcjowK")
    IO.puts("fff")
    #    IO.inspect(x)
    conn
    |> send_resp(200, "Plug!")
  end

end
