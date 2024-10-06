package multiplayer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"otte_main_backend/src/meta"
	"time"
)

type CreateLobbyResponseDTO struct {
	ID uint32 `json:"id"`
}

type HealthCheckResponseDTO struct {
	Status     bool   `json:"status"`
	LobbyCount uint32 `json:"lobbyCount"`
}

// Returns lobbyID, error
func CreateLobby(ownerID uint32, colonyID uint32, appContext *meta.ApplicationContext) (uint32, error) {
	url := fmt.Sprintf("%s/create-lobby?ownerID=%d&encoding=binary&colonyID=%d", appContext.MultiplayerServerAddress, ownerID, colonyID)
	body, err := makeEmptyPostRequest(url)
	if err != nil {
		return 0, fmt.Errorf("error creating lobby: %v", err)
	}
	return body.ID, nil
}

func CheckConnection(appContext *meta.ApplicationContext) (*HealthCheckResponseDTO, error) {
	url := fmt.Sprintf("%s/health", appContext.MultiplayerServerAddress)
	resp, err := makeGetRequest[HealthCheckResponseDTO](url)

	if err != nil {
		return nil, fmt.Errorf("error checking connection: %v", err)
	}
	return resp, nil
}

func makeGetRequest[T any](url string) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error performing multiplayer backend healthcheck: %v", err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing multiplayer backend healthcheck: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from multiplayer backend: %d", resp.StatusCode)
	}
	var dest T
	err = json.NewDecoder(resp.Body).Decode(&dest)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %v", err)
	}

	return &dest, nil
}

func makeEmptyPostRequest(url string) (*CreateLobbyResponseDTO, error) {
	// Create a new POST request with no body
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response body into the ResponseData struct
	var responseData CreateLobbyResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %v", err)
	}

	return &responseData, nil
}
