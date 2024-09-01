package main

import (
	"main/external"
	"main/providers"
	"main/server"
)

func main() {
	external.LoadLibraries()
	defer external.FreeLibraries()

	providers.InstallServices()
	server.StartServer()
}
