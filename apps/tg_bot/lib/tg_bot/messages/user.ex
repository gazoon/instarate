defmodule TGBot.Messages.User do

  alias TGBot.Messages.User
  @type t :: %User{id: integer, name: String.t, username: String.t}
  defstruct id: nil, name: nil, username: ""

  @spec from_data(map()) :: TextMessage.t
  def from_data(user_data) do
    struct(User, user_data)
  end

  @spec is_bot?(User.t) :: boolean
  def is_bot?(user) do
    user.username == Application.get_env(:tg_bot, :bot_username)
  end
end

