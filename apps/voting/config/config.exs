# This file is responsible for configuring your application
# and its dependencies with the aid of the Mix.Config module.
use Mix.Config

# This configuration is loaded before any dependency and is restricted
# to this project. If another project depends on this project, this
# file won't be loaded nor affect the parent project. For this reason,
# if you want to provide default values for your application for
# 3rd-party users, it should be done in your "mix.exs" file.

# You can configure your application as:
#
config :voting,
       Voting,
       girls_storage: Voting.Competitor.Storages.Mongo,
       voters_storage: Voting.Voters.Storages.Mongo,
       profiles_storage: Voting.InstagramProfiles.Storages.Mongo

config :voting,
       Voting.Competitors.Model,
       storage: Voting.Competitors.Storages.Mongo

config :voting,
       Voting.Girl,
       storage: Voting.Competitors.Storages.Mongo

config :utils, Instagram.Client, Instagram.Clients.Http

config :voting,
       mongo_girls: [
         database: "local",
         host: "localhost",
         port: 27017,
       ],
       mongo_voters: [
         database: "local",
         host: "localhost",
         port: 27017,
       ],
       mongo_profiles: [
         database: "local",
         host: "localhost",
         port: 27017,
       ]
#
# and access this configuration in your application as:
#
#     Application.get_env(:voting, :key)
#
# You can also configure a 3rd-party app:
#
#     config :logger, level: :info
#

# It is also possible to import configuration files, relative to this
# directory. For example, you can emulate configuration per environment
# by uncommenting the line below and defining dev.exs, test.exs and such.
# Configuration from the imported file will override the ones defined
# here (which is why it is important to import them last).
#
import_config "#{Mix.env}.exs"
