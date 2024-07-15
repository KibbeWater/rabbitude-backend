package services

import (
	"fmt"
	"main/config"
	"main/structures"
)

func RunGenerative(client *structures.Client, prompt string) string {
	fmt.Println("Running Generative service")
	fmt.Println("BaseLLM: ", config.BaseGenerative)
	fmt.Println("ServiceBase: ", config.ServiceBase)

	if config.BaseGenerative == nil {
		fmt.Println("No LLM provider found")
		return ""
	}

	fmt.Println("Running LLM on text: ", prompt)
	var preventDef bool
	ret, err := config.BaseGenerative.Run(client, []byte(prompt), &preventDef)
	if err != nil {
		fmt.Println("Error running LLM service: ", err)
		return ""
	}
	promptReturn := string(ret)

	fmt.Println("Generative service returned: ", promptReturn)

	if preventDef {
		return ""
	}

	return promptReturn
}
