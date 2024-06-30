package main

import (
	"main/providers"
	"main/server"
)

func main() {
	providers.InstallServices()
	server.StartServer()
}
