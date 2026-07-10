package passwordchanger

import (
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/locker"
	"goxlock/core-fns/unlocker"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ChangePasswordForFolder(cfg *config.Config) (error) {
	if cfg == nil {
		return &config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			ElapsedTime: time.Now(),
			Provider: `passwordchanger.ChangePasswordForFolder`,
		}
	} 

	encryptedfolder := cfg.FolderName
	oldpassword := cfg.Password
	newpassword := cfg.ChangePassword.NewPassword

	
	// - Pre Safety 
	if _,err := os.Stat(encryptedfolder);err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot get stats of the given folder to Change the password : %s`,encryptedfolder),
			ElapsedTime: time.Now(),
			Provider: `passwordchanger.ChangePasswordForFolder`,
		}
	}
	if ext := filepath.Ext(encryptedfolder);ext != config.LockExt {
		return &config.FunctionCancelError{
			Cause: `Wrong extension`,
			Message: fmt.Sprintf(`Make sure that the extension is named as : %s`,config.LockExt),
			ElapsedTime: time.Now(),
			Provider: `passwordchanger.ChangePasswordForFolder`,
		}
	}

	// - config formation 
	// As the function is standalone and dont get a fullfilled config on it own by the upper commands
	// The code have to build its own
	u_cfg := config.Config{
		FolderName: encryptedfolder,
		Password: oldpassword,
		OutputName: strings.TrimSuffix(filepath.Base(encryptedfolder),config.LockExt),
		UserAction: config.UnlockFolder,
		InstructData: config.Instructions{
			DeleteOriginal: true,
		},
	}

	err := unlocker.Unlocker(&u_cfg)
	if err != nil {
		return err
	}

	// - Relocking config 
	ru_cfg := config.Config{
		FolderName: u_cfg.OutputName,
		Password: newpassword,
		OutputName: encryptedfolder,
		UserAction: config.LockFolder,
		InstructData: config.Instructions{
			DeleteOriginal: true,
		},
	}
	err = locker.Locker(&ru_cfg)
	if err != nil {
		return err
	}

	return nil
}