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
    tasks = Enum.map(functions, &Task.async/1)
    Enum.map(tasks, &Task.await/1)
  end
end
