defmodule VotingTest do
  use ExUnit.Case
  doctest Voting

  test "greets the world" do
    assert Voting.hello() == :world
  end
end
