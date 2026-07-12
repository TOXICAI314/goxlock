package profiler

import (
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/mutex"
	"os"
	"path/filepath"
	"time"
)

// - DeleteProfiler
// Will Delete the profile if user want it so
func (pf *Profiler) Delete() error {
	// - Pre Safety
	if err := pf.Validate();err != nil {
		return err
	}

	profileFile := filepath.Join(ProfilerConfigDir, fmt.Sprintf(ProfilePattern, config.Name, pf.Name))
	mut, alrexist, err := mutex.NewMutex(profileFile)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`An internal error has occured while Creating the Mutex for the given folder : %s`,profileFile),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Delete`,
		}
	}
	defer mut.CloseMutex()
	if alrexist {
		return &config.FunctionCancelError{
			Cause:   `Mutex already exist`,
			Message: fmt.Sprintf(`The mutex of the given folder %s is already there`, profileFile),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Delete`,
		}
	}

	err = os.Remove(profileFile)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot delete the given profile file : %s`,profileFile),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Delete`,
		}
	}
	return nil
}
