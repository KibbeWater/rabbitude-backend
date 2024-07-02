package services

import (
	"fmt"
	"main/api"
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
	var preventDef bool
	ret, err := config.BaseLLM.Run(client, []byte(text), &preventDef)
	if err != nil {
		fmt.Println("Error running LLM service: ", err)
		return
	}
	promptReturn := string(ret)

	fmt.Println("LLM service returned: ", promptReturn)

	if preventDef {
		return
	}

	// Send the response back to the client
	api.SendTextResponse(client, promptReturn)
	RunTTS(client, promptReturn)
}
