defmodule Utils.Messages.User do

  alias Utils.Messages.User
  @type t :: %User{id: integer, name: String.t, username: String.t}
  defstruct id: nil, name: nil, username: ""

  @spec from_data(map()) :: TextMessage.t
  def from_data(user_data) do
    user_data = Utils.keys_to_atoms(user_data)
    struct(User, user_data)
  end

  @spec is_bot?(User.t, String.t) :: boolean
  def is_bot?(user, bot_username) do
    user.username == bot_username
  end
end
