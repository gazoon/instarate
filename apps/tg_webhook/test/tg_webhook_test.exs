defmodule TgWebhookTest do
  use ExUnit.Case
  doctest TgWebhook

  test "greets the world" do
    assert TgWebhook.hello() == :world
  end
end
