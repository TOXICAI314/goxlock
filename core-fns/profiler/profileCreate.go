package profiler

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/mutex"
	"os"
	"path/filepath"
)

// - Create
// Is a function that Creates the Profile from Profiler
// This Profile struct is then dumped into the %APPDATA%/Roaming/name of the application
func (pf *Profiler) Create() error {
	// - Pre Safety
	if !filepath.IsAbs(ProfilerAppDataPath) {
		return &config.UserSafetyError{
			Cause: `Cwd folder detected`,
			Message: `Absoulte folder needed for storing the Profile`,
		}
	}
	err := os.MkdirAll(ProfilerAppDataPath, 0700)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `Cannot make a folder to store the profile`,
		}
	}
	if pf.Name == `` {
		return &config.UserSafetyError{
			Cause:   `Invalid profile name`,
			Message: `Provide a Valid name to be taken for the profile`,
		}
	}

	profilerfilename := filepath.Join(ProfilerAppDataPath, fmt.Sprintf(ProfilePattern, config.Name, pf.Name))
	data, err := json.MarshalIndent(pf, ``, ` `)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `Cannot create the json format of the profile provided`,
		}
	}
	mutex,exists,err := mutex.NewMutex(profilerfilename)
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
			Message: fmt.Sprintf(`The mutex of the given folder %s is already there`, profilerfilename),
		}
	}

	err = os.WriteFile(profilerfilename, data, 0700)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `Cannot write the file to store the data of the profile`,
		}
	}

	return nil
}
