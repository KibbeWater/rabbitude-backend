package api

import (
	"encoding/base64"
	"main/structures"
	"time"
)

func SendInitResponse(client *structures.Client) {
	currentTime := time.Now().UTC().Format(time.RFC3339)
	response := map[string]interface{}{
		"global": map[string]interface{}{
			"initialize": map[string]interface{}{
				"currentTime": currentTime,
				"clientIp":    client.Conn.RemoteAddr().String(),
			},
		},
	}

	client.Conn.WriteJSON(response)
}

func SendTextResponse(client *structures.Client, text string) {
	response := map[string]interface{}{
		"kernel": map[string]interface{}{
			"assistantResponse": text,
		},
	}

	client.Conn.WriteJSON(response)
}

func SendAudioResponse(client *structures.Client, audio []byte, text string) {
	// Base64 encode the audio
	audioBase64 := base64.StdEncoding.EncodeToString(audio)

	response := map[string]interface{}{
		"kernel": map[string]interface{}{
			"assistantResponseDevice": map[string]interface{}{
				"text":  text,
				"audio": audioBase64,
			},
		},
	}

	client.Conn.WriteJSON(response)
}
