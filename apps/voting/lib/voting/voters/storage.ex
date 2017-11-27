defmodule Voting.Voters.Storage do
  @moduledoc false

  @callback try_vote(
              voter_id :: String.t,
              girl_one_id :: String.t,
              girl_two_id :: String.t
            ) :: :ok | {:error, String.t}

  @callback can_vote?(
              voter_id :: String.t,
              girl_one_id :: String.t,
              girl_two_id :: String.t
            ) :: boolean
end
