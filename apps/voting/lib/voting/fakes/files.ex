defmodule Voting.Fakes.Files do
  alias Voting.Files.Storage
  @behaviour Storage

  @spec upload(String.t, String.t) :: String.t
  def upload(path, _content_url), do: path

  # not used
  @spec build_url(Strin.t) :: String.t
  def build_url(_path), do: ""
end
