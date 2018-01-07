defmodule Mix.Tasks.DeleteGirls do
  use Mix.Task
  require Logger
  def run(args) do
    usernames = args
    Application.ensure_all_started(:voting)
    Voting.delete_girls(usernames)
    Logger.info("Deleted girls: #{inspect usernames}")
  end
end

