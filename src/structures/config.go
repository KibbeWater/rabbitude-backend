package structures

type Config struct {
	Version   int              `json:"version"`
	General   GeneralConfig    `json:"general"`
	Services  []ServiceConfig  `json:"services"`
	Providers []ProviderConfig `json:"providers"`
}

type GeneralConfig struct {
	Port int `json:"port"`

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

type ProviderConfig struct {
	ProviderName string                 `json:"provider_name"`
	Options      map[string]interface{} `json:"options"`
}
