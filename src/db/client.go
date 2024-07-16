package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"main/structures"
	"main/utils"
	"os"
	"path/filepath"
)

type ClientSaveData struct {
	Imei       string
	AccountKey string

	API_URL string
}

func fromClient(client structures.Client) ClientSaveData {
	return ClientSaveData{
		Imei:       client.Imei,
		AccountKey: client.AccountKey,
		API_URL:    client.DashboardAPIURL,
	}
}

func loadToClient(client *structures.Client, data ClientSaveData) {
	// Verify the client data is valid
	if data.Imei == "" {
		return
	}

	if client.Imei != data.Imei || (client.AccountKey != data.AccountKey && data.AccountKey != "") {
		fmt.Printf("ERR! CLIENT DATA MISMATCH, %s != %s || %s != %s)\n", client.Imei, data.Imei, client.AccountKey, data.AccountKey)
		client.UNVERIFIED = true
		return
	}

	client.DashboardAPIURL = data.API_URL
}

func generateClientFilename(client structures.Client) (string, error) {
	filenameHash := utils.CreateHash([]byte(client.Imei))
	if len(filenameHash) < 7 {
		return "", fmt.Errorf("bad hash")
	}

	return fmt.Sprintf("client_%s_%s.json", client.Imei, filenameHash[:7]), nil
}

func getClientDBPath() (string, error) {
	exeDir, err := utils.GetExecutableDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(exeDir, "db", "clients")
	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return cacheDir, nil
}

func SaveClient(client structures.Client) {
	if client.UNVERIFIED {
		fmt.Println("Client is unverified, not saving")
		return
	}

	// Save the client data to the database
	data := fromClient(client)
	filename, err := generateClientFilename(client)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(data)
	if err != nil {
		return
	}

	// Save the data to the database (a file for now)
	saveDir, err := getClientDBPath()
	if err != nil {
		return
	}

	err = os.WriteFile(filepath.Join(saveDir, filename), buffer.Bytes(), os.ModePerm)
	if err != nil {
		return
	}
}

func LoadClient(client *structures.Client) error {
	if client.Imei == "" {
		return fmt.Errorf("client data is invalid")
	}

	if !client.IsLoggedIn {
		return fmt.Errorf("client is not logged in")
	}

	filename, err := generateClientFilename(*client)
	if err != nil {
		return err
	}

	saveDir, err := getClientDBPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filepath.Join(saveDir, filename))
	if err != nil {
		return err
	}

	var clientData ClientSaveData
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(&clientData)
	if err != nil {
		return err
	}

	loadToClient(client, clientData)
	return nil
}
