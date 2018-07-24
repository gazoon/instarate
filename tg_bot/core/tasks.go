package core

import (
	"context"
	"instarate/tg_bot/messages"
	"instarate/tg_bot/models"
)

func (self *Bot) onNextPairTask(ctx context.Context, chat *models.Chat,
	message *messages.NextPairTask) error {

	logger := self.GetLogger(ctx).WithField("task_match_id", message.LastMatchMessageId)
	if chat.LastMatch == nil {
		logger.Info("Next pair task received, but chat doesn't have last match, skip")
		return nil
	}
	if chat.LastMatch.MessageId != message.LastMatchMessageId {
		logger.WithField("actual_match_id", chat.LastMatch.MessageId).
			Info("Next pair task relates to another match, skip")
		return nil
	}
	// This task fires when the chat is ready for the next pair.
	// No need to do the checks.
	// Even more trySendNextGirlsPair can reject the sending,
	// because of the clock drift in distributed systems.
	return self.sendNextGirlsPair(ctx, chat)
}

func (self *Bot) onDailyActivationTask(ctx context.Context, chat *models.Chat,
	message *messages.DailyActivationTask) error {

	text := self.gettext(chat, "propose_to_vote")
	if _, err := self.messenger.SendText(ctx, chat.Id, text); err != nil {
		return err
	}
	// Daily activation task fires after 24 hours inactivity in the chat,
	// so actually we can call sendNextGirlsPair directly and it will be safe.
	// But it costs nothing to do extra checks.
	return self.trySendNextGirlsPair(ctx, chat)
}
