defmodule VotingTest do
  @moduledoc false
  use ExUnit.Case
  alias Voting.Girls.Girl

  test "winner has higher rating" do
    assert {
             :ok,
             %Girl{username: "svetabily", photo: "https://www.instagram.com/p/BbRv15lltW7/"}
           } = Voting.add_girl("BbRv15lltW7")

  end
end
