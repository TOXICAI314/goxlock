package sessions

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
)

// - ReadSession 
// ReadSession : Will read the given session by its id and returns the desired struct
func ReadSession(id string) (*Session, error) {

	// - Pre Safety 
	if !filepath.IsAbs(Sessionfolder) {
		return nil, &config.UserSafetyError{
			Cause:   `CWD session path`,
			Message: fmt.Sprintf(`Cannot get the AppdataRoaming path -> Got a cwd path %s`, Sessionfolder),
		}
	}
	if id == `` {
		return nil,&config.UserSafetyError{
			Cause: `Empty id string`,
			Message: `Given an empty id to work by`,
		}
	}
	err := os.MkdirAll(Sessionfolder, 0700)
	if err != nil {
		return nil, &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `Making the sessions directory failed`,
		}
	}

	path := filepath.Join(Sessionfolder, id+config.JsonExt)
	filedata, err := os.ReadFile(path)
	if err != nil {
		return nil, &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `Cannot read the file , Check its there or not and have desried perms or not`,
		}
	}
	var s *Session = &Session{}
	err = json.Unmarshal(filedata, s)
	if err != nil {
		return nil, &config.UserSafetyError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot make the session out of the data provided in the file %s`, path),
		}
	}

	return s, nil
}
