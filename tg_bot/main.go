package main

import (
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/localization"
	"path"
)

func main() {
	localesDir := path.Join(utils.GetCurrentFileDir(), "locales")
	lm, err := localization.NewManager(localesDir)
	if err != nil {
		panic(err)
	}
	println(lm.GettextD("ru", "messages", "propose_to_vote"))
}
