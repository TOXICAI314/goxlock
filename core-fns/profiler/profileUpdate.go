package profiler

import (
	"goxlock/config"
	"os"
	"path/filepath"
)

// - Update
// Will update the profiler with the given data and then call `pf.Create` to store all the data
func (pf *Profiler) Update(outname string, instruction *config.Instructions) error {
	// - Pre Safety
	if instruction == nil {
		return &config.UserSafetyError{
			Cause:   `Nil Pointer dereference`,
			Message: `The given instructions is a pointer to nil`,
		}
	}
	if !filepath.IsAbs(ProfilerAppDataPath) {
		return &config.UserSafetyError{
			Cause:   `Cwd folder detected`,
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
