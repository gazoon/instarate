defmodule Voting.Files.Storage do
  @callback upload(path :: String.t, content_url :: String.t) :: String.t
  @callback build_url(path :: String.t) :: String.t
end
