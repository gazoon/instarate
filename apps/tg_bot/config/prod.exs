use Mix.Config

config :tg_bot,
       bot_name: "InstaRate",
       bot_username: "InstaRateBot",
       mongo_chats: [
         database: "local",
         seeds: ["35.189.124.60:27017"],
       ],
       mongo_scheduler: [
         database: "local",
         seeds: ["35.189.124.60:27017"],
       ],
       mongo_cache: [
         database: "local",
         seeds: ["35.189.124.60:27017"],
       ]

config :utils,
       mongo_queue: [
         database: "local",
         seeds: ["35.189.124.60:27017"],
         collection: "insta_queue",
         max_processing_time: 10000
       ]


config :nadia, token: "501332340:AAFqDbDgOx6K4GqfuV0dMlOMW5RzoEObtl4"
