package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/api"
	"main/config"
	"main/services"
	"main/structures"
	"main/utils"
	"net/http"
	"strings"
)

var groqProviderName = "groq"
var groqProvider = structures.Provider{
	ProviderName: groqProviderName,
}

var (
	groqAPIKey string

	groqQuickModel  string
	groqDeepModel   string
	groqSpeechModel string

	groq_setup bool
)

func RegisterGroq(baseServices *[]structures.BaseService, services *[]structures.Service) {
	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    groqProvider,
		ServiceType: structures.BASE_SERVICE,
		Run:         groqBase,
		Setup:       groqSetup,
	})

	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    groqProvider,
		ServiceType: structures.LLM_SERVICE,
		Run:         groqLLM,
		Setup:       groqSetup,
	})

	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    groqProvider,
		ServiceType: structures.SPEECH_SERVICE,
		Run:         groqSpeech,
		Setup:       groqSetup,
	})
}

func groqSetup() {
	if groq_setup {
		return
	}
	groq_setup = true

	cfg := config.GetProviderConfig(groqProviderName)
	if cfg == nil {
		cfg = &structures.ProviderConfig{
			ProviderName: groqProviderName,
			Options:      make(map[string]interface{}),
		}
	}

	// Find the options and if they don't exist, run utils.GetSetupValue
	if _, ok := cfg.Options["api_key"]; !ok {
		cfg.Options["api_key"] = utils.GetSetupValue("Groq Setup - API Key")
	}
	if _, ok := cfg.Options["quick_model"]; !ok {
		cfg.Options["quick_model"] = utils.GetSetupValue("Groq Setup - Quick Model")
	}
	if _, ok := cfg.Options["deep_model"]; !ok {
		cfg.Options["deep_model"] = utils.GetSetupValue("Groq Setup - Deep Model")
	}
	if _, ok := cfg.Options["speech_model"]; !ok {
		cfg.Options["speech_model"] = utils.GetSetupValue("Groq Setup - Speech Model")
	}

	// Set the runtime variables
	groqAPIKey = cfg.Options["api_key"].(string)
	groqQuickModel = cfg.Options["quick_model"].(string)
	groqDeepModel = cfg.Options["deep_model"].(string)
	groqSpeechModel = cfg.Options["speech_model"].(string)

	config.SaveProviderConfig(cfg)
}

func makeGroqRequest(model string, system_prompt string, prompt string) string {
	// body: {"model": model, messages: [{role: "system", "content": system_prompt}, {role: "user", "content": prompt}]}
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{"role": "system", "content": system_prompt},
			{"role": "user", "content": prompt},
		},
	}

	// Marshal the body
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + groqAPIKey,
		"Content-Type":  "application/json",
	}

	// Make the request
	req, err := utils.CreatePostRequest("https://api.groq.com/openai/v1/chat/completions", string(bodyJSON), headers)
	if err != nil {
		log.Fatal(err)
	}

	// Send the request
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error running whisper, status code:", resp.StatusCode)
	}

	// Get choices[0].message.content from the response json
	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Fatal("Error unmarshalling response body: ", err)
	}

	// Get the completion
	var completion string
	choices, ok := response["choices"].([]interface{})
	if ok && len(choices) > 0 {
		message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
		if ok {
			content, ok := message["content"].(string)
			if ok {
				completion = content
			}
		}
	}

	return completion
}

func groqLLM(client *structures.Client, data []byte) {
	system_prompt := "You are an AI powered voice assistant called \"Rabbit\", you are a hardware AI assistant running a specialized OS created by \"Rabbitude\" a modding and jailbreaking community looking to improve the Rabbit R1 hardware AI assistant."
	system_prompt = fmt.Sprintf("%s\nSimilar to Siri, you are to respond in a conversational but concise manner. Do not use emojis", system_prompt)

	prompt := string(data)

	completion := makeGroqRequest(groqDeepModel, system_prompt, prompt)

	fmt.Println("Completion: ", completion)

	api.SendTextResponse(client, completion)
	services.RunTTS(client, completion)
}

func groqBase(client *structures.Client, data []byte) {
	system_prompt := utils.BuildClassificationPrompt("Your responses should only contain the name of a given service")

	prompt := string(data)

	completion := makeGroqRequest(groqQuickModel, system_prompt, prompt)

	fmt.Println("Completion: ", completion)

	// Find text between &s and &e
	intention := strings.ToLower(completion)

	// Find the service that matches the intention
	for _, service := range config.CustomServices {
		if strings.Contains(strings.ToLower(service.Name), intention) {
			// Run custom service
			fmt.Println("Running custom service: ", service.Name)
			return
		}
	}

	if strings.Contains(intention, "llm") {
		fmt.Println("Identified LLM service")
		go services.RunLLM(client, prompt)
		return
	}

	if strings.Contains(intention, "search") {
		fmt.Println("Identified Search service")
		go services.RunSearch(client, prompt)
		return
	}
}

func groqSpeech(client *structures.Client, data []byte) {
	fmt.Println("Running Groq speech service")

	// Create a HTTP post request to the Groq TTS API
	url := "https://api.groq.com/openai/v1/audio/transcriptions"
	headers := map[string]string{
		"Authorization": "Bearer " + groqAPIKey,
		"Content-Type":  "multipart/form-data; boundary=grq",
	}

	// Create form-data body
	bodyForm := map[string]interface{}{
		"model":           groqSpeechModel,
		"temperature":     "0",
		"response_format": "json",
		"language":        "en",
	}

	// Create the body string
	var body string
	for key, value := range bodyForm {
		body += fmt.Sprintf("--grq\nContent-Disposition: form-data; name=\"%s\"\n\n%s\n", key, value)
	}
	body += "--grq\nContent-Disposition: form-data; name=\"file\"; filename=\"sample_audio.wav\"\nContent-Type: audio/wav\n\n"
	body += string(data)
	body += "\n--grq--"

	// Create a http request
	req, err := utils.CreatePostRequest(url, body, headers)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{}

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error sending request: ", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	// Get the "text" from the response json
	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Fatal("Error unmarshalling response body: ", err)
	}

	// Get the completion
	var completion string
	text, ok := response["text"].(string)
	if ok {
		completion = text
	}

	api.SendSpeechRecognised(client, completion)
	services.ClassifyText(client, completion)
}
