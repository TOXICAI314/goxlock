package profiler

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/mutex"
	"os"
	"path/filepath"
	"time"
)

// Create
// Is a function that Creates the Profile from Profiler
// This Profile struct is then dumped into the %APPDATA%/Roaming/name of the application
func (pf *Profiler) Create() error {
	// Pre Safety
	if err := pf.Validate();err != nil {
		return err
	}

	profilerfilename := filepath.Join(ProfilerConfigDir, fmt.Sprintf(ProfilePattern, config.Name, pf.Name))
	data, err := json.MarshalIndent(pf, ``, ` `)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Cannot create the json format of the profile provided : %+v`, *pf),
			ElapsedTime: time.Now(),
			Provider:    `profiler.Profile.Create`,
		}
	}
	mutex, exists, err := mutex.NewMutex(profilerfilename)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`An internal error has occured while Creating the Mutex for the given folder : %s`, profilerfilename),
			ElapsedTime: time.Now(),
			Provider:    `profiler.Profile.Create`,
		}
	}
	defer mutex.CloseMutex()
	if exists {
		return &config.FunctionCancelError{
			Cause:       `Mutex already exist`,
			Message:     fmt.Sprintf(`The mutex of the given folder %s is already there`, profilerfilename),
			ElapsedTime: time.Now(),
			Provider:    `profiler.Profile.Create`,
		}
	}

	err = os.WriteFile(profilerfilename, data, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Cannot write the file to store the data of the profile : %s`, profilerfilename),
			ElapsedTime: time.Now(),
			Provider:    `profiler.Profile.Create`,
		}
	}

	return nil
}
