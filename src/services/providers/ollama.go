package providers

import (
	"main/api"
	"main/structures"
)

var ollamaProvider = structures.Provider{
	ProviderName: "ollama",
}

// Returns (BaseServices, CustomServices)
func RegisterOllama(baseServices *[]structures.BaseService, services *[]structures.Service) {
	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    ollamaProvider,
		ServiceType: structures.BASE_SERVICE,
		Run:         runBase,
	})

	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    ollamaProvider,
		ServiceType: structures.LLM_SERVICE,
		Run:         runLLM,
	})

	*services = append(*services, structures.Service{
		Provider:    ollamaProvider,
		Name:        "uber",
		Description: "Orders a taxi given a location",
		Run:         runUber,
	})
}

func runBase(client *structures.Client, text string) {
	api.SendTextResponse(client, text)
}

func runLLM(client *structures.Client, text string) {
	api.SendTextResponse(client, text)
}

func runUber(client *structures.Client, text string) {
	api.SendTextResponse(client, "No taxi for you lol")
}
