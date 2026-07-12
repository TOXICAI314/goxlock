package scheduler

import (
	"fmt"
	"goxlock/config"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CreateShedule : Will create the shedules task and will run with the help of `sctasks` native to windows
func CreateShedule(sessionID string,instructions config.Instructions) error {
	// - Pre Safety
	if sessionID == `` {
		return &config.FunctionCancelError{
			Cause: `Empty id string`,
			Message: `Given an empty id to work by`,
			ElapsedTime: time.Now(),
			Provider: `profiler.CreateSchedule`,
		}
	}

	// Info : Os.executable is gives the current running executable path of the executable
	var exe string
	var err  error
	exe, err = os.Executable()
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: `Cannot get the executable path of the current exe to schedule the task`,
			ElapsedTime: time.Now(),
			Provider: `profiler.CreateSchedule`,
		}
	}

	taskName := config.Name + "-" + sessionID

	currentTime := time.Now()
	runTime := currentTime.Add(instructions.Timeout)

	if runTime.Before(currentTime.Add(1 * time.Minute)) {
		fmt.Println(`A minimum of 1 min is required by the task scheduler of windows -> Making your Timeduration of 1 min`)
		runTime = currentTime.Add(1 * time.Minute)
	}
	// Info : `schtasks` dont support second precision so rest will be truncated
	date := runTime.Format(`02/01/2006`)
	clock := runTime.Format(`15:04`)

	exclusiontring := strings.Join(instructions.Exclusion, ",")
	var args string = fmt.Sprintf(`--re-lock --session "%s" --exclude "%s"`, sessionID,exclusiontring)

	cmd := exec.Command(
		"schtasks",

		"/Create",

		"/TN", taskName,

		"/TR", fmt.Sprintf(`"%s" %s`, exe, args),

		"/SC", "ONCE",

		"/SD", date,

		"/ST", clock,

		"/F",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Error while scheduling the task to the sctask : %s`, string(out)),
			ElapsedTime: time.Now(),
			Provider: `profiler.CreateSchedule`,
		}
	}

	return nil
}
