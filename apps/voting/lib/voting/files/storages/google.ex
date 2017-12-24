defmodule Voting.Files.Storages.Google do
  alias Voting.Files.Storage
  @behaviour Storage
  @scope  "https://www.googleapis.com/auth/cloud-platform"
  @api_url "https://www.googleapis.com/upload/storage/v1/b/"
  @config Application.get_env(:voting, __MODULE__)
  @bucket @config[:bucket_name]

  @spec upload(String.t, String.t) :: String.t
  def upload(path, content_url) do
    token = case Goth.Token.for_scope(@scope) do
      {:ok, token} -> token
      :error -> raise "Can't get token for #{@scope}"
    end
    query = URI.encode_query(%{uploadType: "media", name: path})
    upload_url = @api_url <> @bucket <> "/o?" <> query
    {body, content_type} = Utils.download_file(content_url)
    case HTTPoison.post!(
           upload_url,
           body,
           [{"Authorization", "#{token.type} #{token.token}"}, {"Content-Type", content_type}]
         ) do
      %HTTPoison.Response{status_code: 200} -> build_url(path)
      %HTTPoison.Response{body: body, status_code: status_code} ->
        raise "Can't upload to GCS #{status_code} #{body}"
    end
  end

  @spec build_url(Strin.t) :: String.t
  def build_url(path) do
    "https://storage.googleapis.com/#{@bucket}/#{path}"
  end
end
