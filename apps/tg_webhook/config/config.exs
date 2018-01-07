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
#     config :tg_webhook, key: :value
#
# and access this configuration in your application as:
#
#     Application.get_env(:tg_webhook, :key)
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
#     import_config "#{Mix.env}.exs"

config :tg_webhook, TGWebhook.Update,
       queue: Utils.Queue.Impls.Mongo

config :tg_webhook,
       port: 8088

config :utils,
       mongo_queue: [
         database: "local",
         seeds: ["35.198.77.219:27017"],
         collection: "insta_queue",
         max_processing_time: 10000
       ]

import_config "#{Mix.env}.exs"
