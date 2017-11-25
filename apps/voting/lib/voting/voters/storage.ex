defmodule Voting.Voters.Storage do
  @moduledoc false

  @callback try_vote(
              voter_id :: String.t,
              first_girl_id :: String.t,
              second_girl_id :: String.t
            ) :: :ok | :error

  @callback can_vote?(
              voter_id :: String.t,
              first_girl_id :: String.t,
              second_girl_id :: String.t
            ) :: boolean
end
