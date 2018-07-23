package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	. "instarate/tg_gateway/config"
	"instarate/tg_gateway/webhook"
	"net/http"
	"net/url"
	"strconv"
)

func main() {
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", Config.BotToken)
	data := url.Values{}
	webhookUrl := Config.PublicUrl + webhook.BuildWebhookPath(Config.BotToken)
	data.Set("url", webhookUrl)
	data.Set("max_connections", strconv.Itoa(100))
	log.WithFields(log.Fields{"token": Config.BotToken, "webhook_url": webhookUrl}).Info("Set webhook")
	resp, err := http.PostForm(apiUrl, data)
	if err != nil {
		panic(errors.Wrap(err, "post form request"))
	}
	log.WithField("status_code", resp.StatusCode).Info()
}
