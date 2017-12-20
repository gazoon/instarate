use Mix.Config

config :tg_bot,
       bot_name: "InstaRate",
       bot_username: "InstaRateBot",
       mongo_chats: [
         database: "local",
         host: "localhost",
         port: 27017,
       ],
       mongo_scheduler: [
         database: "local",
         host: "localhost",
         port: 27017,
       ],
       mongo_queue: [
         database: "local",
         host: "localhost",
         port: 27017,
         collection: "insta_queue",
         max_processing_time: 10000
       ],
       mongo_cache: [
         database: "local",
         host: "localhost",
         port: 27017,
       ]

config :tg_bot, TGBot,
       chats_storage: TGBot.Chats.Storages.Mongo,
       messenger: TGBot.Messengers.NadiaLib,
       scheduler: Scheduler.Impls.Mongo,
       photos_cache: TGBot.Cache.Impls.Mongo,
       admins: [231193206, 309370324]

config :tg_bot, Scheduler.Reader,
       tasks_storage: Scheduler.Impls.Mongo,
       queue: TGBot.Queue.Impls.Mongo

config :tg_bot, TGBot.Queue.Reader,
       queue: TGBot.Queue.Impls.Mongo,
       fetch_delay: 100

config :tg_bot, TGBot.MatchPhotoCache,
       cache: TGBot.Cache.Impls.Mongo

config :nadia, token: "501332340:AAGMi61i2NEYAJR6-GnqwHAE5MYpBKwOjo0"

config :logger, :console,
       metadata: [:request_id, :chat_id]

import_config "#{Mix.env}.exs"
