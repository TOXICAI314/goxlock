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
	if err := pf.Validate();err != nil {
		return err
	}

	x_path := filepath.Join(ProfilerConfigDir, fmt.Sprintf(ProfilePattern, config.Name, pf.Name))
	data, err := os.ReadFile(x_path)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Cannot fetch data from the profile because of the Given Reason : %s`, x_path),
			ElapsedTime: time.Now(),
			Provider:    `profiler.Profiler.Fetch`,
		}
	}

	err = json.Unmarshal(data, pf)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Cannot dump the file data into the struct because of the Given Reason : %+v`, data),
			ElapsedTime: time.Now(),
			Provider:    `profiler.Profiler.Fetch`,
		}
	}

	return nil
}
