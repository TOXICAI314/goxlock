package scheduler

import (
	"fmt"
	"goxlock/config"
	"os/exec"
	"time"
)

// - DeletSchedule()
// Will delete shedule from the task schdules which was scheduled before
func DeleteSchedule(sessionID string) error {

    // - Pre Safety
    if sessionID == `` {
		return &config.FunctionCancelError{
			Cause: `Empty id string`,
			Message: `Given an empty id to work by`,
            ElapsedTime: time.Now(),
			Provider: `profiler.DeleteSchedule`,
		}
	}

    task := "GoxLock-" + sessionID

    cmd := exec.Command(

        "schtasks",

        "/Delete",

        "/TN", task,

        "/F",
    )

    out, err := cmd.CombinedOutput()

    if err != nil {
        return fmt.Errorf(
            "%v\n%s",
            err,
            string(out),
        )
    }

    return nil
}