package services

import (
	"fmt"
	"main/api"
	"main/config"
	"main/structures"
	"main/utils"
	"math/rand/v2"
	"strings"
)

var llmAckResponses = []string{
	"Working on it",
	"Give me a moment",
	"Let me think about that",
	"Let me check",
}

var searchAckResponses = []string{
	"Looking it up for you",
	"Searching",
	"Give me a moment",
	"Let me look that up for you",
}

func sendAckResponse(client *structures.Client, responses []string) {
	fmt.Println("Sending ack response, total responses: ", len(responses))
	response := ""
	if len(responses) > 0 {
		randVal := rand.IntN(len(responses))
		response = responses[randVal]

		fmt.Println("Sending ack response id", randVal, "of", len(responses), ":", response)
	}

	if response == "" {
		return
	}

	cachedResponse, err := utils.GetCachedPromptTTS("ack_", response)
	if err != nil {
		fmt.Println("Error getting cached response: ", err)
		return
	}

	if cachedResponse != nil {
		api.SendAudioResponse(client, cachedResponse.Audio, cachedResponse.TextMetadata)
		api.SendTextResponse(client, response)
		return
	}

	if config.BaseTTS == nil {
		fmt.Println("No base TTS provider found")
		return
	}

	// We will ignore the preventDef flag here but add it so we won't get errors
	var preventDef bool
	api.SendTextResponse(client, response)
	audioData, err := config.BaseTTS.Run(client, []byte(response), &preventDef)
	if err != nil {
		fmt.Println("Error running ack TTS: ", err)
		return
	}

	audio, err := utils.ReadAudioReturn(audioData)
	if err != nil {
		fmt.Println("Error reading audio return")
		return
	}

	api.SendAudioResponse(client, audio.Audio, audio.TextMetadata)
	utils.CachePromptTTS("ack_", response, audio.Audio, audio.TextMetadata)
}

// Classifies which Service to be used to fullfill the request
func ClassifyText(client *structures.Client, text string, speechRecognised bool) {
	fmt.Println("Classifying text: ", text)

	fmt.Println("ServiceBase: ", config.ServiceBase)
	if config.ServiceBase == nil {
		fmt.Println("No base provider found")
		return
	}

	var preventDef bool
	ret, err := config.ServiceBase.Run(client, []byte(text), &preventDef)
	if err != nil {
		fmt.Println("Error running base service: ", err)
		return
	}
	classifiedText := strings.ToLower(string(ret))

	fmt.Println("Base service returned: ", classifiedText)

	if preventDef {
		return
	}

	if strings.Contains(classifiedText, "lang") {
		fmt.Println("Running LLM service")
		go sendAckResponse(client, llmAckResponses)
		RunLLM(client, text, speechRecognised)
		return
	}

	if strings.Contains(classifiedText, "search") {
		fmt.Println("Running search service")
		go sendAckResponse(client, searchAckResponses)
		RunSearch(client, text)
		return
	}
}
