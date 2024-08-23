package config

import "os"

//  To keep the configuration details
type Config struct {
	StorageType string `json:"storage_type"`
	Port        string `json:"port"`
	// add more configuration in future
}

// LoadConfig - confiuration from enviornment varibale
func LoadConfig() (*Config, error) {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "in-memory" // Default to in-memory
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	return &Config{
		StorageType: storageType,
		Port:        port,
	}, nil
}
