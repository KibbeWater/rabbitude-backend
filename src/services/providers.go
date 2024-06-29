package services

import (
	"fmt"
	"main/services/providers"
	"main/structures"
	"main/utils"
	"strings"
)

type Service = structures.Service
type BaseService = structures.BaseService

var (
	ServiceBase *BaseService
	BaseSpeech  *BaseService
	BaseTTS     *BaseService
	BaseLLM     *BaseService
	BaseSearch  *BaseService

	CustomServices []Service
)

func InstallServices() {
	config := utils.GetConfig()
	_, baseServices, _ := DiscoverServices()

	var erroredProviders []string

	ServiceBase = findProvider(config.General.BaseProvider, baseServices)
	if ServiceBase == nil {
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
			ServiceBase = &providers[providerIndex]
			println(utils.ConfigData)
			utils.ConfigData.General.BaseProvider = ServiceBase.Provider.ProviderName
		}
	}

	BaseSpeech = findProvider(config.General.SpeechProvider, baseServices)
	if BaseSpeech == nil {
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
			BaseSpeech = &providers[providerIndex]
			utils.ConfigData.General.SpeechProvider = BaseSpeech.Provider.ProviderName
		}
	}

	BaseTTS = findProvider(config.General.TTSProvider, baseServices)
	if BaseTTS == nil {
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
			BaseTTS = &providers[providerIndex]
			utils.ConfigData.General.TTSProvider = BaseTTS.Provider.ProviderName
		}
	}

	BaseLLM = findProvider(config.General.LLMProvider, baseServices)
	if BaseLLM == nil {
		providers := findProvidersByType(structures.LLM_SERVICE, baseServices)

		// Get array of provider names
		var providerNames []string
		for _, provider := range providers {
			providerNames = append(providerNames, provider.Provider.ProviderName)
		}

		// Render the setting page
		providerIndex := utils.RenderSettingPage("Select a LLM provider", providerNames)

		if providerIndex < 0 {
			erroredProviders = append(erroredProviders, "llm")
		} else {
			BaseLLM = &providers[providerIndex]
			utils.ConfigData.General.LLMProvider = BaseLLM.Provider.ProviderName
		}
	}

	BaseSearch = findProvider(config.General.SearchProvider, baseServices)
	if BaseSearch == nil {
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
			BaseSearch = &providers[providerIndex]
			utils.ConfigData.General.SearchProvider = BaseSearch.Provider.ProviderName
		}
	}

	utils.SaveConfig()

	// Find unique service types
	serviceTypes := []string{}
	for _, service := range CustomServices {
		serviceTypes = append(serviceTypes, service.Name)
	}

	// Render the setting page
	for _, serviceType := range serviceTypes {
		providers := findProvidersByTypeName(serviceType, CustomServices)

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
			CustomServices = append(CustomServices, providers[providerIndex])
		}
	}

	if len(erroredProviders) > 0 {
		fmt.Println("Could not find providers for: " + strings.Join(erroredProviders, ", "))
	}
}

func findProvider(providerName string, services []BaseService) *BaseService {
	for _, service := range services {
		if service.Provider.ProviderName == providerName {
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

// Returns a list of unique providers and a list of services.
func DiscoverServices() ([]string, []BaseService, []Service) {
	var baseServices []BaseService
	var services []Service

	providers.RegisterOllama(&baseServices, &services)

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
