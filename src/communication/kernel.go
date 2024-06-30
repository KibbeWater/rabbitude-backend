package communication

import (
	"encoding/json"
	"fmt"
	"log"
	"main/api"
	"main/services"
	"main/structures"
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

			go services.ClassifyText(req.Client, userTextReq.Kernel.UserText.Text)
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

			fmt.Println(voiceState)

			if voiceState == structures.VOICE_ACTIVITY_RELEASED {
				fmt.Println("Sending audio response")
				api.SendAudioResponse(req.Client, []byte(`123`), "response")
			}
		}
	}
}