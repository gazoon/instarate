package core

import (
	"context"
	"fmt"
	"instarate/tg_bot/concatenation"
)

func (self *Bot) getMatchPhoto(ctx context.Context, leftPhotoId, rightPhotoId string) (string, error) {
	key := buildMatchPhotoCacheKey(leftPhotoId, rightPhotoId)
	value, _, err := self.cache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (self *Bot) setMatchPhoto(ctx context.Context, leftPhotoId, rightPhotoId, tgFileId string) error {
	key := buildMatchPhotoCacheKey(leftPhotoId, rightPhotoId)
	err := self.cache.Set(ctx, key, tgFileId)
	return err
}

func buildMatchPhotoCacheKey(leftPhotoId, rightPhotoId string) string {
	return fmt.Sprintf("%s: %s | %s", concatenation.Version, leftPhotoId, rightPhotoId)
}
