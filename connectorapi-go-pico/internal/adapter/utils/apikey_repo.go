package utils

import (
	"connectorapi-go/pkg/config"
)

// Validation of API keys based on configuration
type APIKeyRepository struct {
	keys map[string]*config.APIKey
}

// // New repository and pre-loads keys into map
// func NewAPIKeyRepository(apiKeys []config.APIKey) *APIKeyRepository {
// 	keyMap := make(map[string]*config.APIKey)
// 	for i := range apiKeys {
// 		keyMap[apiKeys[i].Key] = &apiKeys[i]
// 	}
// 	return &APIKeyRepository{keys: keyMap}
// }

func NewAPIKeyRepository(apiKeys []config.APIKey) *APIKeyRepository {
	keyMap := make(map[string]*config.APIKey)
	for i := range apiKeys {
		for _, k := range apiKeys[i].Key { // loop ทุก key ใน []string
			keyMap[k] = &apiKeys[i]
		}
	}
	return &APIKeyRepository{keys: keyMap}
}

// Validate checks if an API key is valid, active, and has permission
func (r *APIKeyRepository) Validate(apiKey, method, path string) bool {
	clientKey, exists := r.keys[apiKey]
	if !exists || clientKey.Status != "active" {
		return false
	}

	// Check if the key has permission for the specific METHOD:PATH
	routeKey := method + ":" + path
	for _, p := range clientKey.Permissions {
		if p == routeKey {
			return true
		}
	}

	return false
}
