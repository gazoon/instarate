defmodule Utils.Queue.Producer do
  @callback put(chat_id :: integer, message :: any, opts :: Keyword.t) :: any
end
