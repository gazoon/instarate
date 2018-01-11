defmodule Mix.Tasks.DeleteWebhook do
  use Mix.Task
  require Logger

  def run(_args) do
    Application.ensure_all_started(:httpoison)
    bot_tokens = for {token, _} <- Application.get_env(:tg_webhook, :known_bots), do: token
    Enum.each(bot_tokens, &delete_webhook/1)
  end

  def delete_webhook(token) do
    Logger.info("Delete webhook for #{token}")
    resp = HTTPoison.post!("https://api.telegram.org/bot#{token}/deleteWebhook", [])
    Logger.info("API response code=#{resp.status_code} body=#{resp.body}")
  end

end
