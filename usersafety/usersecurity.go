package usersafety

import (
	"fmt"
	"goxlock/config"
)

// - SecureFolder
// - SecureFolder : Will make sure none of the badfolder can enter the lock chart
func SecureFolder(cfg *config.Config) error {
	folder := cfg.FolderName
	if cfg.InstructData.DeleteOriginal {
		if bad,_ := BadFolders[folder];bad {
		return &config.UserSafetyError{
			Cause: fmt.Sprintf("A bad folder input of : %s",folder),
			Message: `Cannot use a bad folder as it posses risk of destroying highly sensible folder trees`,
			}
		}
	}
	return nil
}
// - SecureFolder
