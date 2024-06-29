package api

import (
	"encoding/json"
	"log"
	"main/structures"
	"net/http"

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
		Conn: ws,
	}

	go handleMessages(client)
}

func handleMessages(client structures.Client) {
	for {
		msgType, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("error %s when reading message", err)
			break
		}

		if msgType == websocket.BinaryMessage {
			// We are streaming audio data, ignore this for now
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

func parseJson(msg []byte, client structures.Client) {
	// parse the json object into a map
	var jsonMap map[string]interface{}
	err := json.Unmarshal(msg, &jsonMap)
	if err != nil {
		log.Printf("error %s when parsing json", err)
		return
	}

	req := structures.ServiceRequest{
		Client: &client,
		Data:   msg,
	}

	// is the top-level object key "global" exists?
	if _, ok := jsonMap["global"]; ok {
		HandleGlobal(req)
		return
	}

	// is the top-level object key "kernel" exists?
	if _, ok := jsonMap["kernel"]; ok {
		HandleKernel(req)
		return
	}
}

func StartServer() {
	portNumber := "8080"

	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	http.Handle("/", webSocketHandler)
	log.Print("Starting server on port " + portNumber)
	log.Fatal(http.ListenAndServe("localhost:"+portNumber, nil))
}
