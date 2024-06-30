package config

import "main/structures"

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
