package structures

type Config struct {
	General GeneralConfig `json:"general"`
}

type GeneralConfig struct {
	BaseProvider   string `json:"base_provider"`
	LLMProvider    string `json:"llm_provider"`
	SpeechProvider string `json:"speech_provider"`
	TTSProvider    string `json:"tts_provider"`
	SearchProvider string `json:"search_provider"`
}

type ServiceConfig struct {
	ServiceType string `json:"service_type"`
	Provider    string `json:"provider"`
}
