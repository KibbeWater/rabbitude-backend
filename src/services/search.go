package services

import (
	"main/config"
	"main/structures"
)

func RunSearch(client *structures.Client, text string) {
	if config.BaseSearch == nil {
		return
	}

	config.BaseSearch.Run(client, []byte(text))
}
