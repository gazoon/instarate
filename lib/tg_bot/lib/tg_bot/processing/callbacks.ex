defmodule TGBot.Processing.Callbacks do
  import TGBot.Processing.Common

  require Logger
  alias Utils.Messages.Callback, as: Callback
  alias TGBot.Localization
  alias TGBot.Chats.Chat
  import Localization, only: [get_translation: 3, get_translation: 2]
  @vote_callback "vt"
  @get_top_callback "top"

  @usernames_separator "|"

  @config Application.get_env(:tg_bot, __MODULE__)
  @messenger @config[:messenger]

  @spec on_callback(Callback.t, Chat.t) :: Chat.t
  def on_callback(message, chat) do
    Logger.info("Process callback #{inspect message}")
    callback_name = Callback.get_name(message)
    handler = case callback_name do
      @vote_callback -> &handle_vote_callback/2
      @get_top_callback -> &handle_get_top_callback/2
      _ -> nil
    end
    if handler do
      Logger.info("Handle #{callback_name} callback")
      handler.(message, chat)
    else
      Logger.warn("Unknown callback: #{callback_name}")
      chat
    end
  end

  @spec handle_vote_callback(Callback.t, Chat.t) :: Chat.t
  defp handle_vote_callback(message, chat) do
    callback_args = Callback.get_args(message)
    [winner_username, loser_username] = String.split(callback_args, @usernames_separator)
    voters_group_id = build_voters_group_id(message.chat_id)
    voter_id = build_voter_id(message.user)
    case Voting.vote(
           chat.competition,
           voters_group_id,
           voter_id,
           winner_username,
           loser_username
         ) do
      :ok ->
        @messenger.send_notification(
          message.callback_id,
          get_translation(chat, "success_vote", username: winner_username)
        )
        if chat.last_match.message_id == message.parent_msg_id do
          try_to_send_next_pair(chat)
        else
          chat
        end
      {:error, error} ->
        @messenger.send_notification(message.callback_id, get_translation(chat, "already_voted"))
        Logger.warn("Can't vote by callback: #{error}")
        chat
    end
  end

  @spec handle_get_top_callback(Callback.t, Chat.t) :: Chat.t
  defp handle_get_top_callback(_message, chat) do
    #    @messenger.delete_attached_keyboard(message.chat_id, message.parent_msg_id)
    #    callback_args = Callback.get_args(message)
    #    girl_offset = case Integer.parse(callback_args) do
    #      {offset, ""} -> offset
    #      _ -> raise "Non-int arg for get top callback: #{callback_args}"
    #    end
    #    if chat.current_top_offset == girl_offset - 1 do
    #      chat = send_girl_from_top(chat, girl_offset)
    #      @messenger.answer_callback(message.callback_id)
    #      chat
    #    else
    #      @messenger.send_notification(
    #        message.callback_id,
    #        "Please, continue from the most recent girl."
    #      )
    #      chat
    #    end
    chat
  end
end
