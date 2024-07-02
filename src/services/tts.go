package services

import (
	"fmt"
	"main/api"
	"main/config"
	"main/structures"
	"main/utils"
)

// TODO: Implement stricter text metadata restrictions, currently it will ignore the format which is required on the client
func RunTTS(client *structures.Client, text string) {
	fmt.Println("Running TTS service")
	fmt.Println("BaseTTS: ", config.BaseTTS)
	fmt.Println("ServiceBase: ", config.ServiceBase)

	if config.BaseTTS == nil {
		fmt.Println("No TTS provider found")
		return
	}

	fmt.Println("Running TTS on text: ", text)
	var preventDef bool
	audioInfo, err := config.BaseTTS.Run(client, []byte(text), &preventDef)
	if err != nil {
		fmt.Println("Error running TTS service: ", err)
		return
	}

	if preventDef {
		return
	}

	data, err := utils.ReadAudioReturn(audioInfo)
	if err != nil {
		fmt.Println("Error reading audio return")
		return
	}

	// Send the response back to the client
	api.SendAudioResponse(client, data.Audio, data.TextMetadata)
}
