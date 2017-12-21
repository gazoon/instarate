defmodule Utils.Meter do
  defmacro __using__(_opts) do
    quote do
      require Logger
      import Utils.Meter
    end
  end

  defmacro measure([metric_name: metric_name], do: do_block) do
    quote do
      metric_name = unquote(metric_name)
      start_time = Utils.timestamp_milliseconds()
      try do
        result = unquote(do_block)
        finish_time = Utils.timestamp_milliseconds()
        Logger.info("Action #{metric_name} took: #{finish_time - start_time}")
        result
      rescue
        error ->
          finish_time = Utils.timestamp_milliseconds()
          Logger.warn("Action #{metric_name} failed, took: #{finish_time - start_time}")
          reraise error, System.stacktrace()
      end
    end
  end
end
