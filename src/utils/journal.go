package utils

import (
	"encoding/json"
	"main/structures"
)

func SendTextEntry(client *structures.Client, text string, response string, recognized bool) {
	// Json marshal the data
	textData := structures.JournalTextEntry{
		VoiceMode: recognized,
		Response:  response,
	}
	jsonData, err := json.Marshal(textData)
	if err != nil {
		return
	}

	entry := structures.JournalEntry{
		Type:  structures.JOURNAL_TEXT_ENTRY,
		Title: text,
		Data:  string(jsonData),
	}

	SendJournalEntry(*client, entry)
}
