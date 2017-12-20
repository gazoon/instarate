defmodule TGBot.Cache.Behaviour do
  @type value :: any
  @callback get(key :: String.t) :: value | nil
  @callback set(key :: String.t, value :: value) :: :ok
end
