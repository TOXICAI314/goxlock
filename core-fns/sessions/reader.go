package sessions

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
	"time"
)

// - ReadSession
// ReadSession : Will read the given session by its id and returns the desired struct
func ReadSession(id string) (*Session, error) {

	// - Pre Safety 
	if !filepath.IsAbs(Sessionfolder) {
		return nil, &config.FunctionCancelError{
			Cause:   `CWD session path`,
			Message: fmt.Sprintf(`Cannot get the AppdataRoaming path -> Got a cwd path %s`, Sessionfolder),
			ElapsedTime: time.Now(),
			Provider: `session.ReadSession`,
		}
	}
	if id == `` {
		return nil,&config.FunctionCancelError{
			Cause: `Empty id string`,
			Message: `Given an empty id to work by`,
			ElapsedTime: time.Now(),
			Provider: `session.Delete`,
		}
	}
	err := os.MkdirAll(Sessionfolder, 0700)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Making the sessions directory failed : %s`,Sessionfolder),
			ElapsedTime: time.Now(),
			Provider: `session.Delete`,
		}
	}

	path := filepath.Join(Sessionfolder, id+config.JsonExt)
	filedata, err := os.ReadFile(path)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot read the file , Check its there or not and have desried perms or not : %s`,path),
			ElapsedTime: time.Now(),
			Provider: `session.Delete`,
		}
	}
	var s *Session = &Session{}
	err = json.Unmarshal(filedata, s)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot make the session out of the data provided in the file %s`, path),
			ElapsedTime: time.Now(),
			Provider: `session.Delete`,
		}
	}

	return s, nil
}
