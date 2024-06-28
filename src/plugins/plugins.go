package plugins

import (
	"fmt"
	"log"
	"os"
	"plugin"
	"strings"

	"main/utils"
)

type RabbitPlugin struct {
	plugin   *plugin.Plugin
	metadata RabbitPluginMetadata
	path     string
}

type RabbitPluginMetadata struct {
	Name    string
	Version string
}

var loadedPlugins []RabbitPlugin

func LoadPlugins() {
	fmt.Println("Loading plugins...")
	plugins := findPlugins()
	fmt.Println("Found plugins:", plugins)
	for _, plugin := range plugins {
		fmt.Println("Loading plugin:", plugin)
		loadPlugin(plugin)
	}
}

func loadPlugin(pluginName string) {
	// See if plugin is found inside plu
	exeDir, err := utils.GetExecutableDir()
	if err != nil {
		log.Fatal("Error getting executable directory:", err)
	}

	// Check if plugin exists
	pluginPath := exeDir + "/plugins/" + pluginName
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		fmt.Println("Plugin not found:", pluginName)
		return
	}

	// Load plugin
	fmt.Println("Loading plugin:", pluginName)
	p, err := plugin.Open(pluginPath)
	if err != nil {
		fmt.Println("Error loading plugin:", err)
		return
	}

	// Get plugin metadata
	metadata := getPluginMetadata(p)
	if metadata.Name == "" {
		fmt.Println("Error loading plugin metadata")
		return
	}

	// Create RabbitPlugin
	rabbitPlugin := RabbitPlugin{
		plugin:   p,
		metadata: metadata,
		path:     pluginPath,
	}

	fmt.Println("Loaded plugin:", rabbitPlugin.metadata.Name, rabbitPlugin.metadata.Version)

	// Append to plugins
	loadedPlugins = append(loadedPlugins, rabbitPlugin)

	// Call OnLoad
	symbol, err := p.Lookup("OnLoad")
	if err != nil {
		fmt.Println("Error loading OnLoad:", err)
		return
	}

	onLoad, ok := symbol.(func())
	if !ok {
		fmt.Println("Error loading OnLoad:", err)
		return
	}

	onLoad()
}

func getPluginMetadata(p *plugin.Plugin) RabbitPluginMetadata {
	// Get plugin name
	symbol, err := p.Lookup("Name")
	if err != nil {
		fmt.Println("Error loading plugin name:", err)
		return RabbitPluginMetadata{}
	}

	name, ok := symbol.(*string)
	if !ok {
		fmt.Println("Error loading plugin name:", err)
		return RabbitPluginMetadata{}
	}

	// Get plugin version
	symbol, err = p.Lookup("Version")
	if err != nil {
		fmt.Println("Error loading plugin version:", err)
		return RabbitPluginMetadata{}
	}

	version, ok := symbol.(*string)
	if !ok {
		fmt.Println("Error loading plugin version:", err)
		return RabbitPluginMetadata{}
	}

	return RabbitPluginMetadata{
		Name:    *name,
		Version: *version,
	}
}

func findPlugins() []string {
	exeDir, err := utils.GetExecutableDir()
	if err != nil {
		log.Fatal("Error getting executable directory:", err)
	}

	files, err := os.ReadDir(exeDir + "/plugins")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	var plugins []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			plugins = append(plugins, file.Name())
		}
	}

	return plugins
}
