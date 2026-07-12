package config

import (
	"fmt"
	"goxlock/utils"
	"os"
	"path/filepath"
	"time"
)

// Version : The way that the data will be read in the future
const (
	VersionMajor = 1
	VersionMinor = 0
	Patch        = 1
	Release      = beta
)

var Version string = fmt.Sprintf(`%d.%d.%d`, VersionMajor, VersionMinor, Patch)
var VersionRelease string = fmt.Sprintf(`%s-%s`, Version, Release)

// Release Verion and Path allocators
const (
	// Alpha = Needs testing and feture ammend (for testers)
	alpha = `Alpha`
	// Beta  = Done featuring but need testing
	beta = `Beta`
	// GA 	 = Available fopr 100% potential with no bugs
	stable = `Stable`
)

// Name : The General name by which the app is running
const Name string = `goxlock`

// Banner : The cli ASCII art that will be shown to the user
const Banner string = `
	  ____  _____  _  _  _      ____   ____ _  __
	 / ___|/ _ \ \/ / | | |    / __ \ / ___| |/ /
	| |  _| | | \  /  | | |   | |  | | |   | ' / 
   	| |_| | |_| /  \  | | |___| |__| | |___| . \ 
	 \____|\___/_/\_\ |_|_____|\____/ \____|_|\_\
`

// - Const Data
const (
	// ZipExt -> The extension by which a zip ext will be recognised
	ZipExt string = `.zip`
	// JsonExt -> The extension repr for the json
	JsonExt string = `.json`
	// LockExt -> The extension for the `goxlock` lock file
	LockExt string = `.g-lock`
)

// - User Action
// - These are the actions that are allowed to the user
// - Values are stored in their corresponding `const`
const (
	LockFolder = iota + 1
	UnlockFolder
	RelockFolder
	ChangePassword
	VerifyPassword
	HeaderCheck
)

// - Imp Var 
var (
	// ConfigDir -> Gives the storage location where the application will store its data
	ConfigDir			string
	// GoxLockConfigDir -> A sub folder for the config dir where goxlock can store its data
	GoxLockConfigDir string
)

// init() 
// will run first when this package is needed and will secure the importants details to the variables
func init() {
	var err error
	// Info : Intial Config Directory to the AppData folder
	ConfigDir, err = os.UserConfigDir()
	if err == nil{
		GoxLockConfigDir = filepath.Join(ConfigDir,Name)	
	}
}

// - Config
// Config : Stores the basic user request to be forwarded afterwards
type Config struct {
	SessionID      string
	StartedAt      time.Time
	FolderName     string
	OutputName     string
	Password       string
	ChangePassword ChangePasswordData
	UserAction     int
	InstructData   Instructions
}

// - Instruction
// Instructions : The parts that user want extra other than `locking` and `unlocking`
type Instructions struct {
	DeleteOriginal bool
	Exclusion      []string
	Stats          bool
	Timeout        time.Duration
	LoggerAllowed  bool
	UnSafe         bool
}

// - ChangePasswordData
type ChangePasswordData struct {
	NewPassword string
}

// - Header
// Header : The base header data for the extension `glock`
type Header struct {
	Magic   [7]byte
	Version [5]byte
	Salt    [16]byte
	Nonce   [12]byte
}

// - Structure
// - Overriding the undesirable config data with usable one
func (cfg *Config) Structure() error {
	if cfg.OutputName == `` {
		cfg.OutputName = cfg.FolderName
	}
	if cfg.SessionID == `` {
		cfg.SessionID = utils.CreateSessionID()
	}
	if cfg.StartedAt.IsZero() {
		cfg.StartedAt = time.Now()
	}
	// Vulnerability Checker for the filepath (no escaping quotes)
	cfg.FolderName = filepath.Clean(cfg.FolderName)
	cfg.OutputName = filepath.Clean(cfg.OutputName)
	switch {
	case !filepath.IsAbs(cfg.FolderName):
		if !filepath.IsLocal(cfg.FolderName) {
			return &FunctionCancelError{
				Cause: `The Path have esaping quotes which are non local`,
				Message: fmt.Sprintf(`The give path %s have quotes that are not allowed by the security`,cfg.FolderName),
				ElapsedTime: time.Now(),
				Provider: `config.Config.Structure`,
			}
		}
	case !filepath.IsAbs(cfg.OutputName):
		if !filepath.IsLocal(cfg.OutputName) {
			return &FunctionCancelError{
				Cause: `The Path have esaping quotes which are non local`,
				Message: fmt.Sprintf(`The give path %s have quotes that are not allowed by the security`,cfg.OutputName),
				ElapsedTime: time.Now(),
				Provider: `config.Config.Structure`,
			}
		}
	}

	return nil
}

// - SharedEncryptionData
// SharedEncryptionData : Is the data that is shared between memory so that other functions and data structure can use them
type SharedEncryptionData struct {
	Salt          [16]byte
	Nonce         [12]byte
	EncryptedData []byte
}
