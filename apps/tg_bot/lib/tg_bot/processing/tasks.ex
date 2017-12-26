defmodule TGBot.Processing.Tasks do
  import TGBot.Processing.Common

  require Logger
  alias TGBot.Messages.Task, as: TaskMessage
  alias TGBot.Localization
  alias TGBot.Chats.Chat
  import Localization, only: [get_translation: 2]

  @next_pair_task next_pair_task()
  @daily_activation_task daily_activation_task()

  @spec on_task(TaskMessage.t, Chat.t) :: Chat.t
  def on_task(task, chat)  do
    handler = case task.name do
      @next_pair_task -> &handle_next_pair_task/2
      @daily_activation_task -> &handle_daily_activation_task/2
      _ -> nil
    end
    if handler do
      Logger.info("Handle #{task.name} task")
      handler.(task, chat)
    else
      Logger.info("Received unknown #{task.name} task")
      chat
    end
  end

  @spec handle_next_pair_task(TaskMessage.t, Chat.t) :: Chat.t
  defp handle_next_pair_task(task, chat) do
    task_match_message_id = task.args.last_match_message_id
    actual_match_message_id = chat.last_match.message_id
    if task_match_message_id == actual_match_message_id do
      send_next_girls_pair(chat)
    else
      Logger.info(
        "Skip next pair task, not invalid match message id: #{task_match_message_id}
         actual one: #{actual_match_message_id}"
      )
      chat
    end
  end

  @spec handle_daily_activation_task(TaskMessage.t, Chat.t) :: Chat.t
  defp handle_daily_activation_task(_task, chat) do
    send_next_girls_pair(chat, message_before: get_translation(chat, "propose_to_vote"))
  end
end
