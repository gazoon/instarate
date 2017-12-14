defmodule Voting.InstagramProfiles.Storages.Mongo do
  alias Voting.InstagramProfiles.Storage
  alias Voting.InstagramProfiles.Model, as: Profile
  @behaviour Storage

  @collection "insta_profiles"
  @duplication_code 11_000

  @process_name :mongo_profiles

  @spec child_spec :: tuple
  def child_spec do
    options = [name: @process_name] ++ Application.get_env(:voting, :mongo_profiles)
    Utils.set_child_id(Mongo.child_spec(options), {Mongo, :profiles})
  end

  @spec add(Profile.t) :: {:ok, Profile.t} | {:error, String.t}
  def add(girl) do
    insert_result = Mongo.insert_one(@process_name, @collection, girl)
    case insert_result do
      {:ok, _} ->
        {:ok, girl}
      {:error, %Mongo.Error{code: @duplication_code}} ->
        {:error, "Girl #{girl.username} already added"}
      {:error, error} -> raise error
    end
  end

  @spec get(String.t) :: Profile.t
  def get(username) do
    row = Mongo.find_one(@process_name, @collection, %{username: username})
    if row, do: transform_girl(row), else: raise "dd"
  end

  @spec get_multiple([String.t]) :: [Profile.t]
  def get_multiple(usernames) do
    rows = Mongo.find(
      @process_name,
      @collection,
      %{
        username: %{
          "$in" => usernames
        }
      }
    )
    Enum.map(rows, &transform_profile/1)
  end

  @spec transform_profile(map()) :: Profile.t
  defp transform_profile(row) do
    %Profile{
      username: row["username"],
      added_at: row["added_at"],
      photo: row["photo"],
      followers: row["followers"],
    }
  end
end
