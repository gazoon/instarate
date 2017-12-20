defmodule TGBot.Chats.Chat do
  alias TGBot.Chats.Chat
  alias TGBot.Localization
  defmodule Match do
    alias TGBot.Chats.Chat.Match
    @type t :: %Match{
                 message_id: integer,
                 left_girl: String.t,
                 right_girl: String.t,
                 shown_at: integer,
               }
    defstruct message_id: nil, left_girl: nil, right_girl: nil, shown_at: nil

    @spec new(integer, String.t, String.t) :: Match.t
    def new(message_id, left_girl, right_girl) do
      %Match{
        message_id: message_id,
        left_girl: left_girl,
        right_girl: right_girl,
        shown_at: Utils.timestamp_milliseconds()
      }
    end
  end

  @type t :: %Chat{
               chat_id: integer,
               members_number: integer,
               current_top_offset: integer,
               last_match: Match.t,
               created_at: integer,
               competition: String.t,
               self_activation_allowed: boolean,
               voting_timeout: integer,
               language: String.t,
             }
  defstruct chat_id: nil,
            members_number: nil,
            current_top_offset: 0,
            last_match: nil,
            created_at: nil,
            competition: nil,
            self_activation_allowed: true,
            voting_timeout: 5,
            language: nil

  @spec new(integer, integer) :: Chat.t
  def new(chat_id, members_number) do
    %Chat{
      chat_id: chat_id,
      members_number: members_number,
      created_at: Utils.timestamp(),
      competition: Voting.global_competition(),
      language: Localization.english_lang(),
    }
  end
end
