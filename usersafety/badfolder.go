package usersafety

import (
	"goxlock/config"
	"os"
	"path/filepath"
	"time"
)

// BadFolders -> The folder which can never be changed even if user want
var BadFolders = map[string]bool{}

func InitializeBadFolders() error {
	// Pre Safety
	userHome,err := os.UserHomeDir()
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: `The Home dir cannot be initialized to form the bad folder security`,
			ElapsedTime: time.Now(),
			Provider: `usersafety.InitializeBadFolder`,
		}
	}

	BadFolders = map[string]bool{
		// Empty
		"": true,

		// Drive Roots
		`C:\`: true,
		`D:\`: true,
		`E:\`: true,
		`F:\`: true,

		// Windows Core
		`C:\Windows`:                  true,
		`C:\Windows\System32`:         true,
		`C:\Windows\SysWOW64`:         true,
		`C:\Windows\WinSxS`:           true,
		`C:\Windows\Boot`:             true,
		`C:\Windows\Fonts`:            true,
		`C:\Windows\Installer`:        true,
		`C:\Windows\Security`:         true,
		`C:\Windows\SystemResources`:  true,
		`C:\Windows\servicing`:        true,

		// Program Installation
		`C:\Program Files`:            true,
		`C:\Program Files (x86)`:      true,
		`C:\ProgramData`:              true,

		// Recovery & System
		`C:\Recovery`:                 true,
		`C:\System Volume Information`: true,
		`C:\$Recycle.Bin`:             true,

		// User Profile Critical
		`C:\Users\Default`:            true,
		`C:\Users\Default User`:       true,
		`C:\Users\Public`:             true,

		// AppData 
		filepath.Join(userHome, "AppData")             :  true,
		filepath.Join(userHome, "AppData", "Local")    :  true,
		filepath.Join(userHome, "AppData", "Roaming")  :  true,
		filepath.Join(userHome, "AppData", "LocalLow") :  true,

		// Windows Apps
		`C:\Program Files\WindowsApps`: true,

		// Defender
		`C:\ProgramData\Microsoft\Windows Defender`: true,

		// Driver Store
		`C:\Windows\System32\DriverStore`: true,

		// EFI (if mounted)
		`C:\EFI`: true,
	}
	return nil
}
