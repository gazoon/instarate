defmodule VotingTest do
  use ExUnit.Case
  doctest Voting

  test "greets the world" do
    IO.inspect(Application.get_all_env(:voting))
    IO.inspect Application.get_env(:voting, :foo)
    assert :world == :world
  end
end
