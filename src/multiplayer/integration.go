package multiplayer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"otte_main_backend/src/meta"
)

type CreateLobbyResponseDTO struct {
	ID uint32 `json:"id"`
}

// Returns lobbyID, error
func CreateLobby(ownerID uint32, appContext *meta.ApplicationContext) (uint32, error) {
	url := fmt.Sprintf("%s/create-lobby?ownerID=%d&encoding=binary", appContext.MultiplayerServerAddress, ownerID)
	body, err := makeEmptyPostRequest(url)
	if err != nil {
		return 0, fmt.Errorf("error creating lobby: %v", err)
	}
	return body.ID, nil
}

func makeEmptyPostRequest(url string) (*CreateLobbyResponseDTO, error) {
	// Create a new POST request with no body
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Send the request
	client := &http.Client{}
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
