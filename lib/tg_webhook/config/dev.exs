use Mix.Config

config :tg_webhook,
       known_bots: %{
         "480997285:AAEwT3739sBnTz0RSqhEz8TNh4wvJUuqn20" => "insta_queue_gazon"
       },
       port: 8080,
       public_url: "https://dev-dot-instarate-tg-webhook.appspot.com"

