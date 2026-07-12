package profiler

import (
	"fmt"
	"goxlock/config"
	"goxlock/utils"
	"os"
	"path/filepath"
	"time"
)

// - Var
var (
	// ProfileName -> The name for the profiler by which its subsystem file structure will go
	ProfilerName string = `profile`
	// ProfileAppData -> The App data folder that is made for storing the profiles
	ProfilerConfigDir string = filepath.Join(config.GoxLockConfigDir, ProfilerName)
	// ProfilePattern -> The pattern by which all files will be made
	ProfilePattern string = `%s-Profile-%s` + config.JsonExt
)

// - Profiler
// Strcut based on the data of user preference and control flow
// It can be switched from user to user to get the certainity of what the user wanted before
// No sesitive data will stored in the struct , Only the data that is based on preference is stored
type Profiler struct {
	// Name : The main Name that the profile will be called
	Name string `json:"name"`
	// Output : Taken from the parent `config.Config.OutputName`
	OutputName string `json:"outputname"`
	// Instructions : Taken from the parent `config.Config.InstructionsData`
	Instruction config.Instructions `json:"instructions"`
}

// Validate
// Validates the profile and returns an error upon that
func (pf *Profiler) Validate() error {
	if !filepath.IsAbs(ProfilerConfigDir) {
		return &config.FunctionFailError{
			Cause:       `Cwd folder detected`,
			Message:     fmt.Sprintf(`Absoulte folder needed for storing the Profile : %s`, ProfilerConfigDir),
			ElapsedTime: time.Now(),
			Provider:    `profiler.Profiler.Validate`,
		}
	}
	err := os.MkdirAll(ProfilerConfigDir, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Cannot make a folder to store the profile : %s`, ProfilerConfigDir),
			ElapsedTime: time.Now(),
			Provider:    `profiler.Profiler.Validate`,
		}
	}
	if pf.Name == `` {
		return &config.FunctionCancelError{
			Cause:       `Invalid profile name`,
			Message:     `Provide a Valid name to be taken for the profile`,
			ElapsedTime: time.Now(),
			Provider:    `profiler.Profiler.Validate`,
		}
	}
	valid := utils.ValidateCharacters(pf.Name)
	if !valid {
		return &config.FunctionCancelError{
			Cause: `Invalid Name`,
			Message: fmt.Sprintf(`The given Name - %s does not fit in the base char set`,pf.Name),
			ElapsedTime: time.Now(),
			Provider: `profiler.Profiler.Validate`,
		}
	}
	return nil
}