package logger

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/mutex"
	"os"
	"path/filepath"
	"time"
)

// - Vars
var (
	LoggerName        string = `logs`
	LoggerConfigDir string = filepath.Join(config.GoxLockConfigDir, LoggerName)
	Loggerpattern     string = `goxlock-log-%s` + config.JsonExt
)

// - Logger
// Logger is just a printer that allows printing the extras on the desired `*os.File`
type Logger struct {
	Place                string
	CollectiveLoggerData LoggerData
}

// - LoggerData
// Will keep all the fields form the config that is needed
type LoggerData struct {
	Logs 	[]ConfiguredConfigData		`json:"logs"`
}

// - ConfiguredConfigData
// The safe form of the configuration that can be stored
type ConfiguredConfigData struct {
	SessionID        string              `json:"session"`
	UserAction       int                 `json:"Action"`
	FolderName       string              `json:"folder_name"`
	OutputName       string              `json:"output_name"`
	StartedAt        time.Time           `json:"started_at"`
	CompletedAt      time.Time           `json:"completed_at"`
	Instructions     config.Instructions `json:"instructions"`
	ErrorEncountered string               `json:"error"`
}

// - Write
// Prints the data as normal but in the given destination and Permission
func (lg *Logger) Write() error {
	// - Mutex
	mutex, exists, err := mutex.NewMutex(lg.Place)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`An internal error has occured while Creating the Mutex for the given folder: %s`,lg.Place),
			Provider: `logger.Logger.Write`,	
			ElapsedTime:  time.Now(),
		}
	}
	defer mutex.CloseMutex()
	if exists {
		return &config.FunctionCancelError{
			Cause:   `Mutex already exist`,
			Message: fmt.Sprintf(`The mutex of the given folder %s is already there`, lg.Place),
			Provider: `logger.Logger.Write`,	
			ElapsedTime:  time.Now(),
		}
	}

	data, err := json.MarshalIndent(lg.CollectiveLoggerData, ``, ` `)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot marshal the data into json format for the data %+v`,lg.CollectiveLoggerData),
			Provider: `logger.Logger.Write`,	
			ElapsedTime:  time.Now(),
		}
	}
	err = os.WriteFile(lg.Place, data, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot Write into the file specified : %s`, lg.Place),
			Provider: `logger.Logger.Write`,	
			ElapsedTime:  time.Now(),
		}
	}
	return err
}

// - Read
// It reads the data if the allowance is done
func (lg *Logger) Read(buf *[]byte) error {
	data, err := os.ReadFile(lg.Place)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot Write into the file specified : %s`, lg.Place),
			Provider: `logger.Logger.Read`,	
			ElapsedTime:  time.Now(),
		}
	}
	*buf = data
	return nil
}

// - Log
// The main function which writes into the goxlock appdata folder for log
func Log(cfg *config.Config, errEncountered error) (*Logger, error) {
	if cfg == nil {
		return nil,&config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			ElapsedTime: time.Now(),
			Provider: `logger.Log`,
		}
	} 

	// - Pre Safety
	if !filepath.IsAbs(LoggerConfigDir) {
		return nil, &config.FunctionCancelError{
			Cause:   `Cwd Folder Detected`,
			Message: fmt.Sprintf(`The folder path is a local path not absolute path : %s`,LoggerConfigDir),
			ElapsedTime: time.Now(),
			Provider: `logger.Log`,
		}
	}
	errx := os.MkdirAll(LoggerConfigDir, 0700)
	if errx != nil {
		return nil, &config.FunctionFailError{
			Cause:   errx.Error(),
			Message: fmt.Sprintf(`Cannot Create the folder that is needed to store the logs : %s`,LoggerConfigDir),
			ElapsedTime: time.Now(),
			Provider: `logger.Log`,
		}
	}

	currenttime := time.Now()
	formattedtime := currenttime.Format(`02-01-2006`)
	accessingFile := filepath.Join(LoggerConfigDir, fmt.Sprintf(Loggerpattern, formattedtime))

	file,errx := os.OpenFile(accessingFile,os.O_CREATE|os.O_RDWR,0700)
	if errx != nil {
		return nil,&config.FunctionFailError{
			Cause: errx.Error(),
			Message: fmt.Sprintf(`Cannot open the file for the reading of the data : %s`,accessingFile),
			ElapsedTime: time.Now(),
			Provider: `logger.Log`,
		}
	}
	file.Close()
	previousLoggerdata,errx := os.ReadFile(accessingFile)
	var previouslogs LoggerData = LoggerData{}
	// Info : As the file can be made first time -> unmarshal can give error
	// Application will ignore any error while fetching and will append the data no matter what
	_ = json.Unmarshal(previousLoggerdata,&previouslogs)

	configureddata := ConfiguredConfigData{
		SessionID		:       cfg.SessionID,
		FolderName		:       cfg.FolderName,
		OutputName		:       cfg.OutputName,
		Instructions	:    	cfg.InstructData,
		StartedAt		:       cfg.StartedAt,
		CompletedAt		:      	time.Now(),
		UserAction		:       cfg.UserAction,
	}
	if errEncountered != nil {
		configureddata.ErrorEncountered = errEncountered.Error()
	}
	previouslogs.Logs = append(previouslogs.Logs,configureddata)
	return &Logger{
		Place				:				accessingFile,
		CollectiveLoggerData: 				previouslogs ,
	}, nil
}

// - ReadLogFile 
// This Reads the log file and returns the Logger that is needed 
func ReadLogFile(formattedTime string) (*Logger,error) {
	// - Pre Safety
	if !filepath.IsAbs(LoggerConfigDir) {
		return nil, &config.FunctionCancelError{
			Cause:   `Cwd Folder Detected`,
			Message: fmt.Sprintf(`The folder path is a local path not absolute path : %s`,LoggerConfigDir),
			ElapsedTime: time.Now(),
			Provider: `logger.ReadLogFile`,
		}

	}
	err := os.MkdirAll(LoggerConfigDir, 0700)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot Create the folder that is needed to store the logs : %s`,LoggerConfigDir),
			ElapsedTime: time.Now(),
			Provider: `logger.ReadLogFile`,
		}
	}
	// Info : If the time is not given then current date will be used
	if formattedTime == `` {
		return nil,&config.FunctionCancelError{
			Cause: `Invalid formatting time`,
			Message: `The time given by the user is invalid to use -> Empty string`,
			ElapsedTime: time.Now(),
			Provider: `logger.ReadLogFile`,
		}
	}

	accessingFile := filepath.Join(LoggerConfigDir, fmt.Sprintf(Loggerpattern, formattedTime))
	data,err := os.ReadFile(accessingFile)
	if err != nil {
		return nil,&config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot read from the log file and fetch the logged data : %s`,accessingFile),
			ElapsedTime: time.Now(),
			Provider: `logger.ReadLogFile`,
		}
	} 

	var previouslogs LoggerData = LoggerData{}
	err = json.Unmarshal(data,&previouslogs)
	// Info : Here the error will be counted as there is no restriction of the file to be made or to be fresh
	if err != nil {
		return nil,&config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`The data cannot be unmarshlled into a formatted type : %+v`,data),
			ElapsedTime: time.Now(),
			Provider: `logger.ReadLogFile`,
		}
	}
	return &Logger{
		Place: accessingFile,
		CollectiveLoggerData: previouslogs,
	},nil
}