package sessions

import (
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
	"time"
)

// - Delete
// Delete : deletes the given session with the given id
func Delete(id string) error {

	// - Pre safety 
	if !filepath.IsAbs(Sessionfolder) {
		return &config.FunctionCancelError{
			Cause:   `CWD session path`,
			Message: fmt.Sprintf(`Cannot get the AppdataRoaming path -> Got a cwd path %s`, Sessionfolder),
			ElapsedTime: time.Now(),
			Provider: `session.Delete`,
		}
	}
	if id == `` {
		return &config.FunctionCancelError{
			Cause: `Empty id string`,
			Message: `Given an empty id to work by`,
			ElapsedTime: time.Now(),
			Provider: `session.Delete`,
		}
	}

	path := filepath.Join(Sessionfolder, id+config.JsonExt)
	err := os.Remove(path)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot delete the given file with the id %s`, id),
			ElapsedTime: time.Now(),
			Provider: `session.Delete`,
		}
	}
	return nil
}
