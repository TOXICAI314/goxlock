package relocker

import (
	"goxlock/config"
	"goxlock/core-fns/dpapi"
	"goxlock/core-fns/locker"
	"goxlock/core-fns/scheduler"
	"goxlock/core-fns/sessions"
	"time"
)

// A scheduled action that will relock the given folder as user commands
// It dont need config upport as it can read it own data forming `Session`
func Relocker(sessionId string) error {
	if sessionId == `` {
		return &config.FunctionCancelError{
			Cause: `Empty id string`,
			Message: `Given an empty id to work by`,
			ElapsedTime: time.Now(),
			Provider: `relocker.Relock`,
		}
	}

	session,err := sessions.ReadSession(sessionId)
	if err != nil {
		return err
	}

	password,err := dpapi.Unprotect(session.Password)
	if err != nil {
		return err
	}

	stringPassword := string(password)

	defer sessions.Delete(sessionId)
	defer scheduler.DeleteSchedule(sessionId)
	// Config construct 
	// This will reconstruct config to use `locker/locker.go` 
	// [Maybe my first time writting like this in this codebase but i aint writting whole another logic for relock]
	var cfg *config.Config = &config.Config{
		FolderName: session.Folder,
		UserAction: config.LockFolder,
		
		InstructData: session.InstructionData,
		Password: stringPassword,
		OutputName: session.OutputName,
	}

	return locker.Locker(cfg)
}