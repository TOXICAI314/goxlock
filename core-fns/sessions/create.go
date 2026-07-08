package sessions

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/mutex"
	"os"
	"path/filepath"
)

// - CreateSession
// Makes the session to its Program files and saves the data for further use cases
func (s *Session) CreateSession() error {

	// - Pre Safety Run 
	if !filepath.IsAbs(Sessionfolder) {
		return &config.UserSafetyError{
			Cause:   `CWD session path`,
			Message: fmt.Sprintf(`Cannot get the AppdataRoaming path -> Got a cwd path %s`, Sessionfolder),
		}
	}
	err := os.MkdirAll(Sessionfolder, 0700)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `Making the sessions directory failed`,
		}
	}

	// - Making session file 
	data, err := json.MarshalIndent(s, ``, ` `)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `Error in marsheling the Session , confirm its structure`,
		}
	}

	filename := filepath.Join(Sessionfolder, s.Id+config.JsonExt)
	mutex, exists, err := mutex.NewMutex(filename)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `An internal error has occured while Creating the Mutex for the given folder`,
		}
	}
	defer mutex.CloseMutex()
	if exists {
		return &config.UserSafetyError{
			Cause:   `Mutex already exist`,
			Message: fmt.Sprintf(`The mutex of the given folder %s is already there`, filename),
		}
	}
	// - Data Dump 
	err = os.WriteFile(filename, data, 0700)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `Errors in writting into the file`,
		}
	}

	return nil
}
