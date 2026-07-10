package doctor

import (
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
	"time"
)

// - CreateLockConfig
// This creates a lock for the doctor to check -> Everything will be in temp folder
func (dc *Doctor) CreateTestLockerConfig() (error) {
	tempfolder,err := os.MkdirTemp(``,fmt.Sprintf(`%s-test-%s`,config.Name,dc.config.SessionID))
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot make a temp folder for the test`,
			ElapsedTime: time.Now(),
			Provider: `doctoy.Doctor.CreateTestLockerConfig`,
		}
	}
	r1 := filepath.Join(tempfolder,`r1.txt`)
	err = os.WriteFile(r1,[]byte(`Wassup man , Having a great day ??`),0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot make write into the test folder : %s`,r1),
			ElapsedTime: time.Now(),
			Provider: `doctoy.Doctor.CreateTestLockerConfig`,
		}
	}
	r2 := filepath.Join(tempfolder,`r2.txt`)
	err = os.WriteFile(r2,[]byte(`Wassup man , Having a great day ??`),0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot make write into the test folder : %s`,r2),
			ElapsedTime: time.Now(),
			Provider: `doctoy.Doctor.CreateTestLockerConfig`,
		}
	}
	r3 := filepath.Join(tempfolder,`r3.txt`)
	err = os.WriteFile(r3,[]byte(`Wassup man , Having a great day ??`),0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot make write into the test folder : %s`,r3),
			ElapsedTime: time.Now(),
			Provider: `doctoy.Doctor.CreateTestLockerConfig`,
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