package verify

import (
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/header"
	"os"
	"path/filepath"
	"time"
)

// Will verify the given password and will provide error if there are any
func VerifyUnlock(cfg *config.Config) error {
	// Pre Safety
	if cfg == nil {
		return &config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			ElapsedTime: time.Now(),
			Provider: `unlocker.VerifyUnlock`,
		}
	} 
	if ext := filepath.Ext(cfg.FolderName); ext != config.LockExt {
		return &config.FunctionCancelError{
			Cause:   fmt.Sprintf(`Cannot Verify a non '%s' folder`, config.LockExt),
			Message: fmt.Sprintf(`The given extension is %s ; needed -> %s`, ext, config.LockExt),
			ElapsedTime: time.Now(),
			Provider: `unlocker.VerifyUnlock`,
		}
	}

	data, err := os.ReadFile(cfg.FolderName)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot read from the given file - %s`, cfg.FolderName),
			ElapsedTime: time.Now(),
			Provider: `unlocker.VerifyUnlock`,
		}
	}
	_, err = header.GetUnlockedData(cfg, data)
	if err != nil {
		return err
	}

	return nil
}