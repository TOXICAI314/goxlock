package logger

import (
	"goxlock/config"
	"path/filepath"
	"time"
)

var (
	LoggerName        string = `logs`
	LoggerConfigDir string = filepath.Join(config.GoxLockConfigDir, LoggerName)
	Loggerpattern     string = `goxlock-log-%s` + config.JsonExt
)

// Logger is just a printer that allows printing the extras on the desired `*os.File`
type Logger struct {
	Place                string
	CollectiveLoggerData LoggerData
}

// Will keep all the fields form the config that is needed
type LoggerData struct {
	Logs 	[]ConfiguredConfigData		`json:"logs"`
}

// The safe form of the configuration that can be stored
type ConfiguredConfigData struct {
	VersionInfo		 string				 `json:"version_info"`
	SessionID        string              `json:"session"`
	UserAction       int                 `json:"Action"`
	FolderName       string              `json:"folder_name"`
	OutputName       string              `json:"output_name"`
	StartedAt        time.Time           `json:"started_at"`
	CompletedAt      time.Time           `json:"completed_at"`
	Instructions     config.Instructions `json:"instructions"`
	ErrorEncountered string              `json:"error"`
}