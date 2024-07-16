package services

import (
	"main/api"
	"main/config"
	"main/structures"
)

func RunSpeech(client *structures.Client, audio []byte) {
	if config.BaseSpeech == nil {
		return
	}

	var preventDef bool
	ret, err := config.BaseSpeech.Run(client, audio, &preventDef)
	if err != nil {
		return
	}
	speechTranscript := string(ret)

	if preventDef {
		return
	}

	api.SendSpeechRecognised(client, speechTranscript)
	ClassifyText(client, speechTranscript, true)
}
