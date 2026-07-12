package profiler

import (
	"goxlock/config"
)

// Will update the profiler with the given data and then call `pf.Create` to store all the data
func (pf *Profiler) Update(outname string, instruction *config.Instructions) error {
	// Pre Safety
	if err := pf.Validate();err != nil {
		return err
	}

	err := pf.Fetch()
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
	pf.Instruction.Timeout = instruction.Timeout
	pf.Instruction.UnSafe = instruction.UnSafe
	return pf.Create()
}
