defmodule Utils.Messages.Builder do
  alias Utils.Messages.Message
  @type t :: module
  @callback from_data(data :: map()) :: Message.t
end

