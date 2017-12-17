defmodule Voting.Voters.Storages.Mongo do

  @behaviour Voting.Voters.Storage

  @collection "voters"
  @duplication_code 11_000

  @process_name :mongo_voters

  @spec child_spec :: tuple
  def child_spec do
    options = [name: @process_name, name: DBConnection.Poolboy] ++
              Application.get_env(:voting, :mongo_voters)
    Utils.set_child_id(Mongo.child_spec(options), {Mongo, :voters})
  end

  @spec try_vote(String.t, String.t, String.t, String.t, String.t) :: :ok | {:error, String.t}
  def try_vote(competition, voters_group_id, voter_id, girl_one_id, girl_two_id)  do
    girls_id = to_girls_id(girl_one_id, girl_two_id)
    insert_result = Mongo.insert_one(
      @process_name,
      @collection,
      %{
        competition: competition,
        voters_group: voters_group_id,
        voter: voter_id,
        girls_id: girls_id
      },
      name: DBConnection.Poolboy
    )
    case insert_result do
      {:ok, _} -> :ok
      {:error, %Mongo.Error{code: @duplication_code}} ->
        {
          :error,
          "#{voters_group_id} #{voter_id} already voted for #{girl_one_id} and #{girl_two_id}"
        }
      {:error, error} -> raise error
    end

  end

  @spec new_pair?(String.t, String.t, String.t, String.t) :: boolean
  def new_pair?(competition, voters_group_id, girl_one_id, girl_two_id) do
    girls_id = to_girls_id(girl_one_id, girl_two_id)
    row = Mongo.find_one(
      @process_name,
      @collection,
      %{competition: competition, voters_group: voters_group_id, girls_id: girls_id},
      name: DBConnection.Poolboy
    )
    !row
  end

  @spec to_girls_id(String.t, String.t) :: String.t
  defp to_girls_id(girl_one_id, girl_two_id) do
    [girl_one_id, girl_two_id]
    |> Enum.sort
    |> Enum.join(" | ")
  end

end
