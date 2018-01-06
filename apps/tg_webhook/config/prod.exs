use Mix.Config

config :nadia, token: "501332340:AAFqDbDgOx6K4GqfuV0dMlOMW5RzoEObtl4"

config :tg_webhook,
       serviced_bots: %{
         "501332340:AAFqDbDgOx6K4GqfuV0dMlOMW5RzoEObtl4" => "insta_queue"
       }

config :utils,
       mongo_queue: [
         database: "local",
         seeds: ["35.189.124.60:27017"],
         collection: "insta_queue",
         max_processing_time: 10000
       ]

