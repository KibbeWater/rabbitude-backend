package structures

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn       *websocket.Conn
	Imei       string
	AccountKey string
}

type ServiceRequest struct {
	Client *Client
	Data   []byte
}

type LoginRequest struct {
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

type UserTextRequest struct {
	Kernel struct {
		UserText struct {
			Text string `json:"text"`
		} `json:"userText"`
	} `json:"kernel"`
}

type AssistantResponse struct {
	Kernel struct {
		Response string `json:"assistantResponse"`
	} `json:"kernel"`
}

type AssistantDeviceResponse struct {
	Kernel struct {
		AssistantResponseDevice struct {
			Text  string `json:"text"`
			Audio string `json:"audio"`
		} `json:"assistantResponseDevice"`
	} `json:"kernel"`
}
