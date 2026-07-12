package sessions

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/mutex"
	"os"
	"path/filepath"
	"time"
)

// - CreateSession
// Makes the session to its Program files and saves the data for further use cases
func (s *Session) CreateSession() error {

	// - Pre Safety Run
	if !filepath.IsAbs(SessionConfigDir) {
		return &config.FunctionCancelError{
			Cause:       `CWD session path`,
			Message:     fmt.Sprintf(`Cannot get the AppdataRoaming path -> Got a cwd path %s`, SessionConfigDir),
			ElapsedTime: time.Now(),
			Provider:    `session.Session.CreateSession`,
		}
	}
	err := os.MkdirAll(SessionConfigDir, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Making the sessions directory failed : %s`, SessionConfigDir),
			ElapsedTime: time.Now(),
			Provider:    `session.Session.CreateSession`,
		}
	}

	// - Making session file
	data, err := json.MarshalIndent(s, ``, ` `)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Error in marsheling the Session , confirm its structure : %+v`, *s),
			ElapsedTime: time.Now(),
			Provider:    `session.Session.CreateSession`,
		}
	}

	filename := filepath.Join(SessionConfigDir, s.Id+config.JsonExt)
	mutex, exists, err := mutex.NewMutex(filename)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`An internal error has occured while Creating the Mutex for the given folder : %s`, filename),
			ElapsedTime: time.Now(),
			Provider:    `session.Session.CreateSession`,
		}
	}
	defer mutex.CloseMutex()
	if exists {
		return &config.FunctionCancelError{
			Cause:       `Mutex already exist`,
			Message:     fmt.Sprintf(`The mutex of the given folder %s is already there`, filename),
			ElapsedTime: time.Now(),
			Provider:    `session.Session.CreateSession`,
		}
	}
	// - Data Dump
	err = os.WriteFile(filename, data, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Errors in writting into the file : %s`, filename),
			ElapsedTime: time.Now(),
			Provider:    `session.Session.CreateSession`,
		}
	}

	return nil
}
