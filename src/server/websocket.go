package server

import (
	"encoding/json"
	"fmt"
	"log"
	"main/communication"
	"main/config"
	"main/structures"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	client := structures.Client{
		Conn:       ws,
		Imei:       "",
		AccountKey: "",
		IsLoggedIn: false,
	}

	go handleMessages(&client)
}

func handleMessages(client *structures.Client) {
	for {
		msgType, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("error %s when reading message", err)
			break
		}

		fmt.Println("Received message: ", msgType)

		if msgType == websocket.BinaryMessage {
			fmt.Println("Received binary frame")
			go communication.HandleAudioData(client, msg)
			continue
		}

		if msgType != websocket.TextMessage {
			log.Println("Unhandled message type")
			continue
		}

		// Parse the JSON message and assign an appropiate service
		parseJson(msg, client)
	}
}

func parseJson(msg []byte, client *structures.Client) {
	// parse the json object into a map
	var jsonMap map[string]interface{}
	err := json.Unmarshal(msg, &jsonMap)
	if err != nil {
		log.Printf("error %s when parsing json", err)
		return
	}

	req := structures.ServiceRequest{
		Client: client,
		Data:   msg,
	}

	// is the top-level object key "global" exists?
	if _, ok := jsonMap["global"]; ok {
		go communication.HandleGlobal(req)
		return
	}

	// is the top-level object key "kernel" exists?
	if _, ok := jsonMap["kernel"]; ok {
		go communication.HandleKernel(req)
		return
	}
}

func StartServer() {
	configData := config.GetConfig()
	if configData.General.Port == 0 {
		configData.General.Port = 8080
		config.SaveConfig()
	}

	portNumber := configData.General.Port

	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	http.Handle("/", webSocketHandler)
	log.Print("Starting server on port " + strconv.Itoa(portNumber))
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(portNumber), nil))
}
