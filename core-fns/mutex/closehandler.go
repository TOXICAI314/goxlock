package mutex

import (
	"goxlock/config"
	"time"

	"golang.org/x/sys/windows"
)

// Will close the given folder mutex -> Normally ending
// Even if crashed app , os will free automatically
func (fMutex *FolderMutex) CloseMutex() error {
	if fMutex.handle == 0 {
		// Info : If the handle is already 0 (closed) -> return nothing
		return nil
	}
	// else juts close it and assign it as 0 (freeing the handle)
	// In `C` thats `NULL` -> `nil`
	err := windows.CloseHandle(fMutex.handle)
	if err != nil {
		return  &config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot close the handle of the mutex`,
			ElapsedTime: time.Now(),
			Provider: `mutex.CloseMutex`,
		}
	}
	fMutex.handle = 0
	return nil
}