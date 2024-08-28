package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)

// Config holds the configuration for the SFTP daemon.
var RootConfig Config

// LoadRootConfig parses the config file and returns the configuration.
func LoadRootConfig() *Config {
	configPath := flag.String("config-path", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		log.Println("Please specify the config-path parameter")
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("Loading root configuration...")
	parseRootConfig(*configPath)
	return &RootConfig
}

// parseRootConfig reads the JSON config file and populates the RootConfig variable.
func parseRootConfig(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file at %s: %v", path, err)
	}

	if err := json.Unmarshal(content, &RootConfig); err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}

	// Create user home directories
	for _, acc := range RootConfig.Accounts {
		if _, err := os.Open(RootConfig.BasePath); os.IsNotExist(err) {
			log.Printf("Creating base path directory: %s", RootConfig.BasePath)
			os.Mkdir(RootConfig.BasePath, 0777)
		}
		homeDir := filepath.Join(RootConfig.BasePath, acc.Username)
		if _, err := os.Open(homeDir); os.IsNotExist(err) {
			log.Printf("Creating home directory for user %s: %s", acc.Username, homeDir)
			os.Mkdir(homeDir, 0777)
		}
	}
}
