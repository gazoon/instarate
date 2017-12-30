defmodule TGBot do

  require Logger
  alias Utils.Messages.Text, as: TextMessage
  alias Utils.Messages.Callback, as: Callback
  alias Utils.Messages.Task, as: TaskMessage
  alias Utils.Messages.Message
  alias TGBot.{Localization, UserMessage}
  alias TGBot.Chats.Chat
  import TGBot.Processing.{Callbacks, Tasks, Text}
  import Localization, only: [get_default_translation: 1]

  @config Application.get_env(:tg_bot, __MODULE__)
  @chats_storage @config[:chats_storage]
  @messenger @config[:messenger]

  @spec on_message(map()) :: any
  def on_message(message_container) do
    message_type = String.to_atom(message_container["type"])
    message_data = message_container["data"]
    message_info = case message_type do
      :text -> {TextMessage, &on_text_message/2}
      :callback -> {Callback, &on_callback/2}
      :task -> {TaskMessage, &on_task/2}
      _ -> nil
    end
    case message_info do
      {message_cls, handler_func} ->
        message = message_cls.from_data(message_data)
        process_message(message, handler_func)
        Logger.info("Finish message processing")
      nil -> Logger.error("Unsupported message type: #{message_type}")
    end
  end

  @spec initialize_context(integer) :: any
  defp initialize_context(chat_id) do
    request_id = UUID.uuid4()
    Logger.metadata([request_id: request_id, chat_id: chat_id])
  end

  @spec process_message(Message.t, ((Message.t, Chat.t) -> Chat.t)) :: any
  defp process_message(message, handler) do
    chat_id = Message.chat_id(message)
    initialize_context(chat_id)

    try do
      {chat, is_new} = get_or_create_chat(message)
      chat_after_processing = handler.(message, chat)
      if chat_after_processing != chat || is_new do
        Logger.info("Save updated chat info")
        @chats_storage.save(chat_after_processing)
      end
    rescue
      error ->
        Logger.error(Exception.format(:error, error))
        send_error_msg(chat_id)
        reraise error, System.stacktrace()
    catch
      :exit, {{error, stack}, _from} ->
        Logger.error(Exception.format(:error, error, stack))
        send_error_msg(chat_id)
        reraise error, stack
      :exit, error ->
        Logger.error(Exception.format(:error, "Exit signal #{inspect error}"))
        send_error_msg(chat_id)
        reraise inspect(error), System.stacktrace()
    end
  end

  @spec send_error_msg(integer) :: any
  defp send_error_msg(chat_id) do
    @messenger.send_text(
      chat_id,
      get_default_translation("unknown_error"),
      static_keyboard: :remove
    )
  end

  @spec get_or_create_chat(Message.t) :: {Chat.t, boolean}
  defp get_or_create_chat(message) do
    chat_id = Message.chat_id(message)
    case @chats_storage.get(chat_id) do
      nil ->
        unless UserMessage.impl_for(message) do
          raise "Chat #{chat_id} not found, and it's not a user message #{inspect message}"
        end
        members_count = @messenger.get_chat_members_number(chat_id) - 1
        {Chat.new(chat_id, UserMessage.is_group_chat(message), members_count), true}
      chat -> {chat, false}
    end
  end
end
