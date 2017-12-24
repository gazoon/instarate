defmodule Mix.Tasks.Migrations.UploadPhotos do
  use Mix.Task
  require Logger
  alias Voting.InstagramProfiles.Storages.Mongo, as: Profiles
  alias Voting.Files.Storages.Google, as: FilesStorage
  def run(_args) do
    Application.ensure_all_started(:voting)
    rows = Mongo.find(
      Profiles.process_name(),
      Profiles.collection(),
      %{},
      pool: DBConnection.Poolboy
    )

    Enum.each(
      rows,
      fn (row) ->
        old_photo_url = row["photo"]
        if String.starts_with?(old_photo_url, "http") do
          username = row["username"]
          path = username <> "-" <> UUID.uuid4()
          Logger.info("Process #{username}")
          FilesStorage.upload(path, old_photo_url)
          Mongo.update_one!(
            Profiles.process_name(),
            Profiles.collection(),
            %{username: row["username"]},
            %{
              "$set" => %{
                photo: path
              }
            },
            pool: DBConnection.Poolboy
          )
        end
      end
    )
  end
end
