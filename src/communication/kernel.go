package communication

import (
	"encoding/json"
	"fmt"
	"log"
	"main/services"
	"main/structures"
	"main/utils"
)

type ServiceRequest = structures.ServiceRequest
type UserTextRequest = structures.UserTextRequest

func HandleKernel(req ServiceRequest) {
	// Unmarshal the JSON obj from req.data
	var jsonMap map[string]interface{}
	err := json.Unmarshal(req.Data, &jsonMap)
	if err != nil {
		log.Printf("error %s when parsing json", err)
		return
	}

	for key := range jsonMap["kernel"].(map[string]interface{}) {
		switch key {
		case "userText":
			var userTextReq UserTextRequest
			err := json.Unmarshal(req.Data, &userTextReq)
			if err != nil {
				log.Printf("error %s when parsing json", err)
				return
			}

			log.Printf("User text: %s", userTextReq.Kernel.UserText.Text)

			go services.ClassifyText(req.Client, userTextReq.Kernel.UserText.Text, false)
		case "voiceActivity":
			var voiceActivityReq structures.VoiceActivityRequest
			err := json.Unmarshal(req.Data, &voiceActivityReq)
			if err != nil {
				log.Printf("error %s when parsing json", err)
				return
			}

			voiceState := voiceActivityReq.Kernel.VoiceActivity.State
			if voiceState != structures.VOICE_ACTIVITY_PRESSED &&
				voiceState != structures.VOICE_ACTIVITY_RELEASED &&
				voiceState != structures.VOICE_ACTIVITY_INACTIVE {
				log.Println("Invalid voice activity state")
				return
			}

			if voiceState == structures.VOICE_ACTIVITY_PRESSED {
				fmt.Println("Clearing audio buffer, voice activity pressed")
				req.Client.AudioBuf = [][]byte{}
			}

			if voiceState == structures.VOICE_ACTIVITY_RELEASED {
				fmt.Println("Sending audio response")
				go runAudioService(req.Client)
			}
		}
	}
}

func runAudioService(client *structures.Client) {
	audioBuf := utils.MergeAudioBuffer(client.AudioBuf)
	services.RunSpeech(client, audioBuf)
}
