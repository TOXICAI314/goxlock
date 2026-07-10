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

// - Create
// Is a function that Creates the Profile from Profiler
// This Profile struct is then dumped into the %APPDATA%/Roaming/name of the application
func (pf *Profiler) Create() error {
	// - Pre Safety
	if !filepath.IsAbs(ProfilerAppDataPath) {
		return &config.FunctionCancelError{
			Cause: `Cwd folder detected`,
			Message: `Absoulte folder needed for storing the Profile`,
			ElapsedTime: time.Now(),
			Provider: `profiler.Profile.Create`,
		}
	}
	err := os.MkdirAll(ProfilerAppDataPath, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot make a folder to store the profile : %s`,ProfilerAppDataPath),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profile.Create`,
		}
	}
	if pf.Name == `` {
		return &config.FunctionCancelError{
			Cause:   `Invalid profile name`,
			Message: `Provide a Valid name to be taken for the profile`,
			ElapsedTime: time.Now(),
			Provider: `profiler.Profile.Create`,
		}
	}

	profilerfilename := filepath.Join(ProfilerAppDataPath, fmt.Sprintf(ProfilePattern, config.Name, pf.Name))
	data, err := json.MarshalIndent(pf, ``, ` `)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot create the json format of the profile provided : %+v`,*pf),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profile.Create`,
		}
	}
	mutex,exists,err := mutex.NewMutex(profilerfilename)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`An internal error has occured while Creating the Mutex for the given folder : %s`,profilerfilename),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profile.Create`,
		}
	}
	defer mutex.CloseMutex()
	if exists {
		return &config.FunctionCancelError{
			Cause:   `Mutex already exist`,
			Message: fmt.Sprintf(`The mutex of the given folder %s is already there`, profilerfilename),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profile.Create`,
		}
	}

	err = os.WriteFile(profilerfilename, data, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot write the file to store the data of the profile : %s`,profilerfilename),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profile.Create`,
		}
	}

	return nil
}
