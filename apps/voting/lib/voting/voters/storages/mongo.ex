defmodule Voting.Voters.Storages.Mongo do
  @moduledoc false
  #  alias Voting.Voters.Storages.Mongo, as: MongoStorage
  @behaviour Voting.Voters.Storage

  @collection "voters"
  @duplication_code 11000

  def try_vote(voter_id, girl_one_id, girl_two_id)  do
    girls_id = to_girls_id(girl_one_id, girl_two_id)
    insert_result = Mongo.insert_one(
      :mongo_voters,
      @collection,
      %{voter: voter_id, girls_id: girls_id}
    )
    case insert_result do
      {:ok, _} -> :ok
      {:error, %Mongo.Error{code: @duplication_code}} -> :error
    end

  end

  def can_vote?(voter_id, girl_one_id, girl_two_id) do
    girls_id = to_girls_id(girl_one_id, girl_two_id)
    row = Mongo.find_one(
      :mongo_voters,
      @collection,
      %{voter: voter_id, girls_id: girls_id}
    )
    !!row

  end

  defp to_girls_id(girl_one_id, girl_two_id) do
    [girl_one_id, girl_two_id]
    |> Enum.sort
    |> Enum.join(" | ")
  end

end
