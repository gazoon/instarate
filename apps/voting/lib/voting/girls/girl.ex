defmodule Voting.Girls.Girl do
  @moduledoc false
  alias Voting.Girls.Girl
  @initial_rating 1500

  @type t :: %Girl{username: String.t, photo: String.t, added_at: integer, rating: integer}

  defstruct username: nil, photo: nil, added_at: nil, rating: @initial_rating


end
