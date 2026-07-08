package profiler

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
)

// - ProfileFetcher
// Fetches the profile data and fill it in thee given Instructor Profiler struct
func (pf *Profiler) Fetch() error {
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

	data,err := os.ReadFile(filepath.Join(ProfilerAppDataPath, fmt.Sprintf(ProfilePattern, config.Name, pf.Name)))
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message : `Cannot fetch data from the profile because of the Given Reason`,
		}
	}
	
	err = json.Unmarshal(data,pf)
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot dump the file data into the struct because of the Given Reason`,
		}
	}

	return nil
}