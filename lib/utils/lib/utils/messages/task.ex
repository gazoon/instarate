defmodule Utils.Messages.Task do


  @behaviour Utils.Messages.Builder

  alias Utils.Messages.Task
  @type t :: %Task{
               chat_id: integer,
               name: atom,
               args: map(),
               do_at: integer,
             }
  defstruct chat_id: nil,
            name: nil,
            args: nil,
            do_at: nil

  def type, do: :task

  @spec from_data(map()) :: Task.t
  def from_data(data) do
    data = Utils.keys_to_atoms(data)
    data = %{data | name: String.to_atom(data.name)}

    data = if data.args, do: %{data | args: Utils.keys_to_atoms(data.args)}, else: data
    struct(Task, data)
  end

  @spec new(integer, integer, atom, Keyword.t) :: Task.t
  def new(chat_id, do_at, name, opts \\ []) do
    args = Keyword.get(opts, :args)
    %Task{chat_id: chat_id, name: name, do_at: do_at, args: args}
  end
end
