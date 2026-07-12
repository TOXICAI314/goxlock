package mutex

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// - FolderMutex
// Creates a mutex onto the folder that `goxlock` is currently working on
type FolderMutex struct {
	handle 	windows.Handle
}

// - mutexName 
// mutexName : Create the Name of the mutex that is unique to the folder (it will be a hash)
func mutexname(folder string) (string,error) {

	// - Pre Safety 
	abs,err := filepath.Abs(folder) 
	if err != nil {
		return ``,err
	}

	potabs,err := filepath.EvalSymlinks(abs)
	if err == nil {
		abs = potabs
	}
	hash := sha256.Sum256([]byte(abs))

	return `Global\GoXLock_` + hex.EncodeToString(hash[:]),nil
}

// - NewMutex 
// Creates a new mutex for the folder specified
func NewMutex(object string) (folderMutex *FolderMutex,alreadyexists bool,err error)  {

	// Note : Use this mutexing careful -> One extra mutex marking and the app will crash on itslef
	// Currenly only being used in highly low level function `locker` & `unlocker`
	// And an upper function `verifypassword` (as it reads from the folder) (bypassig unlocker)

	name,err := mutexname(object)
	if err != nil {
		return nil,false,err
	}

	// Info : Making a pointer to uint16 type from the string that is given to us by the hash function
	ptr,err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil,false,err
	} 

	handle,callErr := windows.CreateMutex(nil,false,ptr)
	// Info : err is checked but also if its not what we want then it will get passed
	if handle == 0 {
		return nil,false,fmt.Errorf(`An Internal Error occured while making the mutex %+v`,callErr)
	}

	folderMutex = &FolderMutex{
		handle: handle,
	}

	// Info : alreadyexists is a guard rail for the mutex to work
	// const ERROR_ALREADY_EXISTS = 183
	alreadyexists = callErr == windows.ERROR_ALREADY_EXISTS

	return folderMutex,alreadyexists,nil
}