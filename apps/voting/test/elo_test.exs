defmodule EloTest do
  @moduledoc false
  use ExUnit.Case
  alias Voting.EloRating

  test "winner has higher rating" do
    assert EloRating.recalculate(1500, 1200) == {1502, 1198}
  end

  test "winner has lower rating" do
    assert EloRating.recalculate(1200, 1500) == {1214, 1486}
  end

  test "winner has zero rating" do
    assert EloRating.recalculate(0, 1500) == {16, 1484}
  end

  test "loser has zero rating" do
    assert EloRating.recalculate(120, 0) == {125, -5}
  end

  test "winner has negative rating" do
    assert EloRating.recalculate(-100, 1500) == {-84, 1484}
  end

  test "loser has negative rating" do
    assert EloRating.recalculate(120, -100) == {124, -104}
  end
end
