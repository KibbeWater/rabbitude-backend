package api

import (
	"main/structures"
)

func SendTextResponse(client *structures.Client, text string) {
	response := map[string]interface{}{
		"kernel": map[string]interface{}{
			"assistantResponse": text,
		},
	}

	client.Conn.WriteJSON(response)
}

func SendAudioResponse(client *structures.Client, audio []byte, text string) {
	response := map[string]interface{}{
		"kernel": map[string]interface{}{
			"assistantResponseDevice": map[string]interface{}{
				"text":  text,
				"audio": audio,
			},
		},
	}

	client.Conn.WriteJSON(response)
}
