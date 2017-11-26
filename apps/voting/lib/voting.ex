defmodule Voting do
  @moduledoc false
  alias Voting.Girls.Girl

  @girls_storage Application.get_env(:voting, :girls_storage)
  @voters_storage Application.get_env(:voting, :voters_storage)
  @instagram_client Application.get_env(:voting, :instagram_client)

  @spec add_girl(String.t) :: :ok | {:error, String.t}
  def add_girl(photo_id) do
    photo_code = retrieve_photo_code(photo_id)
  end

  @spec get_next_pair(String.t) :: {:ok, {Girl.t, Girl.t}} | {:error, String.t}
  def get_next_pair(voter_id) do

  end

  @spec get_top(integer) :: [Girl.t]
  def get_top(number) do

  end

  @spec vote(String.t, String.t, String.t) :: :ok
  def vote(voter_id, winner_id, loser_id) do

  end

  @spec retrieve_photo_code(String.t) :: String.t
  defp retrieve_photo_code(photo_link) do
    URI.parse(photo_link).path
    |> Path.split()
    |> List.last()
  end
end
