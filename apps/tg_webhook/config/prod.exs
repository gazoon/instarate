use Mix.Config

config :utils,
       mongo_queue: [
         database: "local",
         host: "35.189.124.60",
         port: 27017,
         collection: "insta_queue",
         max_processing_time: 10000
       ]

config :nadia, token: "501332340:AAFqDbDgOx6K4GqfuV0dMlOMW5RzoEObtl4"
