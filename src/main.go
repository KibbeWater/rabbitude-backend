package main

import (
	"main/api"
	"main/services"
)

func main() {
	services.InstallServices()
	api.StartServer()
}
