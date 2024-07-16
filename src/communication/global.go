package communication

import (
	"encoding/json"
	"fmt"
	"log"
	"main/api"
	"main/db"
	"main/structures"
	"main/utils"
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

			err = db.LoadClient(req.Client)
			if err != nil {
				log.Printf("error %s when loading client data, likely doesn't exist", err)
			}

			// Save the client data
			db.SaveClient(*req.Client)

			fmt.Println(req.Client)

			api.SendInitResponse(req.Client)
		case "web_authenticate":
			var webAuthData structures.WebAuthRequest
			err := json.Unmarshal(req.Data, &webAuthData)
			if err != nil {
				log.Printf("error %s when parsing json", err)
				return
			}

			// Print the web auth data
			log.Printf("Web auth data: %+v", webAuthData)

			// Find the client
			client := utils.FindClientByKey(webAuthData.Global.WebAuthenticate.Key, &Clients)
			if client == nil {
				log.Println("Client not found")
				return
			}

			log.Printf("Client data: %+v", *client)

			// Assign the dashboard API URL
			client.DashboardAPIURL = webAuthData.Global.WebAuthenticate.API_URL

			// Save client data
			fmt.Println(*client)
			db.SaveClient(*client)

			// {global: {web_authenticate: {success: true, error: ""}}}
			response := structures.WebAuthResponse{
				Global: struct {
					WebAuthenticate struct {
						Success bool   `json:"success"`
						Error   string `json:"error"`
					} `json:"web_authenticate"`
				}{
					WebAuthenticate: struct {
						Success bool   `json:"success"`
						Error   string `json:"error"`
					}{
						Success: true,
						Error:   "",
					},
				},
			}
			jsonData, err := json.Marshal(response)
			if err != nil {
				log.Printf("error %s when marshalling json", err)
				return
			}

			// Send the response
			req.Client.Conn.WriteMessage(1, jsonData)
		}
	}
}
