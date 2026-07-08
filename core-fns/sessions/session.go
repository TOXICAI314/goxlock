package sessions

import (
	"goxlock/config"
	"path/filepath"
)

// - Const -
const (
	// SessionName ->  The name for the inside folder of the sessions holder (just the name)
	SessionsName = `sessions`
)

// - Imp vars 
var (
	// Sessionfolder -> The app data folder for the session
	Sessionfolder		 string
)

// - Session 
// Is the struct that will record all the user important details to use them further when needed
type Session struct {
	Id 				string  				`json:"id"`
	Folder			string					`json:"folder"`
	Password 		[]byte					`json:"password"`
	OutputName		string					`json:"outputname"`
	InstructionData config.Instructions 	`json:"instructions"`
}

// - init 
// will run first when this package is needed and will secure the importants details to the variables
func init() {
	Sessionfolder = filepath.Join(config.GoxLockAppDataFolder, SessionsName)
}