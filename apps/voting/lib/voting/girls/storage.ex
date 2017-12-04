defmodule Voting.Girls.Storage do

  alias Voting.Girls.Girl

  @callback get_top(number :: integer, offset :: integer) :: [Girl.t]
  @callback get_random_pair() :: {Girl.t, Girl.t}
  @callback get_girl(username :: String.t) :: {:ok, Girl.t} | {:error, Stringt.t}
  @callback get_higher_ratings_number(rating :: integer) :: integer
  @callback update_girl(girl :: Girl.t) :: Girl.t
  @callback add_girl(girl :: Girl.t) :: {:ok, Girl.t} | {:error, String.t}
end
