package profiler

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
	"time"
)

// - ProfileFetcher
// Fetches the profile data and fill it in thee given Instructor Profiler struct
func (pf *Profiler) Fetch() error {
	// - Pre Safety
	if !filepath.IsAbs(ProfilerAppDataPath) {
		return &config.FunctionCancelError{
			Cause: `Cwd folder detected`,
			Message: fmt.Sprintf(`Absoulte folder needed for storing the Profile : %s`,ProfilerAppDataPath),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Fetch`,
		}
	}
	err := os.MkdirAll(ProfilerAppDataPath, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot make a folder to store the profile : %s`,ProfilerAppDataPath),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Fetch`,
		}
	}
	if pf.Name == `` {
		return &config.FunctionCancelError{
			Cause:   `Invalid profile name`,
			Message: `Provide a Valid name to be taken for the profile`,
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Fetch`,
		}
	}

	x_path := filepath.Join(ProfilerAppDataPath, fmt.Sprintf(ProfilePattern, config.Name, pf.Name))
	data,err := os.ReadFile(x_path)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message : fmt.Sprintf(`Cannot fetch data from the profile because of the Given Reason : %s`,x_path),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Fetch`,
		}
	}
	
	err = json.Unmarshal(data,pf)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot dump the file data into the struct because of the Given Reason : %+v`,data),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Fetch`,
		}
	}

	return nil
}