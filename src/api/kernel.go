package api

import (
	"encoding/json"
	"log"
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
		}
	}
}
