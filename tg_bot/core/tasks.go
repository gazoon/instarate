package core

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"instarate/tg_bot/messages"
	"instarate/tg_bot/models"
)

func (self *Bot) onNextPairTask(ctx context.Context, chat *models.Chat,
	message *messages.NextPairTask) error {
	if message.LastMatchMessageId != chat.LastMatch.MessageId {
		logger := self.GetLogger(ctx)
		logger.WithFields(log.Fields{
			"task_match_id":   message.LastMatchMessageId,
			"actual_match_id": chat.LastMatch.MessageId,
		}).Info("Next pair task relates to another match, skip")
		return nil
	}
	return self.sendNextGirlsPair(ctx, chat)
}

func (self *Bot) onDailyActivationTask(ctx context.Context, chat *models.Chat,
	message *messages.DailyActivationTask) error {

	text := self.gettext(chat, "propose_to_vote")
	if _, err := self.messenger.SendText(ctx, chat.Id, text); err != nil {
		return err
	}
	return self.sendNextGirlsPair(ctx, chat)
}
