package api

import (
	"encoding/json"
	"log"
)

type UserTextRequest struct {
	Kernel struct {
		UserText struct {
			Text string `json:"text"`
		} `json:"userText"`
	} `json:"kernel"`
}

func HandleKernel(req ServiceRequest) {
	// Unmarshal the JSON obj from req.data
	var jsonMap map[string]interface{}
	err := json.Unmarshal(req.data, &jsonMap)
	if err != nil {
		log.Printf("error %s when parsing json", err)
		return
	}

	for key := range jsonMap["kernel"].(map[string]interface{}) {
		switch key {
		case "userText":
			var userTextReq UserTextRequest
			err := json.Unmarshal(req.data, &userTextReq)
			if err != nil {
				log.Printf("error %s when parsing json", err)
				return
			}

			log.Printf("User text: %s", userTextReq.Kernel.UserText.Text)
		}
	}
}
