package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

type Client struct {
	conn       *websocket.Conn
	imei       string
	accountKey string
}

type ServiceRequest struct {
	client *Client
	data   []byte
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	client := Client{
		conn: ws,
	}

	go handleMessages(client)
}

func handleMessages(client Client) {
	for {
		msgType, msg, err := client.conn.ReadMessage()
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

func parseJson(msg []byte, client Client) {
	// parse the json object into a map
	var jsonMap map[string]interface{}
	err := json.Unmarshal(msg, &jsonMap)
	if err != nil {
		log.Printf("error %s when parsing json", err)
		return
	}

	req := ServiceRequest{
		client: &client,
		data:   msg,
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
