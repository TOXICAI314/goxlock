package github

import (
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/logger"
	"goxlock/utils"
	"time"
)

// reports the bug to the github page by retruning a url which can be opened later by the parent
func BugReport() (url string, err error) {
	inp := utils.Input(`(Optional) Enter the date for the logs file. [Enter if you dont want]:`)
	retData := []logger.ConfiguredConfigData{}
	if inp != `` {
		log, err := logger.ReadLogFile(inp)
		if err != nil {
			return ``, err
		}
		logLength := len(log.CollectiveLoggerData.Logs)
		switch {
		case logLength > 0 && logLength <= 5:
			retData = log.CollectiveLoggerData.Logs
		case logLength > 5:
			retData = log.CollectiveLoggerData.Logs[logLength-4:]
		default:
			return ``, &config.FunctionCancelError{
				Cause:       `The length of the log is invalid`,
				Message:     fmt.Sprintf(`The length %d is invalid to be reported`, logLength),
				ElapsedTime: time.Now(),
				Provider:    `github.BugReport`,
			}
		}
	}
	return fmt.Sprintf(`%s?template=bug_report.md&body=%+v`, GithubReportURL, retData), nil
}
