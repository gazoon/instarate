use Mix.Config

config :voting,
       mongo_girls: [
         database: "local",
         seeds: ["35.189.124.60:27017"],
       ],
       mongo_voters: [
         database: "local",
         seeds: ["35.189.124.60:27017"],
       ],
       mongo_profiles: [
         database: "local",
         seeds: ["35.189.124.60:27017"],
       ]
