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
       pictures_concatenator: TGBot.Pictures.Concatenators.ImageMagick,
       admins: [231193206, 309370324]

config :tg_bot, Scheduler.Reader,
       tasks_storage: Scheduler.Impls.Mongo,
       queue: TGBot.Queue.Impls.Mongo

config :tg_bot, TGBot.Queue.Reader,
       queue: TGBot.Queue.Impls.Mongo,
       fetch_delay: 100

config :tg_bot, TGBot.MatchPhotoCache,
       cache: TGBot.Cache.Impls.Mongo

config :nadia, token: "501332340:AAFqDbDgOx6K4GqfuV0dMlOMW5RzoEObtl4"

config :logger, :console,
       metadata: [:request_id, :chat_id]

config :tg_bot, TGBot.Localization,
       disable_translation: false

import_config "#{Mix.env}.exs"
