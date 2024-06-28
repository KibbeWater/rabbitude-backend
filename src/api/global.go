package api

import (
	"encoding/json"
	"log"
	"strings"
)

type LoginData struct {
	Global struct {
		Initialize struct {
			DeviceId  string `json:"deviceId"`
			Evaluate  bool   `json:"evaluate"`
			Greet     bool   `json:"greet"`
			Language  string `json:"language"`
			Listening bool   `json:"listening"`
			Location  struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"location"`
			MimeType string `json:"mimeType"`
			TimeZone string `json:"timeZone"`
			Token    string `json:"token"`
		} `json:"initialize"`
	} `json:"global"`
}

func HandleGlobal(req ServiceRequest) {
	// Unmarshal the JSON obj from req.data
	var jsonMap map[string]interface{}
	err := json.Unmarshal(req.data, &jsonMap)
	if err != nil {
		log.Printf("error %s when parsing json", err)
		return
	}

	for key := range jsonMap["global"].(map[string]interface{}) {
		switch key {
		case "initialize":
			var loginData LoginData
			err := json.Unmarshal(req.data, &loginData)
			if err != nil {
				log.Printf("error %s when parsing json", err)
				return
			}

			req.client.imei = loginData.Global.Initialize.DeviceId

			tokenParts := strings.Split(loginData.Global.Initialize.Token, "+")
			if len(tokenParts) != 2 {
				log.Println("Invalid token")
				return
			}

			req.client.accountKey = tokenParts[1]
		}
	}
}
