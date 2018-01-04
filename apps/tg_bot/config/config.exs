use Mix.Config

config :tg_bot,
       bot_name: "InstaRateLocal",
       bot_username: "InstaRateLocalBot",
       concatenation_service: "http://concatenation-dot-instarate-190012.appspot.com",
       mongo_chats: [
         database: "local",
         seeds: ["localhost:27017"],
       ],
       mongo_scheduler: [
         database: "local",
         seeds: ["localhost:27017"],
       ],
       mongo_cache: [
         database: "local",
         seeds: ["localhost:27017"],
       ]

config :utils,
       mongo_queue: [
         database: "local",
         seeds: ["localhost:27017"],
         collection: "insta_queue",
         max_processing_time: 10000
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
       pictures_concatenator: TGBot.Pictures.Concatenators.HttpService

config :tg_bot, Scheduler.Reader,
       tasks_storage: Scheduler.Impls.Mongo,
       queue: Utils.Queue.Impls.Mongo

config :tg_bot,
       TGBot.QueueReader,
       queue: Utils.Queue.Impls.Mongo,
       fetch_delay: 100

config :tg_bot, TGBot.MatchPhotoCache,
       cache: TGBot.Cache.Impls.Mongo,
       pictures_concatenator: TGBot.Pictures.Concatenators.HttpService

config :nadia, token: "480997285:AAEwT3739sBnTz0RSqhEz8TNh4wvJUuqn20"

config :logger, :console,
       metadata: [:request_id, :chat_id]

config :tg_bot, TGBot.Localization,
       disable_translation: false

config :voting,
       Voting,
       girls_storage: Voting.Competitors.Storages.Mongo,
       voters_storage: Voting.Voters.Storages.Mongo,
       profiles_storage: Voting.InstagramProfiles.Storages.Mongo

config :voting,
       Voting.Competitors.Model,
       storage: Voting.Competitors.Storages.Mongo

config :voting,
       Voting.InstagramProfiles.Model,
       photos_storage: Voting.Files.Storages.Google


config :voting, Voting.Files.Storages.Google,
       bucket_name: "insta-rate-local"

config :goth,
       json: ~s({"private_key": "-----BEGIN PRIVATE KEY-----\\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCxh2lkkyuRHMkp\\nWI325CDJB/tK0lSwwfyuEskYpa3AvGmZwaXTZWPDMR5jKrn4Xa+q1CYqxA+dfi62\\nNEuph5ExQgKHwSKCnHIpgkWwz68WiarEea4QwdZMAD+DW1P1oJyZGOeq4j3y/ej6\\nqYXN3vm7phPMues94rVx9vKDkYOL+kXoh7YsuytTwhE3TtNRnFIbcoHFw3tYDQTt\\nVGrDhSdTLai0qGl0fnWHQ78Uy+cIzBbMCopPZj+sIAiWDtknC32XztJvDhKs+Qf8\\nu+yKvx5xp/yz0FC+4+YXqIft3rog5pUU/TdQvIsIIZ3+GzaxfRXlEFC28pdU8PPH\\n8tPcD8gTAgMBAAECggEAE73zQnP7TE/fLvHiF6qW2sRAdbmeIWnxJ4p/QnZqNBy/\\nwZF/V1rXFUXJE04VEEGa32xoMilLc0AtAaYfI+MnikrE7UPVCqFLMMKD7X2rAt6c\\nVu0RJlvn8Nu2NG1bkN0jjQLwAFjYesUxu25Oqu445104pnOmbWNguiO0JF1yO4aR\\n1JoWUu/skYcTHfKAwqzzSaJ/wfBPcm0XRr9ynZKqGzQw3cdgsWt7ONu0sXKJJOig\\nw9Mm2BFR95+9R4VkmpV1Vw3sXxByGyfVJprtfuu5eLFlpr7gSHBzhkpQxZ3mrBgk\\ndSYQLaHceSPTKEc7t+id5pV5sWOpvTYpghZhIUjWQQKBgQDpMOeWSXSOLk324apT\\ntPPO4Fk7MMFTG166YyKCZ9KPo+Y2J/S0SPf0w0GRQzmrtTStnlQRlp7LUaza5hf7\\nj1OrB86wUGnNh97F59ntUX3UXzvKEbVG3rdUR0Xh1eR4wJKBAwpdUb3GGs1SMhhO\\nU2YmiVhBiTvjfW/Jk1U2qe/WpwKBgQDC5Lr/rRo7Tcp9orkOvr2o75W21IOncDos\\nnjEESGdLCzqotyA7Lwgqn8u/w0RNDLb3NpFI1keGBhZF4VQEF4FPLCCcLw1ToYWX\\n2C6RkJsy3713wbNhSzaL9P4UiG2ZixVk0Kpx5II1OcDlXX2RkMzTSQ3B8NPRUS9v\\nmzAwLP9ctQKBgCJs2NOD9pQC2/mlaOrAnPmefy3QzcmCEeL8PZz4zMntzU+TaHmx\\nCAH3TVevj/T8ZFF5PTn3fWvQm+8Y5tN5XkWyel3nESeJdmbLCo4RaL0Qbhgvsw5K\\nVNA5UWS8meUFsNsg4sfSCG0VidgnkxUFFOB9iCzsoI299+HPQVY3kjOjAoGBAIl1\\n6q8K2MWbSIb0jrHntr3AvkgF/BXNAjsWGFx89N3pPaZiA0m7End93aeTgkkV/ra+\\ntho5iJjvEiaXlzqLmZjN9vIx/aRO+Hrw72ecJtrrFCezZ2HoOsDcO5kf4K27e4tv\\n4cgS9AO2iGc+WaKiDtW3YQy5X6zzJhIB0ysnkbVBAoGBANrsTLz85Hlo72BWOfjj\\nI7xA17ollLah1F/OkFbrCWAMXG8oxTTUp1C+5XhmqHnLzCtTnPGiE2V/Vp+UTQhK\\nShW+kb0azYR6K84HmK9VECfKtYASHjVmN9Ci7ZjYTA9WIjK+YAkCJ8963AjL94wD\\ntt+eRnSq0v1smvXgtZ6oklN5\\n-----END PRIVATE KEY-----\\n", "client_id": "103474766428157250352", "token_uri": "https://accounts.google.com/o/oauth2/token", "project_id": "instarate-190012", "private_key_id": "2ad1a78758bf6d1805a168afd942ad7fee8cd4cf", "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs", "type": "service_account", "auth_uri": "https://accounts.google.com/o/oauth2/auth", "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/general%40instarate-190012.iam.gserviceaccount.com", "client_email": "general@instarate-190012.iam.gserviceaccount.com"})

config :voting,
       mongo_girls: [
         database: "local",
         seeds: ["localhost:27017"],
       ],
       mongo_voters: [
         database: "local",
         seeds: ["localhost:27017"],
       ],
       mongo_profiles: [
         database: "local",
         seeds: ["localhost:27017"],
       ]

import_config "#{Mix.env}.exs"
