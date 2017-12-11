defmodule Utils.Reader do
  @callback fetch :: any
  @callback process(data :: any) :: any
  defmacro __using__(opts) do

    quote bind_quoted: [
            opts: opts
          ] do

      require Logger
      use GenServer
      @behaviour Utils.Reader

      @fetch_delay Keyword.get(opts, :fetch_delay, 100)
      @workers_number Keyword.get(opts, :workers_number, 1)

      def init(state) do
        Process.flag(:trap_exit, true)
        next_fetch()
        Logger.info("#{__MODULE__} server started")
        {:ok, state}
      end

      defp next_fetch do
        send(self(), :fetch)
      end

      def start_link(opts) do
        GenServer.start_link(__MODULE__, nil, opts)
      end

      def handle_info(_msg, state) do
        data = fetch()
        if data do
          Logger.info("Fetched #{inspect data} start processing")
          process(data)
        else
          Process.sleep(@fetch_delay)
        end
        next_fetch()
        {:noreply, state}
      end

      def children_spec do
        0..@workers_number - 1
        |> Enum.map(
             fn (i) ->
               process_name = String.to_atom("#{__MODULE__}.#{i}")
               Supervisor.child_spec({__MODULE__, name: process_name}, id: {__MODULE__, i})
             end
           )
      end

    end
  end
end

