defmodule Voting.Fakes.InstagramProfiles do
  alias Voting.InstagramProfiles.Storage
  alias Voting.InstagramProfiles.Model, as: Profile
  @behaviour Storage

  @spec add(Profile.t) :: {:ok, Profile.t} | {:error, String.t}
  def add(girl) do
    if String.starts_with?(girl.username, "existent_girl") do
      {:error, "already added"}
    else
      {:ok, girl}
    end
  end

  @spec get_multiple([String.t]) :: [Profile.t]
  def get_multiple(usernames) do
    for username <- usernames, do: %Profile{username: username}
  end

  # not used
  @spec get(String.t) :: Profile.t
  def get(_username), do: %Profile{}

  # not used
  @spec delete([String.t]) :: :ok
  def delete(_usernames), do: :ok
end
