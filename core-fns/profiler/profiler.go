package profiler

import (
	"goxlock/config"
	"path/filepath"
)

// - Var
var (
	// ProfileName -> The name for the profiler by which its subsystem file structure will go
	ProfilerName string = `profile`
	// ProfileAppData -> The App data folder that is made for storing the profiles
	ProfilerAppDataPath string = filepath.Join(config.GoxLockAppDataFolder,ProfilerName)
	// ProfilePattern -> The pattern by which all files will be made
	ProfilePattern string = `%s-Profile-%s` + config.JsonExt
)

// - Profiler
// Strcut based on the data of user preference and control flow
// It can be switched from user to user to get the certainity of what the user wanted before
// No sesitive data will stored in the struct , Only the data that is based on preference is stored
type Profiler struct {
	// Name : The main Name that the profile will be called
	Name 			string					`json:"name"`
	// Output : Taken from the parent `config.Config.OutputName`
	OutputName 		string					`json:"outputname"`
	// Instructions : Taken from the parent `config.Config.InstructionsData`
	Instruction 	config.Instructions		`json:"instructions"`
}