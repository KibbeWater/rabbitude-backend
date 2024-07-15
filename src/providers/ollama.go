package providers

import (
	"context"
	"fmt"
	"log"
	"main/api"
	"main/config"
	"main/structures"
	"main/utils"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const ollamaProviderName = "ollama"

var ollamaProvider = structures.Provider{
	ProviderName: ollamaProviderName,
}

// Config variables
var (
	ollama_quick_model string
	ollama_deep_model  string
	ollama_setup       bool = false
)

// Returns (BaseServices, CustomServices)
func RegisterOllama(baseServices *[]structures.BaseService, services *[]structures.Service) {
	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    ollamaProvider,
		ServiceType: structures.BASE_SERVICE,
		Run:         ollamaBase,
		Setup:       ollamaSetup,
	})

	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    ollamaProvider,
		ServiceType: structures.LLM_SERVICE,
		Run:         ollamaLLM,
		Setup:       ollamaSetup,
	})

	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    ollamaProvider,
		ServiceType: structures.GENERATIVE_SERVICE,
		Run:         ollamaGen,
		Setup:       ollamaSetup,
	})

	*services = append(*services, structures.Service{
		Provider:    ollamaProvider,
		Name:        "uber",
		Description: "Orders a taxi given a location",
		Run:         runUber,
		Setup:       ollamaSetup,
	})
}

func ollamaSetup() {
	if ollama_setup {
		return
	}
	ollama_setup = true

	cfg := config.GetProviderConfig(ollamaProviderName)
	if cfg == nil {
		cfg = &structures.ProviderConfig{
			ProviderName: ollamaProviderName,
			Options:      make(map[string]interface{}),
		}
	}

	// Find the options and if they don't exist, run utils.GetSetupValue
	if _, ok := cfg.Options["quick_model"]; !ok {
		cfg.Options["quick_model"] = utils.GetSetupValue("Ollama Setup - Quick Model")
	}
	if _, ok := cfg.Options["deep_model"]; !ok {
		cfg.Options["deep_model"] = utils.GetSetupValue("Ollama Setup - Deep Model")
	}

	// Set the runtime variables
	ollama_quick_model = cfg.Options["quick_model"].(string)
	ollama_deep_model = cfg.Options["deep_model"].(string)

	config.SaveProviderConfig(cfg)
}

func ollamaBase(client *structures.Client, data []byte, preventDef *bool) ([]byte, error) {
	text := string(data)

	fmt.Println("Running Ollama base service")
	llm, err := ollama.New(ollama.WithModel(ollama_quick_model))
	if err != nil {
		log.Fatal(err)
	}

	query := utils.BuildClassificationPrompt(utils.SequenceReturnPrompt("&s", "&e"))
	query = fmt.Sprintf("%s\nWhat intention does this query have: %s", query, text)

	fmt.Println("Query: ", query)

	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Completion: ", completion)

	// Find text between &s and &e
	intention, found := utils.FindSubstring(completion, "&s", "&e")
	if !found {
		api.SendTextResponse(client, "I'm not sure what you're asking")
		return nil, nil
	}

	intention = strings.ToLower(intention)

	// Find the service that matches the intention
	for _, service := range config.CustomServices {
		if strings.Contains(strings.ToLower(service.Name), intention) {
			// Run custom service
			fmt.Println("Running custom service: ", service.Name)
			return nil, fmt.Errorf("feature not implemented")
		}
	}

	return []byte(text), nil
}

func ollamaLLM(client *structures.Client, data []byte, preventDef *bool) ([]byte, error) {
	text := string(data)

	fmt.Println("Running Ollama LLM service")
	llm, err := ollama.New(ollama.WithModel(ollama_deep_model))
	if err != nil {
		log.Fatal(err)
	}

	query := "You are an AI powered voice assistant called \"Rabbit\", you are a hardware AI assistant running a specialized OS created by \"Rabbitude\" a modding and jailbreaking community looking to improve the Rabbit R1 hardware AI assistant."
	query = fmt.Sprintf("%s\nSimilar to Siri, you are to respond in a conversational but concise manner. Do not use emojis", query)
	query = fmt.Sprintf("%s\nQuery: %s", query, text)

	fmt.Println("Query: ", query)

	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Completion: ", completion)

	return []byte(completion), nil
}

func ollamaGen(client *structures.Client, data []byte, preventDef *bool) ([]byte, error) {
	text := string(data)

	fmt.Println("Running Ollama Generative service")
	llm, err := ollama.New(ollama.WithModel(ollama_deep_model))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Completion: ", completion)

	return []byte(completion), nil
}

func runUber(client *structures.Client, data []byte) {
	api.SendTextResponse(client, "No taxi for you lol")
}
