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
	if err := SessionValidate(id);err != nil {
		return nil,err
	}

	err := os.MkdirAll(SessionConfigDir, 0700)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Making the sessions directory failed : %s`, SessionConfigDir),
			ElapsedTime: time.Now(),
			Provider:    `session.ReadSession`,
		}
	}

	path := filepath.Join(SessionConfigDir, id+config.JsonExt)
	filedata, err := os.ReadFile(path)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Cannot read the file , Check its there or not and have desried perms or not : %s`, path),
			ElapsedTime: time.Now(),
			Provider:    `session.ReadSession`,
		}
	}
	var s *Session = &Session{}
	err = json.Unmarshal(filedata, s)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Cannot make the session out of the data provided in the file %s`, path),
			ElapsedTime: time.Now(),
			Provider:    `session.Delete`,
		}
	}

	return s, nil
}
