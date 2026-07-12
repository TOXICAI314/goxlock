package usersafety

import (
	"fmt"
	"goxlock/config"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	// getDiskFreeSpaceEx -> Is a proc that gives the free space details from the os
	getDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")
)

// A struct made for the disk space system which gets a secure way to complete the download
type DiskSpace struct {
	Path			string
	Volume 			string
	// Total space in the given volume
	TotalSpace 		uint64
	// The available for the drive
	AvailableSpace  uint64
	// The allocated for the current user in the drive
	FreeSpcae   	uint64
}

// Gets the disk space available to your deivce and make decision if the other data shall be made or not
func WindowsPartitionSpace(path string) (*DiskSpace,error) {

	var absPath string
	var err error
	if filepath.IsAbs(path) {
		absPath = path
	} else {
		absPath,err = filepath.Abs(path)
		if err != nil {
			return nil,&config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Cannot get the absolute path of the given path - %s`,path),
				ElapsedTime: time.Now(),
				Provider: `usersafety.WindowsPartitionProvider`,
			}
		}
	}

	volume := filepath.VolumeName(absPath)
	if volume == `` {
		return nil,&config.FunctionCancelError{
			Cause: `Empty Voulme`,
			Message: fmt.Sprintf(`Got an empty volume for - %s`,absPath),
			ElapsedTime: time.Now(),
			Provider: `usersafety.WindowsPartitionProvider`,
		}
	}

	// Info : This checks that the volume got the end backslashes for windows understandings
	if volume[len(volume)-1] != '\\' {
		volume += "\\"
	}

	pathPtr, err := windows.UTF16PtrFromString(volume)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot get the pointer of the volume string - %s`,absPath),
			ElapsedTime: time.Now(),
			Provider: `usersafety.WindowsPartitionProvider`,
		}
	}

	var availableBytes uint64
	var totalBytes uint64
	var freeBytes uint64

	r1,_,err := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&availableBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&freeBytes)),
	)

	if r1 == 0 {
		return nil,&config.FunctionFailError{
			Cause: err.Error(),
			Message: `The os program didnt went well as expected`,
			ElapsedTime: time.Now(),
			Provider: `usersafety.WindowsPartitionProvider`,
		}
	}

	return &DiskSpace{
		Path:        absPath,
		Volume:  	volume,
		AvailableSpace: availableBytes,
		TotalSpace:     totalBytes,
		FreeSpcae:      freeBytes,
	}, nil
}

// Checks the space for the data that is going to be stored
func CheckSpaceObject(cfg *config.Config) (error) {
	// Pre Safety
	if cfg == nil {
		return &config.FunctionCancelError{
			Cause: `Nil pointer reference`,
			Message: `Given a nil pointer of config for the space checking`,
			ElapsedTime: time.Now(),
			Provider: `usersafety.CheckSpaceObject`,
		}
	}

	stat,err := os.Stat(cfg.FolderName)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot get the stats of the object - %s`,cfg.FolderName),
			ElapsedTime: time.Now(),
			Provider: `usersafety.CheckSpaceObject`,
		}
	}

	diskSpace,err := WindowsPartitionSpace(cfg.OutputName)
	if err != nil {
		return err
	} 

	// Info : For now the only care is about free space
	// The 1.2 is the margin in bytes for the data to be cared of
	spaceRequired := uint64(float64(stat.Size()) * 1.2) 
	if diskSpace.FreeSpcae < spaceRequired {
		return &config.FunctionCancelError{
			Cause: fmt.Sprintf(`No disk Space left on the volumne %s`,diskSpace.Volume),
			Message: fmt.Sprintf(`The space required - %v is more than the space free for the user`,spaceRequired),
			ElapsedTime: time.Now(),
			Provider: `usersafety.CheckSpaceObject`,
		}
	}
	return nil
}