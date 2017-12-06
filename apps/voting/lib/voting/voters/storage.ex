defmodule Voting.Voters.Storage do
  @type t :: module
  @callback try_vote(
              voters_group_id :: String.t,
              voter_id :: String.t,
              girl_one_id :: String.t,
              girl_two_id :: String.t
            ) :: :ok | {:error, String.t}

  @callback new_pair?(
              voters_group_id :: String.t,
              girl_one_id :: String.t,
              girl_two_id :: String.t
            ) :: boolean
end
