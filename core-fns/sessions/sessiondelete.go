package sessions

import (
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/mutex"
	"os"
	"path/filepath"
	"time"
)

// Delete : deletes the given session with the given id
func Delete(id string) (err error) {

	// Pre safety
	if err := SessionValidate(id);err != nil {
		return err
	}

	path := filepath.Join(SessionConfigDir, id+config.JsonExt)
	mut, alrexist, err := mutex.NewMutex(path)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`An internal error has occured while Creating the Mutex for the given folder : %s`,path),
			ElapsedTime: time.Now(),
			Provider: `unlocker.Unlocker`,
		}
	}
	defer func() {
		closeErr := mut.CloseMutex()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if alrexist {
		return &config.FunctionCancelError{
			Cause:   `Mutex already exist`,
			Message: fmt.Sprintf(`The mutex of the given folder %s is already there`, path),
			ElapsedTime: time.Now(),
			Provider: `unlocker.Unlocker`,
		}
	}

	err = os.Remove(path)
	if err != nil {
		return &config.FunctionFailError{
			Cause:       err.Error(),
			Message:     fmt.Sprintf(`Cannot delete the given file with the id %s`, id),
			ElapsedTime: time.Now(),
			Provider:    `session.Delete`,
		}
	}
	return nil
}
