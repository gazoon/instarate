defmodule TgBotTest do
  use ExUnit.Case
  doctest TgBot

  test "greets the world" do
    assert TgBot.hello() == :world
  end
end
