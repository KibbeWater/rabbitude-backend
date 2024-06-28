package main

import (
	"fmt"
	"main/plugins"
)

func main() {
	fmt.Println("Hello, World!")
	plugins.LoadPlugins()
}
