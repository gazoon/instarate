defmodule Utils do
  def set_child_id(spec, child_id) do
    spec
    |> Tuple.delete_at(0)
    |> Tuple.insert_at(0, child_id)
  end
end
