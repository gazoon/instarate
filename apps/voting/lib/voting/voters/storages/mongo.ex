defmodule Voting.Voters.Storages.Mongo do
  @moduledoc false
  alias Voting.Utils
  @behaviour Voting.Voters.Storage

  @collection "voters"
  @duplication_code 11000

  @process_name :mongo_voters
  def process_name, do: @process_name

  def init() do
    mongo_options = [name: @process_name] ++ Application.get_env(:voting, :voters_girls)
    Utils.set_child_id(Mongo.child_spec(mongo_options), {Mongo, :voters})
  end

  @spec try_vote(String.t, String.t, String.t) :: :ok | {:error, String.t}
  def try_vote(voter_id, girl_one_id, girl_two_id)  do
    girls_id = to_girls_id(girl_one_id, girl_two_id)
    insert_result = Mongo.insert_one(
      @process_name,
      @collection,
      %{voter: voter_id, girls_id: girls_id}
    )
    case insert_result do
      {:ok, _} -> :ok
      {:error, %Mongo.Error{code: @duplication_code}} ->
        {:error, "#{voter_id} already voted for #{girl_one_id} and #{girl_two_id}"}
    end

  end

  # TODO: add specs everywhere
  @spec can_vote?(String.t, String.t, String.t) :: boolean
  def can_vote?(voter_id, girl_one_id, girl_two_id) do
    girls_id = to_girls_id(girl_one_id, girl_two_id)
    row = Mongo.find_one(
      @process_name,
      @collection,
      %{voter: voter_id, girls_id: girls_id}
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
