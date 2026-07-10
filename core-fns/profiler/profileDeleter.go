package profiler

import (
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
	"time"
)

// - DeleteProfiler
// Will Delete the profile if user want it so
func (pf *Profiler) Delete() error {
	// - Pre Safety
	if !filepath.IsAbs(ProfilerAppDataPath) {
		return &config.FunctionFailError{
			Cause: `Cwd folder detected`,
			Message: fmt.Sprintf(`Absoulte folder needed for storing the Profile : %s`,ProfilerAppDataPath),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Delete`,
		}
	}
	err := os.MkdirAll(ProfilerAppDataPath, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot make a folder to store the profile : %s`,ProfilerAppDataPath),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Delete`,
		}
	}
	if pf.Name == `` {
		return &config.FunctionCancelError{
			Cause:   `Invalid profile name`,
			Message: `Provide a Valid name to be taken for the profile`,
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Delete`,
		}
	}
	
	return os.Remove(filepath.Join(ProfilerAppDataPath, fmt.Sprintf(ProfilePattern, config.Name, pf.Name)))
}