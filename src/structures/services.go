package structures

type ServiceFunction func(client *Client, data []byte)
type BaseServiceFunction func(client *Client, data []byte, preventDef *bool) ([]byte, error)
type ServiceSetupFunction func()
type Service struct {
	Provider Provider

	Name        string
	Description string

	Run   ServiceFunction
	Setup ServiceSetupFunction
}

const (
	BASE_SERVICE       int = 0
	SPEECH_SERVICE     int = 1
	TTS_SERVICE        int = 2
	LLM_SERVICE        int = 3
	SEARCH_SERVICE     int = 4
	GENERATIVE_SERVICE int = 5
)

type BaseService struct {
	Provider Provider

	ServiceType int

	Run   BaseServiceFunction
	Setup ServiceSetupFunction
}

type Provider struct {
	ProviderName string
}

type ProviderAudioResponse struct {
	TextMetadata string
	Audio        []byte
}
