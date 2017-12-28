use Mix.Config

config :tg_bot,
       bot_name: "InstaRateLocal",
       bot_username: "InstaRateLocalBot",
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
       messenger: TGBot.Messengers.NadiaLib

config :tg_bot,
       TGBot.Processing.Callbacks,
       messenger: TGBot.Messengers.NadiaLib

config :tg_bot,
       TGBot.Processing.Text,
       messenger: TGBot.Messengers.NadiaLib,
       scheduler: Scheduler.Impls.Mongo,
       admins: [231193206, 309370324]

config :tg_bot,
       TGBot.Processing.Common,
       messenger: TGBot.Messengers.NadiaLib,
       scheduler: Scheduler.Impls.Mongo,
       photos_cache: TGBot.Cache.Impls.Mongo,
       pictures_concatenator: TGBot.Pictures.Concatenators.ImageMagick

config :tg_bot, Scheduler.Reader,
       tasks_storage: Scheduler.Impls.Mongo,
       queue: TGBot.Queue.Impls.Mongo

config :tg_bot, TGBot.Queue.Reader,
       queue: TGBot.Queue.Impls.Mongo,
       fetch_delay: 100

config :tg_bot, TGBot.MatchPhotoCache,
       cache: TGBot.Cache.Impls.Mongo,
       pictures_concatenator: TGBot.Pictures.Concatenators.ImageMagick

config :nadia, token: "480997285:AAEwT3739sBnTz0RSqhEz8TNh4wvJUuqn20"

config :logger, :console,
       metadata: [:request_id, :chat_id]

config :tg_bot, TGBot.Localization,
       disable_translation: false

import_config "#{Mix.env}.exs"
