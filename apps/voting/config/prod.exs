use Mix.Config

config :voting,
       mongo_girls: [
         database: "local",
         host: "35.189.124.60",
         port: 27017,
       ],
       mongo_voters: [
         database: "local",
         host: "35.189.124.60",
         port: 27017,
       ],
       mongo_profiles: [
         database: "local",
         host: "35.189.124.60",
         port: 27017,
       ]
