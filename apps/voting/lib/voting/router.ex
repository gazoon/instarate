defmodule Voting.Router do
  @moduledoc false

  use Plug.Router
  alias Voting.Voters.Storages.Mongo, as: VotersStorage
  alias Voting.Girls.Storages.Mongo, as: GirlsStorage
  alias Voting.Girls.Girl


  plug :match
  plug :dispatch

  get "/" do
    x = VotersStorage.can_vote?("1", "1", "2")
    IO.inspect(x)
#    x = VotersStorage.try_vote("1", "3", "2")
#    IO.inspect(x)
    x = GirlsStorage.get_girl("1")
    x = %{x | rating: x.rating + 10}
    GirlsStorage.update_girl(x)
    g = %Girl{username: "32", rating: 222, photo: "llll", added_at: 2222}
    x = GirlsStorage.add_girl(g)
    IO.inspect(x)
    conn
    |> send_resp(200, "Plug!")
  end

end
