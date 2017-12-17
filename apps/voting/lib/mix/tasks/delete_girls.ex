defmodule Mix.Tasks.DeleteGirls do
  use Mix.Task
  alias Voting.Competitors.Storages.Mongo, as: CompetitorsStorage
  alias Voting.InstagramProfiles.Storages.Mongo, as: ProfilesStorage
  require Logger
  def run(args) do
    usernames = args
    Application.ensure_all_started(:voting)
    Mongo.delete_many!(
      CompetitorsStorage.process_name(),
      CompetitorsStorage.collection(),
      %{
        username: %{
          "$in" => usernames
        }
      }
    )
    Mongo.delete_many!(
      ProfilesStorage.process_name(),
      ProfilesStorage.collection(),
      %{
        username: %{
          "$in" => usernames
        }
      }
    )
    Logger.error("Deleted girls: #{inspect usernames}")
  end
end

