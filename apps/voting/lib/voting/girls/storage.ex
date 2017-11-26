defmodule Voting.Girls.Storage do
  @moduledoc false
  alias Voting.Girls.Girl

  @callback get_top(number :: integer) :: [Girl.t]
  @callback get_random_pair() :: {Girl.t, Girl.t}
  @callback get_girl(username :: String.t) :: Girl.t | nil
  @callback update_girl(girl :: Girl.t) :: Girl.t
  @callback add_girl(girl :: Girl.t) :: {:ok, Girl.t} | {:error, String.t}
end
