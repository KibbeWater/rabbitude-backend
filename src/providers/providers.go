package providers

import (
	"fmt"
	"main/config"
	"main/structures"
	"main/utils"
	"strings"
)

type Service = structures.Service
type BaseService = structures.BaseService

// Returns a list of unique providers and a list of services.
func DiscoverServices() ([]string, []BaseService, []Service) {
	var baseServices []BaseService
	var services []Service

	RegisterOllama(&baseServices, &services)
	RegisterElevenlabs(&baseServices, &services)
	RegisterWhisper(&baseServices, &services)
	RegisterGroq(&baseServices, &services)
	RegisterApple(&baseServices, &services)

	// Find all unique providers
	providerMap := make(map[string]bool)
	for _, baseService := range baseServices {
		providerMap[baseService.Provider.ProviderName] = true
	}
	for _, service := range services {
		providerMap[service.Provider.ProviderName] = true
	}

	providers := make([]string, 0, len(providerMap))
	for provider := range providerMap {
		providers = append(providers, provider)
	}

	return providers, baseServices, services
}

func InstallServices() {
	cfg := config.GetConfig()
	_, baseServices, _ := DiscoverServices()

	var erroredProviders []string

	config.ServiceBase = findProvider(cfg.General.BaseProvider, structures.BASE_SERVICE, baseServices)
	if config.ServiceBase == nil {
		providers := findProvidersByType(structures.BASE_SERVICE, baseServices)

		// Get array of provider names
		var providerNames []string
		for _, provider := range providers {
			providerNames = append(providerNames, provider.Provider.ProviderName)
		}

		// Render the setting page
		providerIndex := utils.RenderSettingPage("Select a base provider", providerNames)

		if providerIndex < 0 {
			erroredProviders = append(erroredProviders, "base")
		} else {
			config.ServiceBase = &providers[providerIndex]
			config.ConfigData.General.BaseProvider = config.ServiceBase.Provider.ProviderName
		}
	}

	config.BaseSpeech = findProvider(cfg.General.SpeechProvider, structures.SPEECH_SERVICE, baseServices)
	if config.BaseSpeech == nil {
		providers := findProvidersByType(structures.SPEECH_SERVICE, baseServices)

		// Get array of provider names
		var providerNames []string
		for _, provider := range providers {
			providerNames = append(providerNames, provider.Provider.ProviderName)
		}

		// Render the setting page
		providerIndex := utils.RenderSettingPage("Select a speech provider", providerNames)

		if providerIndex < 0 {
			erroredProviders = append(erroredProviders, "speech")
		} else {
			config.BaseSpeech = &providers[providerIndex]
			config.ConfigData.General.SpeechProvider = config.BaseSpeech.Provider.ProviderName
		}
	}

	config.BaseTTS = findProvider(cfg.General.TTSProvider, structures.TTS_SERVICE, baseServices)
	if config.BaseTTS == nil {
		providers := findProvidersByType(structures.TTS_SERVICE, baseServices)

		// Get array of provider names
		var providerNames []string
		for _, provider := range providers {
			providerNames = append(providerNames, provider.Provider.ProviderName)
		}

		// Render the setting page
		providerIndex := utils.RenderSettingPage("Select a TTS provider", providerNames)

		if providerIndex < 0 {
			erroredProviders = append(erroredProviders, "tts")
		} else {
			config.BaseTTS = &providers[providerIndex]
			config.ConfigData.General.TTSProvider = config.BaseTTS.Provider.ProviderName
		}
	}

	config.BaseLLM = findProvider(cfg.General.LLMProvider, structures.LLM_SERVICE, baseServices)
	if config.BaseLLM == nil {
		providers := findProvidersByType(structures.LLM_SERVICE, baseServices)

		// Get array of provider names
		var providerNames []string
		for _, provider := range providers {
			providerNames = append(providerNames, provider.Provider.ProviderName)
		}

		// Render the setting page
		providerIndex := utils.RenderSettingPage("Select a LLM provider", providerNames)

		if providerIndex < 0 {
			erroredProviders = append(erroredProviders, "lang")
		} else {
			config.BaseLLM = &providers[providerIndex]
			config.ConfigData.General.LLMProvider = config.BaseLLM.Provider.ProviderName
		}
	}

	config.BaseSearch = findProvider(cfg.General.SearchProvider, structures.SEARCH_SERVICE, baseServices)
	if config.BaseSearch == nil {
		providers := findProvidersByType(structures.SEARCH_SERVICE, baseServices)

		// Get array of provider names
		var providerNames []string
		for _, provider := range providers {
			providerNames = append(providerNames, provider.Provider.ProviderName)
		}

		// Render the setting page
		providerIndex := utils.RenderSettingPage("Select a search provider", providerNames)

		if providerIndex < 0 {
			erroredProviders = append(erroredProviders, "search")
		} else {
			config.BaseSearch = &providers[providerIndex]
			config.ConfigData.General.SearchProvider = config.BaseSearch.Provider.ProviderName
		}
	}

	config.BaseGenerative = findProvider(cfg.General.GenerativeProvider, structures.GENERATIVE_SERVICE, baseServices)
	if config.BaseGenerative == nil {
		providers := findProvidersByType(structures.GENERATIVE_SERVICE, baseServices)

		// Get array of provider names
		var providerNames []string
		for _, provider := range providers {
			providerNames = append(providerNames, provider.Provider.ProviderName)
		}

		// Render the setting page
		providerIndex := utils.RenderSettingPage("Select a generative provider", providerNames)

		if providerIndex < 0 {
			erroredProviders = append(erroredProviders, "generative")
		} else {
			config.BaseGenerative = &providers[providerIndex]
			config.ConfigData.General.GenerativeProvider = config.BaseGenerative.Provider.ProviderName
		}
	}

	// Run setups
	if config.ServiceBase != nil {
		config.ServiceBase.Setup()
	}
	if config.BaseSpeech != nil {
		config.BaseSpeech.Setup()
	}
	if config.BaseTTS != nil {
		config.BaseTTS.Setup()
	}
	if config.BaseLLM != nil {
		config.BaseLLM.Setup()
	}
	if config.BaseSearch != nil {
		config.BaseSearch.Setup()
	}
	if config.BaseGenerative != nil {
		config.BaseGenerative.Setup()
	}

	config.SaveConfig()

	// Find unique service types
	serviceTypes := []string{}
	for _, service := range config.CustomServices {
		serviceTypes = append(serviceTypes, service.Name)
	}

	// TODO: Look this over, we do not have custom services implemented yet so can't properly test
	// Render the setting page
	for _, serviceType := range serviceTypes {
		providers := findProvidersByTypeName(serviceType, config.CustomServices)

		// Get array of provider names
		var providerNames []string
		for _, provider := range providers {
			providerNames = append(providerNames, provider.Provider.ProviderName)
		}

		// Render the setting page
		providerIndex := utils.RenderSettingPage("Select a "+serviceType+" provider", providerNames)

		if providerIndex < 0 {
			continue
		} else {
			// Add to custom services
			config.CustomServices = append(config.CustomServices, providers[providerIndex])
		}
	}

	if len(erroredProviders) > 0 {
		fmt.Println("Could not find providers for: " + strings.Join(erroredProviders, ", "))
	}
}

func findProvider(providerName string, serviceType int, services []BaseService) *BaseService {
	for _, service := range services {
		if service.Provider.ProviderName == providerName && service.ServiceType == serviceType {
			return &service
		}
	}

	return nil
}

func findProvidersByType(serviceType int, services []BaseService) []BaseService {
	var providers []BaseService

	for _, service := range services {
		if service.ServiceType == serviceType {
			providers = append(providers, service)
		}
	}

	return providers
}

func findProvidersByTypeName(serviceType string, services []Service) []Service {
	var providers []Service

	for _, service := range services {
		if service.Name == serviceType {
			providers = append(providers, service)
		}
	}

	return providers
}
