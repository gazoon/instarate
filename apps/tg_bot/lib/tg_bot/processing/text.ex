defmodule TGBot.Processing.Text do
  import TGBot.Processing.Common
  require Logger
  alias TGBot.Messages.Text, as: TextMessage
  alias TGBot.Localization
  alias TGBot.Chats.Chat
  import Localization, only: [get_translation: 3, get_translation: 2]

  @start_cmd "start"
  @add_girl_cmd "addGirl"
  @get_top_cmd "showTop"
  @next_top_cmd "Next girl"
  @get_girl_info_cmd "girlInfo"
  @help_cmd "help"
  @chat_settings_cmd "chatSettings"
  @left_vote_cmd "left"
  @right_vote_cmd "right"
  @global_competition_cmd "globalCompetition"
  @celebrities_competition_cmd "celebritiesCompetition"
  @models_competition_cmd "modelsCompetition"
  @normal_competition_cmd "normalCompetition"
  @enable_daily_activation_cmd "enableNotification"
  @disable_daily_activation_cmd "disableNotification"
  @set_voting_timeout_cmd "votingTimeout"
  @delete_girls_cmd "deleteGirls"
  @set_russian_cmd "setRussian"
  @set_english_cmd "setEnglish"

  @min_voting_timeout 5
  @session_duration_seconds 1200

  @config Application.get_env(:tg_bot, __MODULE__)
  @scheduler @config[:scheduler]
  @messenger @config[:messenger]
  @admins @config[:admins]

  @spec on_text_message(TextMessage.t, Chat.t) :: Chat.t
  def on_text_message(message, chat) do
    if TextMessage.appeal_to_bot?(message) || TextMessage.reply_to_bot?(message)
       || !message.is_group_chat do
      process_text_message(message, chat)
    else
      Logger.info("Skip message #{inspect message} it's not an appeal, reply to the bot or private")
      chat
    end
  end

  @spec process_text_message(TextMessage.t, Chat.t) :: Chat.t
  defp process_text_message(message, chat) do
    Logger.info("Process text message #{inspect message}")
    commands = [
      {@start_cmd, &handle_start_cmd/2},
      {@add_girl_cmd, &handle_add_girl_cmd/2},
      {@get_top_cmd, &handle_get_top_cmd/2},
      {@next_top_cmd, &handle_next_top_cmd/2},
      {@get_girl_info_cmd, &handle_get_girl_info_cmd/2},
      {@help_cmd, &handle_help_cmd/2},
      {@left_vote_cmd, &handle_left_vote_cmd/2},
      {@right_vote_cmd, &handle_right_vote_cmd/2},
      {@global_competition_cmd, &handle_global_competition_cmd/2},
      {@celebrities_competition_cmd, &handle_celebrities_competition_cmd/2},
      {@models_competition_cmd, &handle_models_competition_cmd/2},
      {@normal_competition_cmd, &handle_normal_competition_cmd/2},
      {@enable_daily_activation_cmd, &handle_enable_activation_cmd/2},
      {@disable_daily_activation_cmd, &handle_disable_activation_cmd/2},
      {@set_voting_timeout_cmd, &handle_set_voting_timeout_cmd/2},
      {@delete_girls_cmd, &handle_delete_girls_cmd/2},
      {@set_russian_cmd, &handle_set_russian_cmd/2},
      {@set_english_cmd, &handle_set_english_cmd/2},
      {@chat_settings_cmd, &handle_chat_settings_cmd/2},
    ]
    message_text = message.text_lowercase
    command = commands
              |> Enum.find(
                   fn ({cmd_name, _}) -> String.contains?(message_text, String.downcase(cmd_name))
                   end
                 )
    case command do
      {command_name, handler} -> Logger.info("Handle #{command_name} command")
                                 handler.(message, chat)
      nil ->
        Logger.info("Message #{message_text} doesn't contain commands, handle as regular message")
        handle_regular_message(message, chat)
    end
  end

  @spec handle_regular_message(TextMessage.t, Chat.t) :: Chat.t
  defp handle_regular_message(message, chat) do
    if TextMessage.appeal_to_bot?(message) do
      @messenger.send_text(chat.chat_id, get_translation(chat, "dont_get_you"))
    else
      photo_links = String.split(message.text, "\n")
      case photo_links do
        [photo_link] -> add_single_girl(chat, photo_link)
        _ ->
          functions = for photo_link <- photo_links do
            fn ->
              Voting.add_girl(photo_link)
            end
          end
          Utils.parallelize_tasks(functions)
          @messenger.send_text(chat.chat_id, "All girls were processed")
      end
    end
    chat
  end

  @spec handle_start_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_start_cmd(_message, chat) do
    try_to_send_next_pair(chat)
  end

  @spec handle_delete_girls_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_delete_girls_cmd(message, chat) do
    girl_uris = TextMessage.get_command_args(message)
    if Enum.member?(@admins, message.user.id) do
      Voting.delete_girls(girl_uris)
      @messenger.send_text(message.chat_id, "Girls were deleted")
    else
      Logger.warn("Non-admin user #{message.user.id} tried to delete girls")
    end
    chat
  end

  @spec handle_add_girl_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_add_girl_cmd(message, chat) do
    photo_link = TextMessage.get_command_arg(message)
    if photo_link do
      add_single_girl(chat, photo_link)
    else
      @messenger.send_text(message.chat_id, get_translation(chat, "add_girl_no_link"))
    end
    chat
  end

  @spec handle_get_top_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_get_top_cmd(message, chat) do
    optional_start_position = TextMessage.get_command_arg(message) || ""
    offset = case Integer.parse(optional_start_position) do
      {start_position, ""} when start_position > 0 -> start_position - 1
      _ -> 0
    end
    send_girl_from_top(chat, offset)
  end

  @spec handle_enable_activation_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_enable_activation_cmd(_message, chat) do
    @messenger.send_text(chat.chat_id, get_translation(chat, "daily_notification_enabled"))
    %Chat{chat | self_activation_allowed: true}
  end

  @spec handle_disable_activation_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_disable_activation_cmd(_message, chat) do
    @messenger.send_text(chat.chat_id, get_translation(chat, "daily_notification_disabled"))
    @scheduler.delete_task(chat.chat_id, daily_activation_task())
    %Chat{chat | self_activation_allowed: false}
  end

  @spec handle_next_top_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_next_top_cmd(_message, chat) do
    next_offset = chat.current_top_offset + 1
    send_girl_from_top(chat, next_offset)
  end

  @spec handle_left_vote_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_left_vote_cmd(message, chat)  do
    if chat.last_match do
      winer_username = chat.last_match.left_girl
      loser_username = chat.last_match.right_girl
      process_vote_message(message, chat, winer_username, loser_username)
    else
      chat
    end
  end

  @spec handle_right_vote_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_right_vote_cmd(message, chat)  do
    if chat.last_match do
      loser_username = chat.last_match.left_girl
      winer_username = chat.last_match.right_girl
      process_vote_message(message, chat, winer_username, loser_username)
    else
      chat
    end
  end

  @spec handle_set_voting_timeout_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_set_voting_timeout_cmd(message, chat) do
    arg = TextMessage.get_command_arg(message) || ""
    case Integer.parse(arg) do
      {timeout, ""} when @min_voting_timeout <= timeout and timeout < @session_duration_seconds ->
        text = get_translation(chat, "voting_timeout_is_set", timeout: timeout)
        @messenger.send_text(chat.chat_id, text)
        %Chat{chat | voting_timeout: timeout}
      {_, ""} ->
        lower_bound = @min_voting_timeout - 1
        upper_bound = div @session_duration_seconds, 60
        @messenger.send_text(
          chat.chat_id,
          get_translation(
            chat,
            "set_voting_timeout_out_of_range",
            lower_bound: lower_bound,
            upper_bound: upper_bound
          )
        )
        chat
      _ ->
        @messenger.send_text(chat.chat_id, get_translation(chat, "set_voting_timeout_no_number"))
        chat
    end
  end

  @spec handle_help_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_help_cmd(message, chat) do
    @messenger.send_text(message.chat_id, get_translation(chat, "help_message"))
    chat
  end

  @spec handle_chat_settings_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_chat_settings_cmd(message, chat) do
    @messenger.send_text(message.chat_id, get_translation(chat, "chat_settings_commands"))
    chat
  end

  @spec handle_global_competition_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_global_competition_cmd(_message, chat) do
    @messenger.send_text(
      chat.chat_id,
      get_translation(chat, "global_competition_enabled"),
      static_keyboard: :remove
    )
    %Chat{chat | competition: Voting.global_competition()}
  end

  @spec handle_celebrities_competition_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_celebrities_competition_cmd(_message, chat) do
    @messenger.send_text(
      chat.chat_id,
      get_translation(chat, "celebrities_competition_enabled"),
      static_keyboard: :remove
    )
    %Chat{chat | competition: Voting.celebrities_competition()}
  end

  @spec handle_models_competition_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_models_competition_cmd(_message, chat) do
    @messenger.send_text(
      chat.chat_id,
      get_translation(chat, "models_competition_enabled"),
      static_keyboard: :remove
    )
    %Chat{chat | competition: Voting.models_competition()}
  end

  @spec handle_normal_competition_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_normal_competition_cmd(_message, chat) do
    @messenger.send_text(
      chat.chat_id,
      get_translation(chat, "normal_competition_enabled"),
      static_keyboard: :remove
    )
    %Chat{chat | competition: Voting.normal_competition()}
  end

  @spec handle_set_russian_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_set_russian_cmd(_message, chat) do
    chat = %Chat{chat | language: Localization.russian_lang()}
    @messenger.send_text(chat.chat_id, get_translation(chat, "switch_to_language"))
    chat
  end

  @spec handle_set_english_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_set_english_cmd(_message, chat) do
    chat = %Chat{chat | language: Localization.english_lang()}
    @messenger.send_text(chat.chat_id, get_translation(chat, "switch_to_language"))
    chat
  end

  @spec handle_get_girl_info_cmd(TextMessage.t, Chat.t) :: Chat.t
  defp handle_get_girl_info_cmd(message, chat) do
    girl_link = TextMessage.get_command_arg(message)
    if girl_link do
      case Voting.get_girl(chat.competition, girl_link) do
        {:ok, girl} -> display_girl_info(chat, girl)
        {:error, error_msg} -> @messenger.send_text(message.chat_id, error_msg)
      end
    else
      @messenger.send_text(message.chat_id, get_translation(chat, "get_girl_no_username"))
    end
    chat
  end

  @spec process_vote_message(TextMessage.t, Chat.t, String.t, String.t) :: Chat.t
  defp process_vote_message(message, chat, winner_username, loser_username) do
    voters_group_id = build_voters_group_id(message.chat_id)
    voter_id = build_voter_id(message.user)
    case Voting.vote(
           chat.competition,
           voters_group_id,
           voter_id,
           winner_username,
           loser_username
         ) do
      {:ok, _} -> try_to_send_next_pair(chat)
      {:error, error} ->
        Logger.warn("Can't vote by message: #{error}")
        chat
    end
  end
end
