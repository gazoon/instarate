defmodule Voting do
  @moduledoc false
  alias Voting.Girls.Girl
  alias Voting.EloRating
  alias Instagram.Media
  require Logger

  @girls_storage Application.get_env(:voting, :girls_storage)
  @voters_storage Application.get_env(:voting, :voters_storage)
  @instagram_client Application.get_env(:voting, :instagram_client)

  @max_random_attempt 10

  @spec add_girl(String.t) :: {:ok, Girl.t} | {:error, String.t}
  def add_girl(photo_uri) do
    photo_code = @instagram_client.parse_media_code(photo_uri)

    with {:ok, media_info = %Media{is_photo: true}} <- @instagram_client.get_media_info(photo_code),
         {:ok, girl} <- Girl.new(media_info.owner, media_info.url)
                        |> @girls_storage.add_girl do
      {:ok, girl}
    else
      {:error, error} -> {:error, error}
      {:ok, %Media{is_photo: false}} -> {:error, "#{photo_code} is not a photo"}
    end
  end

  @spec get_next_pair(String.t) :: {Girl.t, Girl.t}
  def get_next_pair(voters_group_id) do
    attempt = 0
    get_next_pair(voters_group_id, attempt)
  end

  @spec get_girl(String.t) :: {:ok, Girl.t} | {:error, Stringt.t}
  def get_girl(girl_uri) do
    girl_username = @instagram_client.parse_username(girl_uri)
    @girls_storage.get_girl(girl_username)
  end

  @spec get_top(integer) :: [Girl.t]
  def get_top(number) do
    @girls_storage.get_top(number)
  end

  @spec vote(String.t, String.t, String.t, String.t) :: :ok | {:error, String.t}
  def vote(voters_group_id, voter_id, winner_username, loser_username) do
    with :ok <- @voters_storage.try_vote(
      voters_group_id,
      voter_id,
      winner_username,
      loser_username
    ),
         {:ok, winner} <- @girls_storage.get_girl(winner_username),
         {:ok, loser} <- @girls_storage.get_girl(loser_username) do
      process_vote(winner, loser)
      :ok
    else
      error -> error
    end
  end

  @spec process_vote(Girl.t, Girl.t) :: any
  defp process_vote(winner, loser) do
    {new_winner_rating, new_loser_rating} = EloRating.recalculate(winner.rating, loser.rating)
    winner = %{
      winner |
      rating: new_winner_rating,
      matches: winner.matches + 1,
      wins: winner.wins + 1
    }
    loser = %{
      loser |
      rating: new_loser_rating,
      matches: loser.matches + 1,
      loses: loser.loses + 1
    }
    @girls_storage.update_girl(winner)
    @girls_storage.update_girl(loser)
  end

  defp get_next_pair(_voters_group_id, _attempt = @max_random_attempt) do
    raise "Can't get girl to vote to, it seems you've voted for a lot of girls"
  end
  @spec get_next_pair(String.t, integer) :: {Girl.t, Girl.t}
  defp get_next_pair(voters_group_id, attempt) do
    {girl_one, girl_two} = @girls_storage.get_random_pair()
    if @voters_storage.new_pair?(voters_group_id, girl_one.username, girl_two.username) do
      {girl_one, girl_two}
    else
      get_next_pair(voters_group_id, attempt + 1)
    end
  end
end
