use Mix.Config

config :voting,
       Voting,
       girls_storage: Voting.Fakes.Competitors,
       voters_storage: Voting.Fakes.Voters,
       profiles_storage: Voting.Fakes.InstagramProfiles

config :voting,
       Voting.InstagramProfiles.Model,
       photos_storage: Voting.Fakes.Files


config :utils, Instagram.Client, Voting.Fakes.InstagramClient

