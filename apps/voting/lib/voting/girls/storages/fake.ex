defmodule Voting.Girls.Storages.Fake do
  @moduledoc false
  use Agent

  def start_link(_opt) do
    Agent.start_link(fn -> %{} end, name: __MODULE__)
  end

  def add_girl(_girl) do
  end


end
