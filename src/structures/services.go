package structures

type ServiceFunction func(client *Client, data []byte)
type ServiceSetupFunction func()
type Service struct {
	Provider Provider

	Name        string
	Description string

	Run   ServiceFunction
	Setup ServiceSetupFunction
}

const (
	BASE_SERVICE   int = 0
	SPEECH_SERVICE int = 1
	TTS_SERVICE    int = 2
	LLM_SERVICE    int = 3
	SEARCH_SERVICE int = 4
)

type BaseService struct {
	Provider Provider

	ServiceType int

	Run   ServiceFunction
	Setup ServiceSetupFunction
}

type Provider struct {
	ProviderName string
}
