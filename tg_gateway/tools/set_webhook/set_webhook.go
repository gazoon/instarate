package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/go-utils"
	"instarate/tg_gateway/config"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

func main() {
	conf := &config.Config{}
	configPath := path.Join(utils.GetCurrentFileDir(), "../../config")
	err := utils.ParseConfig(configPath, conf)
	if err != nil {
		panic(err)
	}
	for token := range conf.KnownBots {
		apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token)
		data := url.Values{}
		data.Set("url", conf.PublicUrl+"/update/"+token)
		data.Set("max_connections", strconv.Itoa(100))
		log.WithField("token", token).Info("Set webhook")
		resp, err := http.PostForm(apiUrl, data)
		if err != nil {
			panic(err)
		}
		log.WithField("status_code", resp.StatusCode).Info()
	}
}
