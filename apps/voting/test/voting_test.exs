defmodule VotingTest do
  use ExUnit.Case

  import Voting
  alias Voting.InstagramProfiles.Model, as: Profile
  alias Voting.Girl

  test "simple add" do
    assert {:ok, %Profile{username: "media_code_owner"}} = add_girl("media_code")
  end


end
