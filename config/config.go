package config

import (
	"goxlock/utils"
	"os"
	"path/filepath"
	"time"
)

// Version : The way that the data will be read in the future
const Version float32 = 1.0

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
	// LockExt -> The extension for the `goxlock` lock fiulke
	LockExt string = `.g-lock`
)

// - User Action
// - These are the actions that are allowed to the user
// - Values are stored in thier corresponding `const`
const (
	LockFolder     int = 0001
	UnlockFolder   int = 0002
	RelockFolder   int = 0003
	ChangePassword int = 0004
	VerifyPassword int = 0005
	HeaderCheck    int = 0006
)

// - Imp Var - //
var (
	// AppDataRoamingFolder -> The user config directory
	AppDataRoamingFolder string
	// GoxlockAppDataFolder -> A sub folder for the config dir where goxlock can store its data
	GoxLockAppDataFolder string
)

// - init
// init() will run first when this package is needed and will secure the importants details to the variables
func init() {
	var err error
	AppDataRoamingFolder, err = os.UserConfigDir()
	if err != nil || AppDataRoamingFolder == "" {
		// Info : Fallback to environment variable if UserConfigDir fails
		AppDataRoamingFolder = os.Getenv("APPDATA")
	}
	GoxLockAppDataFolder = filepath.Join(AppDataRoamingFolder, Name)
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
	Version [4]byte
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
	return nil
}

// - SharedEncryptionData
// SharedEncryptionData : Is the data that is shared between memory so that other functions and data structure can use them
type SharedEncryptionData struct {
	Salt          []byte
	Nonce         []byte
	EncryptedData []byte
}
