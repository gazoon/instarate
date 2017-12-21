defmodule Utils do
  require Logger
  @spec set_child_id(tuple, any) :: tuple
  def set_child_id(spec, child_id) do
    spec
    |> Tuple.delete_at(0)
    |> Tuple.insert_at(0, child_id)
  end

  @spec timestamp :: integer
  def timestamp do
    :os.system_time(:seconds)
  end

  @spec timestamp_milliseconds :: integer
  def timestamp_milliseconds do
    :os.system_time(:milli_seconds)
  end

  @spec keys_to_atoms(map()) :: map()
  def keys_to_atoms(input_map) do
    input_map
    |> Enum.map(fn ({k, v}) -> {String.to_atom(k), v} end)
    |> Map.new()
  end

  @spec parallelize_tasks([(() -> any)]) :: [any]
  def parallelize_tasks(functions) do
    tasks = functions
            |> Enum.map(&wrap_with_rescue/1)
            |> Enum.map(&Task.async/1)
    tasks
    |> Enum.map(&Task.await/1)
    |> Enum.map(&bagrify/1)
  end

  @spec bagrify(tuple) :: any
  defp bagrify(result)

  defp bagrify({:ok, result}), do: result
  defp bagrify({:error, msg, stacktrace}), do: reraise(msg, stacktrace)

  @spec wrap_with_rescue((() -> any)) :: (() -> any)
  defp wrap_with_rescue(func) do
    fn ->
      try do
        {:ok, func.()}
      rescue
        e -> {:error, e, System.stacktrace()}
      end
    end
  end
end
