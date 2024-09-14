package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/util"
	"time"
)

type SessionToken string
type Session struct {
	ID     uint32       `json:"id" gorm:"primaryKey"` // ID
	Player uint32       `json:"player" gorm:"foreignKey:player;references:ID"`
	Token  SessionToken `json:"token"`
	//MS
	ValidDuration uint32    `json:"validDuration" gorm:"column:validDuration"` //Column defaults to 1h, or 3600000ms
	CreatedAt     time.Time `json:"createdAt" gorm:"column:createdAt"`         //Column defaults to NOW()
	LastCheckIn   time.Time `json:"lastCheckIn" gorm:"column:lastCheckIn"`     //Column defaults to NOW()
}

func (s *Session) TableName() string {
	return "Session"
}

// MS
const DEFAULT_VALID_DURATION = 3600000 //1 hour

func CreateSessionForPlayer(playerID uint32, appContext *meta.ApplicationContext, idOfPreviousSession int) (*Session, error) {
	var token, generationErr = generateBase32String(64)
	if generationErr != nil {
		return nil, fmt.Errorf("unable to generate session token")
	}

	var session = Session{
		ID:            uint32(util.Ternary(idOfPreviousSession == -1, 0, idOfPreviousSession)),
		Token:         SessionToken(token),
		Player:        playerID,
		ValidDuration: DEFAULT_VALID_DURATION,
		CreatedAt:     time.Now(),
		LastCheckIn:   time.Now(),
	}

	if createSessionError := appContext.PlayerDB.Save(&session).Error; createSessionError != nil {
		return nil, fmt.Errorf("unable to save session")
	}
	return &session, nil
}

func generateBase32String(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes), nil
}
