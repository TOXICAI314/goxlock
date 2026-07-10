package profiler

import (
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
	"time"
)

// - Update
// Will update the profiler with the given data and then call `pf.Create` to store all the data
func (pf *Profiler) Update(outname string, instruction *config.Instructions) error {
	// - Pre Safety
	if instruction == nil {
		return &config.FunctionCancelError{
			Cause:   `Nil Pointer dereference`,
			Message: `The given instructions is a pointer to nil`,
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Update`,
		}
	}
	if !filepath.IsAbs(ProfilerAppDataPath) {
		return &config.FunctionCancelError{
			Cause:   `Cwd folder detected`,
			Message: fmt.Sprintf(`Absoulte folder needed for storing the Profile : %s`,ProfilerAppDataPath),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Update`,
		}
	}
	err := os.MkdirAll(ProfilerAppDataPath, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot make a folder to store the profile: %s`,ProfilerAppDataPath),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Update`,
		}
	}
	if pf.Name == `` {
		return &config.FunctionCancelError{
			Cause:   `Invalid profile name`,
			Message: `Provide a Valid name to be taken for the profile`,
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Update`,
		}
	}

	err = pf.Fetch()
	if err != nil {
		return err
	}

	if outname != `` {
		pf.OutputName = outname
	}
	if pf.Instruction.DeleteOriginal != instruction.DeleteOriginal {
		pf.Instruction.DeleteOriginal = instruction.DeleteOriginal
	}
	if len(instruction.Exclusion) != 0 {
		pf.Instruction.Exclusion = instruction.Exclusion
	}
	if pf.Instruction.Timeout != instruction.Timeout {
		pf.Instruction.Timeout = instruction.Timeout
	}
	if pf.Instruction.UnSafe != instruction.UnSafe {
		pf.Instruction.UnSafe = instruction.UnSafe
	}
	return pf.Create()
}
