package communication

import (
	"fmt"
	"main/structures"
)

func HandleAudioData(client *structures.Client, audio []byte) {
	if client == nil || !client.IsLoggedIn {
		fmt.Println("Client not logged in")
		return
	}

	fmt.Println("Received audio data from client")
	bufLen := len(client.AudioBuf)
	client.AudioBuf = append(client.AudioBuf, audio)
	fmt.Println("Audio buffer length", bufLen, "->", len(client.AudioBuf))
}
