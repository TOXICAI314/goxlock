package locker

import (
	"fmt"
	"goxlock/config"
	corefns "goxlock/core-fns"
	"goxlock/core-fns/mutex"
	"goxlock/usersafety"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Will the lock the file in its `.g-lock` via encryption of the data
// Once the data has been encrypted the only way to decrypt is via the password and nothing else
func Locker(cfg *config.Config) error {

	if cfg == nil {
		return &config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			Provider: `locker.Locker`,
			ElapsedTime: time.Now(),
		}
	} 
	// Pre Safety 
	stats, err := os.Stat(cfg.FolderName)
	if err == nil {
		if !stats.IsDir() {
			return &config.FunctionCancelError{
				Cause: `Non Folder file`,
				Message: fmt.Sprintf("Cannot Walk over a file like structure : %s", cfg.FolderName),
				Provider: `locker.Locker`,
				ElapsedTime: time.Now(),
			}
		}
	} else {
		return &config.FunctionCancelError{
			Cause: err.Error(),
			Message: fmt.Sprintf("Error while confirming the given path : %s - %s", cfg.FolderName, err.Error()),
			Provider: `locker.Locker`,
			ElapsedTime: time.Now(),
		}
	}

	// Info : For safe zipping the following part will be done
	selfDir := filepath.Dir(cfg.OutputName)
	requiredFolderCreation := ``
	if !filepath.IsAbs(cfg.OutputName) {
		folderParent := filepath.Dir(cfg.FolderName)
		relativeWorkingFolder := filepath.Join(folderParent,selfDir)
		requiredFolderCreation = relativeWorkingFolder
		cfg.OutputName = filepath.Join(relativeWorkingFolder,strings.TrimSuffix(filepath.Base(cfg.OutputName), filepath.Ext(cfg.OutputName))+config.ZipExt)
		
	} else {
		zipFolderName := strings.TrimSuffix(filepath.Base(cfg.OutputName), filepath.Ext(cfg.OutputName)) + config.ZipExt
		requiredFolderCreation = selfDir
		cfg.OutputName = filepath.Join(filepath.Dir(cfg.OutputName),zipFolderName)
	}
	err = os.MkdirAll(requiredFolderCreation,0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot create the folder %s for the locking data to be stored`,selfDir),
			Provider: `locker.Locker`,
			ElapsedTime: time.Now(),
		}
	}

	// Mutex Locking 
	// Means getting the mutex for the writing stuff of the data  
	// The mutex is from the os and cant be penetrated easily

	mut, alrexist, err := mutex.NewMutex(cfg.FolderName)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`An internal error has occured while Creating the Mutex for the given folder - %s`,cfg.FolderName),
		}
	}
	defer mut.CloseMutex()
	if alrexist {
		return &config.FunctionCancelError{
			Cause:   `Mutex already exist`,
			Message: fmt.Sprintf(`The mutex of the given folder %s is already there`, cfg.FolderName),
			Provider: `locker.Locker`,
			ElapsedTime: time.Now(),
		}
	}

	if !cfg.InstructData.UnSafe {
		// Info : Folder security check
		err = usersafety.SecureFolder(cfg)
		if err != nil {
			return err
		}
		// Info : Space check for the device
		err = usersafety.CheckSpaceObject(cfg)
		if err != nil {
			return err
		}
	}

	// Info  : Its only for presentation  purpose 
	// Its gives the time where the main zipping and encryption of the data starts
	benchmarktimestart := time.Now()


	err = corefns.Zip(cfg)
	if err != nil {
		return err
	}

	err = corefns.EncryptFileWithHeader(cfg)
	if err != nil {
		return err
	}

	// Info :  This is just done to save the efficiency of the logger
	// `if logger.Allowed` is used for the efficiency
	if cfg.InstructData.Stats {
		var count int 
		var foldersize int64
		err = filepath.Walk(cfg.FolderName,func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			foldersize += info.Size()
			count++
			return nil
		})
		elapsedTime := time.Since(benchmarktimestart)
		
		msg := fmt.Sprintf(`
				--- --- ---
				Subject Name : %s
				Folder Size : %d B
				Folder Material count : %d
				Elapsed Time : %s
				Average Speed : %.4f B/seconds
				--- --- ---
			`,cfg.FolderName,foldersize,count,elapsedTime.String(),(float64(foldersize))/elapsedTime.Seconds())
		fmt.Println(msg)
	}

	// Changes 
	err = corefns.ReplaceZipwithGLock(cfg)
	if err != nil {
		return err
	}

	return nil
}
