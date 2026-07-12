package performer

import (
	"encoding/base64"
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/header"
	"goxlock/core-fns/locker"
	"goxlock/core-fns/logger"
	passwordchanger "goxlock/core-fns/password-changer"
	"goxlock/core-fns/relocker"
	"goxlock/core-fns/unlocker"
	"goxlock/core-fns/verify"
	"time"
)

// Will translate user `int` action instructions -> Performable sections
func PerformAction(cfg *config.Config) (err error) {
	// Pre Safety 
	if cfg == nil {
		return &config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			ElapsedTime: time.Now(),
			Provider: `performer.PerformAction`,
		}
	} 
 
	defer func() {
		// Logger Entry
		// Logger logs here to put the end results
		if cfg.InstructData.LoggerAllowed {
			var loggerx *logger.Logger
			loggerx,logerr := logger.Log(cfg,err)
			if logerr != nil && err == nil {
				err = logerr
			}
			logerr = loggerx.Write()
			if logerr != nil && err == nil {
				err = logerr
			}
		}
	}()

	switch cfg.UserAction {
	case config.LockFolder:
		err =  locker.Locker(cfg)
		if err == nil {
			fmt.Println(`Password Integrity successful -> Remeber your password`)
			fmt.Println(`
			⚠ This password cannot be recovered.
			⚠ Store it safely in a password manager.
			⚠ goxlock is not responsible for forgotten password.
			`)		
		}
	case config.UnlockFolder:
		err = unlocker.Unlocker(cfg)
	case config.RelockFolder:
		err = relocker.Relocker(cfg.SessionID)
	case config.ChangePassword:
		err = passwordchanger.ChangePasswordForFolder(cfg)
		if err == nil {
			fmt.Println(`Password Integrity successful -> Remeber your password : If its lost , you cant recover the data`)
			fmt.Println(`
			⚠ This password cannot be recovered.
			⚠ Store it safely in a password manager.
			⚠ goxlock is not responsible for forgotten password.
			`)	
		}
	case config.VerifyPassword:
		err := verify.VerifyUnlock(cfg)
		if err == nil {
			// Info : Here the err is not returned -> it is used for soft message
			fmt.Println("Correct Password : Your given password is totally correct")
		} else {
			fmt.Printf(`Incorrect Password : The given password is wrong
				        Reason : %s`,err.Error())
		}
	case config.HeaderCheck:
		var header_x *config.Header
		header_x,err = header.Header(cfg.FolderName)
		if err == nil {
			salt := base64.StdEncoding.EncodeToString(header_x.Salt[:])
			nonce := base64.StdEncoding.EncodeToString(header_x.Nonce[:])
			fmt.Printf(
				`
				Name - %s
				Version - %v
				Salt (base64) - %s
				Nonce (base64) - %s
				`,string(header_x.Magic[:]),string(header_x.Version[:]),salt,nonce,
			)
		}
	default:
		err = fmt.Errorf("Unknow action `%d`", cfg.UserAction)
	}

	return err
}
