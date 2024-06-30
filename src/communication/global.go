package communication

import (
	"encoding/json"
	"fmt"
	"log"
	"main/api"
	"main/structures"
	"strings"
)

func HandleGlobal(req structures.ServiceRequest) {
	// Unmarshal the JSON obj from req.data
	var jsonMap map[string]interface{}
	err := json.Unmarshal(req.Data, &jsonMap)
	if err != nil {
		log.Printf("error %s when parsing json", err)
		return
	}

	for key := range jsonMap["global"].(map[string]interface{}) {
		switch key {
		case "initialize":
			var loginData structures.LoginRequest
			err := json.Unmarshal(req.Data, &loginData)
			if err != nil {
				log.Printf("error %s when parsing json", err)
				return
			}

			// Print the login data
			log.Printf("Login data: %+v", loginData)

			req.Client.Imei = loginData.Global.Initialize.DeviceId

			tokenParts := strings.Split(loginData.Global.Initialize.Token, "+")
			if len(tokenParts) != 2 {
				log.Println("Invalid token")
				return
			}

			req.Client.AccountKey = tokenParts[1]

			req.Client.IsLoggedIn = true

			fmt.Println(req.Client)

			api.SendInitResponse(req.Client)
		}
	}
}
