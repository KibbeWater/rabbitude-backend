package utils

import (
	"main/structures"
	"strconv"
	"time"
)

// The client key is a hashed using sha256 from the string `${imei}_${floor(timeInSeconds/10)}`
func GenerateClientKey(client structures.Client) string {
	time := time.Now().Unix()
	seconds := time / 10
	seconds = int64(seconds)

	return CreateHash([]byte(client.Imei + "_" + strconv.FormatInt(seconds, 10)))
}

func FindClientByKey(key string, clients *[]structures.Client) *structures.Client {
	if key == "" {
		return nil
	}

	for _, client := range *clients {
		clientKey := GenerateClientKey(client)
		if client.AccountKey == "" && clientKey == key {
			return &client
		}
	}
	return nil
}
