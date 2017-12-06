use Mix.Config

config :tg_bot,
       bot_name: "Instatop",
       bot_username: "InstaToppBot",
       mongo_chats: [
         database: "local",
         host: "localhost",
         port: 27017,
       ]

config :tg_bot, TGBot,
       chats_storage: TGBot.Chats.Storages.Mongo

