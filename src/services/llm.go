package services

import (
	"fmt"
	"main/config"
	"main/structures"
)

func RunLLM(client *structures.Client, text string) {
	fmt.Println("Running LLM service")
	fmt.Println("BaseLLM: ", config.BaseLLM)
	fmt.Println("ServiceBase: ", config.ServiceBase)

	if config.BaseLLM == nil {
		fmt.Println("No LLM provider found")
		return
	}

	fmt.Println("Running LLM on text: ", text)
	config.BaseLLM.Run(client, []byte(text))
}
