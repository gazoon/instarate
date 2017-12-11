defmodule TGBot.Messages.Task do


  @behaviour TGBot.MessageBuilder

  alias TGBot.Messages.Task
  @type t :: %Task{
               chat_id: integer,
               name: String.t,
               args: map(),
               do_at: integer,
               unique_mark: String.t
             }
  defstruct chat_id: nil,
            name: nil,
            args: nil,
            do_at: nil,
            unique_mark: nil

  def type, do: :task

  @spec from_data(map()) :: Task.t
  def from_data(data) do
    data = Utils.keys_to_atoms(data)
    data = %{data | name: String.to_atom(data.name)}

    data = if data.args, do: %{data | args: Utils.keys_to_atoms(data.args)}, else: data
    data = if data.unique_mark,
              do: %{data | unique_mark: String.to_atom(data.unique_mark)},
              else: data
    struct(Task, data)
  end


  @spec new(integer, integer, String.t, Keyword.t) :: Task.t
  def new(chat_id, do_at, name, opts \\ []) do
    args = Keyword.get(opts, :args)
    unique_mark = Keyword.get(opts, :unique_mark)
    %Task{chat_id: chat_id, name: name, do_at: do_at, args: args, unique_mark: unique_mark}
  end
end
