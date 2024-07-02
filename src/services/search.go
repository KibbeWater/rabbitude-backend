package services

import (
	"main/api"
	"main/config"
	"main/structures"
)

func RunSearch(client *structures.Client, text string) {
	if config.BaseSearch == nil {
		return
	}

	var preventDef bool
	ret, err := config.BaseSearch.Run(client, []byte(text), &preventDef)
	if err != nil {
		return
	}
	searchResult := string(ret)

	if preventDef {
		return
	}

	api.SendTextResponse(client, searchResult)
	RunTTS(client, searchResult)
}
