package services

import (
	"fmt"
	"main/config"
	"main/structures"
)

func RunTTS(client *structures.Client, text string) {
	fmt.Println("Running TTS service")
	fmt.Println("BaseTTS: ", config.BaseTTS)
	fmt.Println("ServiceBase: ", config.ServiceBase)

	if config.BaseTTS == nil {
		fmt.Println("No TTS provider found")
		return
	}

	fmt.Println("Running TTS on text: ", text)
	config.BaseTTS.Run(client, text)
}
