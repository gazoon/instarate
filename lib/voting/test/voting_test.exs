defmodule VotingTest do
  use ExUnit.Case

  import Voting
  alias Voting.InstagramProfiles.Model, as: Profile
  alias Voting.Competitors.Model, as: Competitor

  @global_competiton global_competition()
  @normal_competiton normal_competition()
  @models_competiton models_competition()
  @celebrities_competition celebrities_competition()

  test "add normal girl" do
    username = "media_code_owner"
    photo_code = "media_code"
    expected = %Profile{username: username, photo_code: photo_code}
    assert {:ok, got} = add_girl(photo_code)
    assert got.username == expected.username
    assert got.photo_code == expected.photo_code
    assert_received {
      :competitor_added,
      %Competitor{username: ^username, competition: @global_competiton}
    }
    assert_received {
      :competitor_added,
      %Competitor{username: ^username, competition: @normal_competiton}
    }
    refute_received {:competitor_added, _}
  end

  test "add model girl" do
    username = "model_media_owner"
    photo_code = "model_media"
    expected = %Profile{username: username, photo_code: photo_code}
    assert {:ok, got} = add_girl(photo_code)
    assert got.username == expected.username
    assert got.photo_code == expected.photo_code
    assert_received {
      :competitor_added,
      %Competitor{username: ^username, competition: @global_competiton}
    }
    assert_received {
      :competitor_added,
      %Competitor{username: ^username, competition: @models_competiton}
    }
    refute_received {:competitor_added, _}
  end

  test "add celebrity girl" do
    username = "celebrity_media_owner"
    photo_code = "celebrity_media"
    expected = %Profile{username: username, photo_code: photo_code}
    assert {:ok, got} = add_girl(photo_code)
    assert got.username == expected.username
    assert got.photo_code == expected.photo_code
    assert_received {
      :competitor_added,
      %Competitor{username: ^username, competition: @global_competiton}
    }
    assert_received {
      :competitor_added,
      %Competitor{username: ^username, competition: @celebrities_competition}
    }
    refute_received {:competitor_added, _}
  end

  test "add not photo" do
    assert {:error, _} = add_girl("not_photo_media")
  end

  test "add existent girl" do
    assert {:error, _} = add_girl("existent_girl_media")
  end

  test "get next pair, can get" do
    competition = "some competition"
    assert {girl_one, girl_two} = get_next_pair(competition, "some voters")
    assert girl_one.username == "first_girl"
    assert girl_one.competition == competition
    assert girl_two.username == "second_girl"
    assert girl_two.competition == competition
  end

  test "get next pair, no more girls" do
    competition = "some competition"
    assert get_next_pair(competition, "already_saw_pair") == :error
  end

  test "vote first time" do
    competition = "some competition"
    assert {:ok, {winner, loser}} = Voting.vote(
             competition,
             "some voters",
             "voter",
             "first_girl",
             "second_girl"
           )
    assert winner.wins == 1
    assert winner.loses == 0
    assert winner.matches == 1
    assert loser.wins == 0
    assert loser.loses == 1
    assert loser.matches == 1
    assert winner.rating == 1508
    assert loser.rating == 1492
  end

end
