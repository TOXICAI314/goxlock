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

// Prints the data as normal but in the given destination and Permission
func (lg *Logger) Write() (err error) {
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
	defer func() {
		closeErr := mutex.CloseMutex()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()
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

	// Pre Safety
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
		VersionInfo		: 		cfg.Version,
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