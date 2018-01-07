defmodule Voting.EloRating do

  @significance_coeff 16

  @spec recalculate(integer, integer) :: {integer, integer}
  def recalculate(winner_rating, loser_rating) do
    new_winner_rating = calculate_new_rating(winner_rating, loser_rating, is_winner: true)
    new_loser_rating = calculate_new_rating(loser_rating, winner_rating, is_winner: false)
    {new_winner_rating, new_loser_rating}
  end

  @spec expected_probabilty(integer, integer) :: float
  defp expected_probabilty(rating_a, rating_b) do
    1 / (1 + :math.pow(10, (rating_b - rating_a) / 400))
  end

  #  @spec calculate_new_rating(integer,integer,is_winner)
  defp calculate_new_rating(rating_a, rating_b, is_winner: is_winner) do
    ea = expected_probabilty(rating_a, rating_b)
    actual_result = if is_winner, do: 1, else: 0
    new_rating = round(rating_a + @significance_coeff * (actual_result - ea))
    new_rating
  end
end
