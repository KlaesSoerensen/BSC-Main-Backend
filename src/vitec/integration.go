package vitec

import (
	"fmt"
	"log"
	"otte_main_backend/src/config"
)

type SessionInitiationDTO struct {
	UserIdentifier      string `json:"userIdentifier"`
	CurrentSessionToken string `json:"currentSessionToken"`
}

type CrossVerificationType string

const (
	CrossVerificationNever  CrossVerificationType = "never"
	CrossVerificationAlways CrossVerificationType = "always"
)

var UNVERIFIABLE_USER error = fmt.Errorf("user could not be verified")

func neverVerifyUser(integration *VitecIntegration, initiationDTO *SessionInitiationDTO) error {
	log.Println("[MV INT] User verification currently mocked")
	return nil
}

func alwaysVerifyUser(integration *VitecIntegration, initiationDTO *SessionInitiationDTO) error {
	return fmt.Errorf("user could not be verified: Integration Not Implemented")
}

func CreateNewVitecIntegration() (*VitecIntegration, error) {
	crossVerificationType := config.GetOr("VITEC_CROSS_VERIFICATION", "always")
	log.Println("[MV INT] Establishing Vitec Cross Verification. Type: ", crossVerificationType)

	var integration = &VitecIntegration{} //Struct stepwise initialized in this function

	switch CrossVerificationType(crossVerificationType) {
	case CrossVerificationNever:
		integration.VerifyUser = func(initiationDTO *SessionInitiationDTO) error {
			return neverVerifyUser(integration, initiationDTO)
		}
	case CrossVerificationAlways:
		ip, ipErr := config.LoudGet("VITEC_MV_AUTH_IP")
		if ipErr != nil {
			return nil, ipErr
		}
		port, portErr := config.GetInt("VITEC_MV_AUTH_PORT")
		if portErr != nil {
			return nil, portErr
		}
		if err := authEndpointTest(ip, port); err != nil {
			return nil, err
		}
		integration.ip = ip
		integration.port = port
		integration.VerifyUser = func(initiationDTO *SessionInitiationDTO) error {
			return alwaysVerifyUser(integration, initiationDTO)
		}
	default:
		return nil, fmt.Errorf("invalid cross verification type: %s", crossVerificationType)
	}

	return integration, nil
}

func authEndpointTest(ip string, port int) error {
	//var fullurl = fmt.Sprintf("https://%s:%d/cross-verify", ip, port)
	return fmt.Errorf("authEndpointTest not implemented")
}

type VitecIntegration struct {
	ip         string
	port       int
	VerifyUser func(initiationDTO *SessionInitiationDTO) error
}
