defmodule Voting.Girls.Girl do
  @moduledoc false
  alias Voting.Girls.Girl
  @initial_rating 1500

  @type t :: %Girl{
               username: String.t,
               photo: String.t,
               added_at: integer,
               rating: integer,
               matches: integer,
               wins: integer,
               loses: integer
             }

  defstruct username: nil,
            photo: nil,
            added_at: nil,
            rating: @initial_rating,
            matches: 0,
            wins: 0,
            loses: 0

  def new(username, photo) do
    current_time = DateTime.utc_now()
                   |> DateTime.to_unix()
    %Girl{username: username, photo: photo, added_at: current_time}
  end


end
