package doctor

import (
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
)

// - CreateLockConfig 
// This creates a lock for the doctor to check -> Everything will be in temp folder
func (dc *Doctor) CreateTestLockerConfig() (error) {
	tempfolder,err := os.MkdirTemp(``,fmt.Sprintf(`%s-test-%s`,config.Name,dc.config.SessionID))
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot make a temp folder for the test`,
		}
	}
	err = os.WriteFile(filepath.Join(tempfolder,`r1.txt`),[]byte(`Wassup man , Having a great day ??`),0700)
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Above error shows why the file cant be made`,
		}
	}
	err = os.WriteFile(filepath.Join(tempfolder,`r2.txt`),[]byte(`Wassup man , Having a great day ??`),0700)
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Above error shows why the file cant be made`,
		}
	}
	err = os.WriteFile(filepath.Join(tempfolder,`r3.txt`),[]byte(`Wassup man , Having a great day ??`),0700)
	if err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Above error shows why the file cant be made`,
		}
	}
	dc.config.FolderName = tempfolder
	dc.config.Password = config.Name
	dc.config.UserAction = config.LockFolder
	dc.config.InstructData = config.Instructions{
			DeleteOriginal: true,
	 	}
	return nil
}