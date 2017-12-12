defmodule TGBot.Queue.Consumer do
  @callback get_next :: {any, String.t}
  @callback finish_processing(processing_id :: String.t) :: any
end
