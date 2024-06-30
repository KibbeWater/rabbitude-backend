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
	"os"
)

const whisperProviderName = "whisper"

var whisperProvider = structures.Provider{
	ProviderName: whisperProviderName,
}

// Config variables
var (
	whisper_model string

	whisper_setup bool
)

func RegisterWhisper(baseServices *[]structures.BaseService, services *[]structures.Service) {
	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    whisperProvider,
		ServiceType: structures.SPEECH_SERVICE,
		Run:         whisperSpeech,
		Setup:       whisperSetup,
	})
}

func whisperSpeech(client *structures.Client, data []byte) {
	filename := utils.GenerateUniqueID() + ".wav"
	exeDir, err := utils.GetExecutableDir()
	if err != nil {
		log.Println("Error getting executable directory:", err)
		return
	}

	filePath := fmt.Sprintf("%s/%s", exeDir, filename)

	// Save the audio to a file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Println("Error writing audio to file:", err)
		return
	}

	body := map[string]interface{}{
		"path": filePath,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		os.Remove(filePath)
		return
	}

	// Make the request
	req, err := utils.CreatePostRequest("http://localhost:8118/api/whisper", string(bodyJSON), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Remove(filePath)
		return
	}

	// Send the request
	httpClient := &http.Client{}

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error sending request: ", err)
		os.Remove(filePath)
	}

	// Check the response
	if resp.StatusCode != 200 {
		fmt.Println("Error running whisper, status code:", resp.StatusCode)
		os.Remove(filePath)
	}

	// Read the response body {"status": "success", "message": "<message>"}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
		os.Remove(filePath)
	}

	// unmarshal the response body
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		log.Fatal("Error unmarshalling response body: ", err)
		os.Remove(filePath)
	}

	if response["status"] != "success" {
		fmt.Println("Error running whisper, message:", response["message"])
		os.Remove(filePath)
	}

	// Delete the audio file
	if err := os.Remove(filename); err != nil {
		log.Println("Error deleting audio file:", err)
	}

	message := response["message"].(string)

	api.SendSpeechRecognised(client, message)
	services.ClassifyText(client, message)
}

func whisperSetup() {
	if whisper_setup {
		return
	}
	whisper_setup = true

	cfg := config.GetProviderConfig(whisperProviderName)
	if cfg == nil {
		cfg = &structures.ProviderConfig{
			ProviderName: whisperProviderName,
			Options:      make(map[string]interface{}),
		}
	}

	if _, ok := cfg.Options["model"]; !ok {
		cfg.Options["model"] = utils.GetSetupValue("Whisper Setup - Model")
	}
	whisper_model = cfg.Options["model"].(string)

	config.SaveProviderConfig(cfg)

	go initializeWhisper()
}

func initializeWhisper() {
	// Make a POST request to localhost:8118/api/init with {"model": whisper_model}
	body := []byte(`{"model": "` + whisper_model + `"}`)

	// Make the request
	req, err := utils.CreatePostRequest("http://localhost:8118/api/init", string(body), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Send the request
	httpClient := &http.Client{}

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error sending request: ", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != 200 {
		fmt.Println("Error initializing whisper, status code:", resp.StatusCode)
	}
}
