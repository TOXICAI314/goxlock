package sessions

import (
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
)

// - Delete 
// Delete : deletes the given session with the given id
func Delete(id string) error {

	// - Pre safety 
	if !filepath.IsAbs(Sessionfolder) {
		return &config.UserSafetyError{
			Cause:   `CWD session path`,
			Message: fmt.Sprintf(`Cannot get the AppdataRoaming path -> Got a cwd path %s`, Sessionfolder),
		}
	}
	if id == `` {
		return &config.UserSafetyError{
			Cause: `Empty id string`,
			Message: `Given an empty id to work by`,
		}
	}

	path := filepath.Join(Sessionfolder, id+config.JsonExt)
	err := os.Remove(path)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot delete the given file with the id %s`, id),
		}
	}
	return nil
}
