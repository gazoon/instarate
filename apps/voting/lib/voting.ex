defmodule Voting do
  @moduledoc false
  alias Voting.Girls.Girl
  alias Voting.EloRating
  require Logger

  @girls_storage Application.get_env(:voting, :girls_storage)
  @voters_storage Application.get_env(:voting, :voters_storage)
  @instagram_client Application.get_env(:voting, :instagram_client)

  @max_random_attempt 10

  @spec add_girl(String.t) :: {:ok, Girl.t} | {:error, String.t}
  def add_girl(photo_id) do
    photo_code = retrieve_photo_code(photo_id)

    with {:ok, username} <- @instagram_client.get_media_owner(photo_code),
         true <- @instagram_client.is_photo?(photo_code),
         {:ok, girl} <- photo_code
                        |> @instagram_client.build_media_url()
                        |> (&Girl.new(username, &1)).()
                        |> @girls_storage.add_girl do
      {:ok, girl}
    else
      {:error, error} -> {:error, error}
      false -> {:error, "#{photo_code} is not a photo"}
    end
  end

  @spec get_next_pair(String.t) :: {Girl.t, Girl.t}
  def get_next_pair(voter_id) do
    attempt = 0
    get_next_pair(voter_id, attempt)
  end


  @spec get_top(integer) :: [Girl.t]
  def get_top(number) do
    @girls_storage.get_top(number)
  end

  @spec vote(String.t, String.t, String.t) :: :ok
  def vote(voter_id, winner_username, loser_username) do
    with :ok <- @voters_storage.try_vote(voter_id, winner_username, loser_username),
         {:ok, winner} <- @girls_storage.get_girl(winner_username),
         {:ok, loser} <- @girls_storage.get_girl(loser_username) do
      process_vote(winner, loser)
    else
      {:error, error} -> Logger.warn "Can't vote: #{error}"
    end
  end

  @spec process_vote(Girl.t, Girl.t) :: :ok
  defp process_vote(winner, loser) do
    {new_winner_rating, new_loser_rating} = EloRating.recalculate(winner.rating, loser.rating)
    winner = %{
      winner |
      rating: new_winner_rating,
      matches: winner.matches + 1,
      wins: winner.wins + 1
    }
    loser = %{loser | rating: new_loser_rating, matches: loser.matches + 1, loses: loser.loses + 1}
    @girls_storage.update_girl(winner)
    @girls_storage.update_girl(loser)
    :ok
  end

  @spec retrieve_photo_code(String.t) :: String.t
  defp retrieve_photo_code(photo_link) do
    URI.parse(photo_link).path
    |> Path.split()
    |> List.last()
  end

  defp get_next_pair(_voting, _attempt = @max_random_attempt) do
    raise "Can't get girl to vote to, it seems you've voted for a lot of girls"
  end
  @spec get_next_pair(String.t, integer) :: {Girl.t, Girl.t}
  defp get_next_pair(voter_id, attempt) do
    {girl_one, girl_two} = @girls_storage.get_random_pair()
    if @voters_storage.can_vote?(voter_id, girl_one.username, girl_two.username) do
      {girl_one, girl_two}
    else
      get_next_pair(voter_id, attempt + 1)
    end
  end
end
