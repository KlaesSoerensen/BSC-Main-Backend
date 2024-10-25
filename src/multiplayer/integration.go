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
	Message    string `json:"message"`
}

type ClientStateResponseDTO struct {
	LastKnownPosition uint32 `json:"lastKnownPosition"`
	MsOfLastMessage   uint64 `json:"msOfLastMessage"`
}

type ClientResponseDTO struct {
	ID    uint32                 `json:"id"`
	IGN   string                 `json:"IGN"`
	Type  string                 `json:"type"`
	State ClientStateResponseDTO `json:"state"`
}

type LobbyStateResponseDTO struct {
	ColonyID uint32              `json:"colonyID"`
	Closing  bool                `json:"closing"`
	Phase    uint32              `json:"phase"`
	Encoding string              `json:"encoding"`
	Clients  []ClientResponseDTO `json:"clients"`
}

// Returns lobbyID, error
func CreateLobby(ownerID uint32, colonyID uint32, appContext *meta.ApplicationContext) (uint32, error) {
	url := fmt.Sprintf("%s/create-lobby?ownerID=%d&encoding=binary&colonyID=%d", appContext.InternalMultiplayerServerAddress, ownerID, colonyID)
	body, err := makeEmptyPostRequest(url)
	if err != nil {
		return 0, fmt.Errorf("error creating lobby: %v", err)
	}
	return body.ID, nil
}

func CheckConnection(appContext *meta.ApplicationContext) *HealthCheckResponseDTO {
	url := fmt.Sprintf("%s/health", appContext.InternalMultiplayerServerAddress)
	resp, err := makeGetRequest[HealthCheckResponseDTO](url)

	if err != nil {
		return &HealthCheckResponseDTO{
			Status:     false,
			LobbyCount: 0,
			Message:    fmt.Sprintf("Error checking connection: %s", err.Error()),
		}
	}
	return resp
}

func GetLobbyState(lobbyID uint32, appContext *meta.ApplicationContext) (*LobbyStateResponseDTO, error) {
	url := fmt.Sprintf("%s/lobby/%d", appContext.InternalMultiplayerServerAddress, lobbyID)
	resp, err := makeGetRequest[LobbyStateResponseDTO](url)
	if err != nil {
		return nil, fmt.Errorf("error getting lobby state: %v", err)
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
