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

const (
	VOICE_ACTIVITY_PRESSED  string = "pttButtonPressed"
	VOICE_ACTIVITY_RELEASED string = "pttButtonReleased"
	VOICE_ACTIVITY_INACTIVE string = "inactive"
)

type VoiceActivityRequest struct {
	Kernel struct {
		VoiceActivity struct {
			State string `json:"state"`
		} `json:"voiceActivity"`
	} `json:"kernel"`
}
