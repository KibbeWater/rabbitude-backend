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

	// Check if request header "device-id" exists, if so, set variable "imei" to the value of the header
	imei := r.Header.Get("deviceId")
	fmt.Println("Received connection from device: ", imei)

	// TODO: We need a DB to store persistent client data
	client := structures.Client{
		Conn:            ws,
		Imei:            imei,
		AccountKey:      "",
		DashboardAPIURL: "",
		IsLoggedIn:      false,
	}

	fmt.Printf("New connection, created client: %+v\n", client)

	// Append the client to the list of clients
	communication.Clients = append(communication.Clients, client)

	// Get the pointer to the client
	clientPtr := &communication.Clients[len(communication.Clients)-1]

	// Our client is now mutable, yippie!
	go handleMessages(clientPtr)
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
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(portNumber), nil))
}
