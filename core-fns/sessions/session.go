package sessions

import (
	"fmt"
	"goxlock/config"
	"path/filepath"
	"time"
)

const (
	// SessionName ->  The name for the inside folder of the sessions holder (just the name)
	SessionsName = `sessions`
)

var (
	// SessionConfigDir -> The app data folder for the session
	SessionConfigDir string
)

// Is the struct that will record all the user important details to use them further when needed
type Session struct {
	VersionInfo		string				`json:"version_info"`
	Id              string              `json:"id"`
	Folder          string              `json:"folder"`
	Password        []byte              `json:"password"`
	OutputName      string              `json:"outputname"`
	InstructionData config.Instructions `json:"instructions"`
}

// will run first when this package is needed and will secure the importants details to the variables
func init() {
	SessionConfigDir = filepath.Join(config.GoxLockConfigDir, SessionsName)
}

// Validates the session data and its content
func SessionValidate(id string) error {
	if !filepath.IsAbs(SessionConfigDir) {
		return &config.FunctionCancelError{
			Cause:       `CWD session path`,
			Message:     fmt.Sprintf(`Cannot get the AppdataRoaming path -> Got a cwd path %s`, SessionConfigDir),
			ElapsedTime: time.Now(),
			Provider:    `session.Session.Validate`,
		}
	}
	if id == `` {
		return &config.FunctionCancelError{
			Cause:       `Empty id string`,
			Message:     `Given an empty id to work by`,
			ElapsedTime: time.Now(),
			Provider:    `session.Session.Validate`,
		}
	}
	return nil
}