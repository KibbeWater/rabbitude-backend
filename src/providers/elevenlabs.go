package providers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/api"
	"main/config"
	"main/structures"
	"main/utils"
	"net/http"
	"os"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var elevenlabsProviderName = "elevenlabs"
var elevenlabsProvider = structures.Provider{
	ProviderName: elevenlabsProviderName,
}

var (
	elevenlabsAPIKey  string
	elenvlabsModelID  string
	elevenlabsVoiceID string
)

type ElevenLabsResponse struct {
	Audio     string `json:"audio_base64"`
	Alignment struct {
		Characters            []string  `json:"characters"`
		Character_start_times []float64 `json:"character_start_times_seconds"`
		Character_end_times   []float64 `json:"character_end_times_seconds"`
	} `json:"normalized_alignment"`
}

func RegisterElevenlabs(baseServices *[]structures.BaseService, services *[]structures.Service) {
	elevenlabsSetup()

	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    elevenlabsProvider,
		ServiceType: structures.TTS_SERVICE,
		Run:         elevenlabsTTS,
	})
}

func elevenlabsSetup() {
	cfg := config.GetProviderConfig(elevenlabsProviderName)
	if cfg == nil {
		cfg = &structures.ProviderConfig{
			ProviderName: elevenlabsProviderName,
			Options:      make(map[string]interface{}),
		}
	}

	if _, ok := cfg.Options["api_key"]; !ok {
		cfg.Options["api_key"] = utils.GetSetupValue("Elevenlabs Setup - API Key")
	}
	if _, ok := cfg.Options["model_id"]; !ok {
		cfg.Options["model_id"] = utils.GetSetupValue("Elevenlabs Setup - Model ID")
	}
	if _, ok := cfg.Options["voice_id"]; !ok {
		cfg.Options["voice_id"] = utils.GetSetupValue("Elevenlabs Setup - Voice ID")
	}

	elevenlabsAPIKey = cfg.Options["api_key"].(string)
	elenvlabsModelID = cfg.Options["model_id"].(string)
	elevenlabsVoiceID = cfg.Options["voice_id"].(string)

	config.SaveProviderConfig(cfg)
}

func elevenlabsTTS(client *structures.Client, text string) {
	// Create a HTTP post request to the ElevenLabs TTS API
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s/stream/with-timestamps", elevenlabsVoiceID)
	headers := map[string]string{
		"Content-Type": "application/json",
		"xi-api-key":   elevenlabsAPIKey,
	}

	body := map[string]interface{}{
		"text":     text,
		"model_id": elenvlabsModelID,
		"voice_settings": map[string]interface{}{
			"stability":         0.7,
			"similarity_boost":  0.3,
			"use_speaker_boost": true,
		},
	}

	// Marshal the body to JSON without utils
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a http request
	req, err := utils.CreatePostRequest(url, string(jsonBody), headers)
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

	fmt.Println("ElevenLabs response received")

	var responses []ElevenLabsResponse
	lines := strings.Split(string(respBody), "\n")
	for _, line := range lines {
		// if line is not empty, add it to the contentfulLines
		if line != "" {
			fmt.Println("Parsing line of length: ", len(line))
			var response ElevenLabsResponse
			if err := json.Unmarshal([]byte(line), &response); err != nil {
				log.Fatal(err)
			}

			fmt.Println("Appending response")
			responses = append(responses, response)
		}
	}

	fmt.Println("Response count: ", len(responses))

	// For each response, add the audio to the audioBytes
	var audioBytes []byte
	var characters []string
	var characterStartTimes []float64
	var characterEndTimes []float64
	for _, response := range responses {
		audio, err := base64.StdEncoding.DecodeString(response.Audio)
		if err != nil {
			log.Fatal(err)
		}
		audioBytes = append(audioBytes, audio...)

		characters = append(characters, response.Alignment.Characters...)
		characterStartTimes = append(characterStartTimes, response.Alignment.Character_start_times...)
		characterEndTimes = append(characterEndTimes, response.Alignment.Character_end_times...)
	}

	// Create a new characterDurationTimes array with the duration of each character in ms
	var characterStartTimesMs []float64
	var characterEndTimesMs []float64
	for i := 0; i < len(characterStartTimes); i++ {
		characterStartTimesMs = append(characterStartTimesMs, characterStartTimes[i]*1000)
		characterEndTimesMs = append(characterEndTimesMs, characterEndTimes[i]*1000)
	}

	var characterDurationTimes []float64
	for i := 0; i < len(characterStartTimesMs); i++ {
		characterDurationTimes = append(characterDurationTimes, characterEndTimesMs[i]-characterStartTimesMs[i])
	}

	fmt.Println("Performed character duration calculations")

	// Marshal a json of following format {language: "en", chars: characters, char_start_times_ms: characterStartTimesMs, char_durations_ms: characterDurationTimes}
	characterData := map[string]interface{}{
		"language":            "en",
		"chars":               characters,
		"char_start_times_ms": characterStartTimesMs,
		"char_durations_ms":   characterDurationTimes,
	}
	characterDataJson, err := json.Marshal(characterData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Marshalled character data: ", string(characterDataJson))

	// Save audio data to exeDir /tts.wav
	exeDir, err := utils.GetExecutableDir()
	if err != nil {
		log.Fatal(err)
	}

	audioPath := fmt.Sprintf("%s/tts.mp3", exeDir)
	if err := os.WriteFile(audioPath, audioBytes, 0644); err != nil {
		log.Fatal(err)
	}

	newAudioPath := fmt.Sprintf("%s/tts.wav", exeDir)

	// ffmpeg command: ffmpeg -i tts.mp3 -acodec pcm_s16le -ar 16000 -ac 1 tts.wav
	err = ffmpeg.Input(audioPath).Output(newAudioPath, ffmpeg.KwArgs{
		"acodec": "pcm_s16le",
		"ar":     44100,
		"ac":     1,
		"b:a":    "256k",
	}).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		os.Remove(audioPath)
		os.Remove(newAudioPath)
		log.Fatal(err)
	}

	// Read the new audio file
	audioBytes, err = os.ReadFile(newAudioPath)
	if err != nil {
		os.Remove(audioPath)
		os.Remove(newAudioPath)
		log.Fatal(err)
	}

	// Remove the audio files
	os.Remove(audioPath)
	os.Remove(newAudioPath)

	api.SendAudioResponse(client, audioBytes, string(characterDataJson))
}