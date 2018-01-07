defmodule Voting.Fakes.Voters do
  @behaviour Voting.Voters.Storage

  @spec try_vote(String.t, String.t, String.t, String.t, String.t) :: :ok | {:error, String.t}
  def try_vote(_competition, _voters_group_id, _voter_id, _girl_one_id, _girl_two_id), do: :ok

  @spec new_pair?(String.t, String.t, String.t, String.t) :: boolean
  def new_pair?(_competition, voters_group_id, _girl_one_id, _girl_two_id) do
    if voters_group_id == "already_saw_pair", do: false, else: true
  end
end
