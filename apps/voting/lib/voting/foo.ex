defmodule Voting.Foo do
  @moduledoc false

  use GenServer

  # Callbacks
  def start_link(opts) do
    GenServer.start_link(__MODULE__, %{}, opts)
  end

  def handle_call(:pop, _from, [h | t]) do
    {:reply, h, t}
  end

  def handle_cast({:push, item}, state) do
    {:noreply, [item | state]}
  end
end
