package database

import (
	"fmt"

	"gorm.io/gorm"
)

type PlayerDB = *gorm.DB
type ColonyAssetDB = *gorm.DB
type LanguageDB = *gorm.DB

type DBLoggingLoudness string

const (
	DBLoggingVerbose DBLoggingLoudness = "verbose"
	DBLoggingMinimal DBLoggingLoudness = "minimal"
)

type DBDSN struct {
	Host     string
	Port     uint64
	Username string
	Password string
	Database string
	SSLMode  string
}

func (dsn DBDSN) FullString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dsn.Host, dsn.Port, dsn.Username, dsn.Password, dsn.Database, dsn.SSLMode)
}
func (dsn DBDSN) SafeString() string {
	return fmt.Sprintf("host=%s port=%d user=******** password=******** dbname=%s sslmode=%s",
		dsn.Host, dsn.Port, dsn.Database, dsn.SSLMode)
}
