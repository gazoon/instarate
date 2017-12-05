defmodule TGBot.Chats.Chat do
  alias TGBot.Chats.Chat
  defmodule Match do
    @type t :: %Match{
                 message_id: integer,
                 left_girl: String.t,
                 right_girl: String.t,
                 shown_at: integer
               }
    defstruct message_id: nil, left_girl: nil, right_girl: nil, shown_at: nil
  end

  @type t :: %Chat{chat_id: integer, current_top_offset: integer, last_match: Match.t}
  defstruct chat_id: nil, current_top_offset: 0, last_match: nil
end
