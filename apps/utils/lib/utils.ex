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
    tasks = Enum.map(functions, &Task.Supervisor.async_nolink(:dynamic_tasks_supervisor, &1))
    Enum.map(tasks, &Task.await/1)
  end

  def tasks_supervisor_spec, do: {Task.Supervisor, name: :dynamic_tasks_supervisor}


  @spec download_file(String.t, Keyword.t) :: {binary, String.t}
  def download_file(url, opts \\ []) do
    pool_name = Keyword.get(opts, :http_pool, :file_downloading)
    case HTTPoison.get!(
           url,
           [],
           hackney: [
             pool: pool_name
           ]
         ) do
      %HTTPoison.Response{body: body, headers: headers, status_code: 200} ->
        {body, Map.get(Map.new(headers), "Content-Type", "")}
      %HTTPoison.Response{status_code: status_code} ->
        raise "Can't download file #{url} #{status_code}"
    end
  end
end
