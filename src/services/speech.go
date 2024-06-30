package services

import (
	"main/config"
	"main/structures"
)

func RunSpeech(client *structures.Client, audio []byte) {
	if config.BaseSpeech == nil {
		return
	}

	config.BaseSpeech.Run(client, audio)
}
