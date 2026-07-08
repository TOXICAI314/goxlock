package performer

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/locker"
	"goxlock/core-fns/logger"
	passwordchanger "goxlock/core-fns/password-changer"
	"goxlock/core-fns/relocker"
	"goxlock/core-fns/unlocker"
	"goxlock/utils"
	"math"
	"os"
	"strings"
)

// - Const 
const (
	// deleteuppercursordata -> The ANSII keyword that deletes the upper written line from the buffer of the writter
	deleteuppercursordata string = "\033[A\033[K\033[A\033[K"
)

// -PerformAction
// Will translate user `int` action instructions -> Performable sections
func PerformAction(cfg *config.Config) error {
	// - Pre Safety 
	if cfg == nil {
		return &config.UserSafetyError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
		}
	} 

	var err error
	switch cfg.UserAction {
	case config.LockFolder:
		err =  locker.Locker(cfg)
		if err == nil {
			fmt.Println(`Password Integrity successful -> Remeber your password : If its lost , you cant recover the data`)
			ans := utils.Input(`Do you want to see the password ? (y/n): `)
			if strings.TrimSpace(strings.ToLower(ans)) == `y` {
				if strings.TrimSpace(strings.ToLower(ans)) == `y` {
					fmt.Printf("Your current password : %s\n",cfg.Password)
					fmt.Println(`This password is Visble to other users as it has gone to powershell buffer`)
					fmt.Println(`Press [ENTER] to erase it`)
					// Info : This will wait for the user to press the `ENTER` key to erase the password
					// Then an ANSII code string will be ran to move the cursor up
					bufio.NewReader(os.Stdin).ReadString('\n')
					fmt.Printf("%s%s",deleteuppercursordata,deleteuppercursordata) 
				}
			}
		}
	case config.UnlockFolder:
		err = unlocker.Unlocker(cfg)
	case config.RelockFolder:
		err = relocker.Relocker(cfg.SessionID)
	case config.ChangePassword:
		err = passwordchanger.ChangePasswordForFolder(cfg)
		if err == nil {
			fmt.Println(`Password Integrity successful -> Remeber your password : If its lost , you cant recover the data`)
			ans := utils.Input(`Do you want to see the password ? (y/n): `)
			if strings.TrimSpace(strings.ToLower(ans)) == `y` {
				fmt.Printf("Your current password : %s\n",cfg.ChangePassword.NewPassword)
				fmt.Println(`This password is Visble to other users as it has gone to powershell buffer`)
				fmt.Println(`Press [ENTER] to erase it`)
				// Info : This will wait for the user to press the `ENTER` key to erase the password
				// Then an ANSII code string will be ran to move the cursor up
				bufio.NewReader(os.Stdin).ReadString('\n')
				fmt.Printf("%s%s",deleteuppercursordata,deleteuppercursordata) 
			}
		}
	case config.VerifyPassword:
		err := unlocker.VerifyUnlock(cfg)
		if err == nil {
			// Info : Here the err is not returned -> it is used for soft message
			fmt.Println("Correct Password : Your given password is totally correct")
		} else {
			fmt.Printf(`Incorrect Password : The given password is wrong
				        Reason : %s`,err.Error())
		}
	case config.HeaderCheck:
		var header *config.Header
		header,err = unlocker.Header(cfg.FolderName)
		if err != nil {
			return err
		}
		salt := base64.StdEncoding.EncodeToString(header.Salt[:])
		nonce := base64.StdEncoding.EncodeToString(header.Nonce[:])
		version := math.Float32frombits(binary.LittleEndian.Uint32(header.Version[:]))
		fmt.Printf(
			`
			Name - %s
			Version - %v
			Salt (base64) - %s
			Nonce (base64) - %s
			`,string(header.Magic[:]),version,salt,nonce,
		)
		return nil
	default:
		err = fmt.Errorf("Unknow action `%d`", cfg.UserAction)
	}

	defer func() {
		// - Logger Entry
		// Logger logs here to put the end results
		var loggerx *logger.Logger
		loggerx,_ = logger.Log(cfg,err)
		loggerx.Write()
	}()

	return err
}
