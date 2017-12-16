defmodule Mix.Tasks.MigrateGirls do
  use Mix.Task
  alias Voting.Competitors.Storages.Mongo, as: CompetitorsStorage
  alias Voting.InstagramProfiles.Storages.Mongo, as: ProfilesStorage
  alias Voting.Competitors.Model, as: Competitor
  alias Voting.InstagramProfiles.Model, as: Profile
  require Logger
  def run(_) do
    Application.ensure_all_started(:voting)
    rows = Mongo.find(
      CompetitorsStorage.process_name(),
      CompetitorsStorage.collection(),
      %{}
    )
    Enum.each(
      rows,
      fn (row) ->
        competitor = %Competitor{
          username: row["username"],
          rating: row["rating"],
          matches: row["matches"],
          wins: row["wins"],
          loses: row["loses"],
        }
        competitor_all = %{competitor | competition: Voting.global_competition()}
        competitor_celebrities = %{competitor | competition: Voting.celebrities_competition()}
        profile = %Profile{
          username: row["username"],
          photo: row["photo"],
          added_at: row["added_at"],
          followers: 1_000_000
        }
        case ProfilesStorage.add(profile) do
          {:error, error} -> Logger.warn(error)
          _ ->
            CompetitorsStorage.add_girl(competitor_all)
            CompetitorsStorage.add_girl(competitor_celebrities)
            Logger.info("Girl #{profile.username} migrated")
        end
      end
    )
  end

end
