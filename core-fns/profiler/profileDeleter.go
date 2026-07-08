package profiler

import (
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
)

// - DeleteProfiler
// Will Delete the profile if user want it so
func (pf *Profiler) Delete() error {
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
	
	return os.Remove(filepath.Join(ProfilerAppDataPath, fmt.Sprintf(ProfilePattern, config.Name, pf.Name)))
}