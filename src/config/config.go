package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var cache = map[string]string{}

type ServiceConstants struct {
	ServicePort int64
	DDH         string
	AuthToken   string
}

func NewServiceConstants() (*ServiceConstants, error) {
	servicePortStr, err := LoudGet("SERVICE_PORT")
	if err != nil {
		return nil, err
	}
	servicePortInt, err := strconv.ParseInt(servicePortStr, 10, 32)
	if err != nil {
		return nil, err
	}

	authTokenName, err := LoudGet("AUTH_TOKEN_NAME")
	if err != nil {
		return nil, err
	}

	return &ServiceConstants{
		ServicePort: servicePortInt,
		DDH:         GetOr("DEFAULT_DEBUG_HEADER", "DEFAULT-DEBUG-HEADER"),
		AuthToken:   authTokenName,
	}, nil
}

// Checks the exec args and loads the first one found. Ignores the rest.
// Accepts flags: "--dev", "--prod" (prod is handled by the docker-compose file currentyly)
func DetectAndApplyENV() error {
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return fmt.Errorf("[config] Error getting working directory: %s", wdErr.Error())
	}
	log.Printf("[config] Working directory: %s\n", wd)
	args := os.Args
	var envErr error
	for _, arg := range args[1:] {
		if arg == "--dev" {
			log.Println("[config] --dev flag found, loading dev config")
			envErr = LoadDevConfig()
		}
		if arg == "--prod" {
			log.Println("[config] --prod flag found, loading prod config")
			envErr = LoadProdConfig()
		}

		if envErr != nil {
			return envErr
		}
	}

	return nil
}

// Overwrites any env variables currently set in environment
func LoadDevConfig() error {
	err := LoadCustomConfig("dev.env")
	if err != nil {
		return err
	}
	return loadDevCredentials()
}
func loadDevCredentials() error {
	return LoadCustomConfig("dev.credentials")
}

// Overwrites any env variables currently set in environment
func LoadProdConfig() error {
	return LoadCustomConfig("prod.env")
}

// Overwrites any env variables currently set in environment
func LoadCustomConfig(nameOfFile string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("[config] Error getting working directory: %s", err.Error())
	}
	qualifiedPath := filepath.Join(wd, nameOfFile)
	err = godotenv.Overload(qualifiedPath) // Overwrites all env variables with the ones in the .env file
	if err != nil {
		return fmt.Errorf("[config] Error loading .env file into environment: %s", err.Error())
	}
	return nil
}

// LoudGet func to get env value, will return an error on empty string
// The value of the key will be trimmed/stripped/whitespace removed
func LoudGet(key string) (string, error) {
	if val, exists := cache[key]; exists {
		return val, nil
	}
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return "", fmt.Errorf("[config] Tried to access env: %s but failed", key)
	}
	cache[key] = val
	return val, nil
}

// Get func to get env value, will log on error but return the empty value
// The value of the key will be trimmed/stripped/whitespace removed
func Get(key string) string {
	val, err := LoudGet(key)
	if err != nil {
		log.Println(err.Error())
	}
	return val
}

// The value of the key will be trimmed/stripped/whitespace removed
func GetOr(key string, defaultValue string) string {
	val, err := LoudGet(key)
	if err != nil {
		return defaultValue
	}
	return val
}
