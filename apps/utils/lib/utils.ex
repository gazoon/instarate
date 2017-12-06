defmodule Utils do
  @spec set_child_id(tuple, any) :: tuple
  def set_child_id(spec, child_id) do
    spec
    |> Tuple.delete_at(0)
    |> Tuple.insert_at(0, child_id)
  end

  @spec current_timestamp :: integer
  def current_timestamp do
    DateTime.utc_now()
    |> DateTime.to_unix()
  end
end
