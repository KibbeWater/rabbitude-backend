package services

import (
	"fmt"
	"main/config"
	"main/structures"
)

// Classifies which Service to be used to fullfill the request
func ClassifyText(client *structures.Client, text string) {
	fmt.Println("Classifying text: ", text)

	fmt.Println("ServiceBase: ", config.ServiceBase)
	if config.ServiceBase == nil {
		fmt.Println("No base provider found")
		return
	}
	config.ServiceBase.Run(client, text)
}
