package logger

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
	"time"
)


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