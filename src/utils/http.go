package utils

import (
	"encoding/json"
	"fmt"
	"main/structures"
	"net/http"
	"strings"
)

func CreatePostRequest(url string, body string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func SendJournalEntry(client structures.Client, data structures.JournalEntry) {
	if !client.IsLoggedIn || client.DashboardAPIURL == "" {
		fmt.Println("Not sending journal entry, not logged in / no endpoint")
		return
	}

	// Json marshal the data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	// Create a new request
	req, err := CreatePostRequest(client.DashboardAPIURL+"/entry", string(jsonData), map[string]string{
		"Content-Type": "application/json",
		"Device-Id":    client.Imei,
		"Account-Key":  client.AccountKey,
	})
	if err != nil {
		return
	}

	// Send the request
	httpClient := http.Client{}
	httpClient.Do(req)
}
