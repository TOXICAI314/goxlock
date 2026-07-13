package github

import (
	"fmt"
	"goxlock/config"
	"goxlock/utils"
	"os/exec"
	"time"
)

// reports to the provided github link by the url manipulationn and input system
func Report() (err error) {
	fmt.Print(`Choose the integral value to use for the report`)
	fmt.Println(`
		1. 🐞 Bug Report
		2. ✨ Feature Request
		3. 💬 Support / Question
	`)
	inp := utils.Input(`Your choice : `)
	var targetURL string
	switch inp {
	case `1`:
		targetURL, err = BugReport()
		if err != nil {
			return err
		}
	case `2`:
		targetURL,err = FeatureRequest()
		if err != nil {
			return err
		}
	case `3`:
		targetURL,err = QuestionAndSupport()
		if err != nil {
			return err
		}
	default:
		return &config.FunctionCancelError{
			Cause:       `Invalid input report`,
			Message:     fmt.Sprintf(`Its seems like your given option %s does not match the given answer set. `, inp),
			ElapsedTime: time.Now(),
			Provider:    `gihub.Report`,
		}
	}
	err = exec.Command("cmd", "/c", "start", targetURL).Start()
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Cannot open the url - %s on the browser`, targetURL),
			ElapsedTime: time.Now(),
			Provider:    `gihub.Report`,
		}
	}
	return nil
}
