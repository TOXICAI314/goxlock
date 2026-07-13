package unlocker

import (
	"fmt"
	"goxlock/config"
	corefns "goxlock/core-fns"
	"goxlock/core-fns/dpapi"
	"goxlock/core-fns/header"
	"goxlock/core-fns/mutex"
	"goxlock/core-fns/scheduler"
	"goxlock/core-fns/sessions"
	"goxlock/usersafety"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// unlocks the locked extension file -> only if the header is correct && file has not been altered
func Unlocker(cfg *config.Config) error {
	if cfg == nil {
		return &config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			ElapsedTime: time.Now(),
			Provider: `unlocker.Unlocker`,
		}
	}
	targetglockfile := &cfg.FolderName
	outFolder := &cfg.OutputName

	// Info : Protected psasword can be used globally to spread across multiple sessions
	protectedPassword, err := dpapi.Protect([]byte(cfg.Password))
	if err != nil {
		return err
	}

	// Pre Safety 
	if ext := filepath.Ext(*targetglockfile); ext != config.LockExt {
		return &config.FunctionCancelError{
			Cause:   `Invalid extension`,
			Message: fmt.Sprintf("goxlock cant decrypt the data that is not in its native extension - %s", config.LockExt),
			ElapsedTime: time.Now(),
			Provider: `unlocker.Unlocker`,
		}
	}
	if _,err := os.Stat(*targetglockfile);err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Invalid target file which stats cant be fetched out : %s`,*targetglockfile),
			ElapsedTime: time.Now(),
			Provider: `unlocker.Unlocker`,
		}
	}

	selfDir := filepath.Dir(*outFolder)
	requiredFolderCreation := ``
	if !filepath.IsAbs(*outFolder) {
		folderParent := filepath.Dir(*targetglockfile)
		relativeWorkingFolder := filepath.Join(folderParent,selfDir)
		requiredFolderCreation = relativeWorkingFolder
		*outFolder = filepath.Join(relativeWorkingFolder,strings.TrimSuffix(filepath.Base(*outFolder), config.LockExt)+config.ZipExt)
	} else {
		zipFolderName := strings.TrimSuffix(filepath.Base(*outFolder), config.LockExt) + config.ZipExt
		requiredFolderCreation = selfDir
		cfg.OutputName = filepath.Join(filepath.Dir(cfg.OutputName),zipFolderName)
	}
	err = os.MkdirAll(requiredFolderCreation,0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot create the folder %s for the unlocking data to be stored`,selfDir),
			ElapsedTime: time.Now(),
			Provider: `unlocker.Unlocker`,
		}
	} 

	// Mutex Locking 
	mut, alrexist, err := mutex.NewMutex(*targetglockfile)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`An internal error has occured while Creating the Mutex for the given folder : %s`,*targetglockfile),
			ElapsedTime: time.Now(),
			Provider: `unlocker.Unlocker`,
		}
	}
	defer func() {
		closeErr := mut.CloseMutex()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if alrexist {
		return &config.FunctionCancelError{
			Cause:   `Mutex already exist`,
			Message: fmt.Sprintf(`The mutex of the given folder %s is already there`, cfg.FolderName),
			ElapsedTime: time.Now(),
			Provider: `unlocker.Unlocker`,
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

	benchmarktimestart := time.Now()
	// Timeout Saefty 
	if cfg.InstructData.Timeout > 0 {
		if cfg.SessionID == `` {
				return &config.FunctionCancelError{
					Cause:   `Cannot create a session out of an empty string`,
					Message: `The given config Session is invalid as it is empty`,
					ElapsedTime: time.Now(),
				Provider: `unlocker.Unlocker`,
				}
			}

		scheduleID := cfg.SessionID

		s := &sessions.Session{
			VersionInfo: 	cfg.Version,
			Id:             scheduleID,
			Folder:         strings.TrimSuffix(*outFolder, config.ZipExt),
			Password:       protectedPassword,
			InstructionData: cfg.InstructData,
			OutputName:     cfg.FolderName,
		}
		err = scheduler.CreateShedule(scheduleID, cfg.InstructData)
		if err != nil {
			return err
		}
		err = s.CreateSession()
		if err != nil {
			return err
		}
	}

	lockfile := &cfg.FolderName
	data, err := os.ReadFile(*lockfile)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot read from the locked file : %s`,*lockfile),
			ElapsedTime: time.Now(),
			Provider: `unlocker.Unlocker`,
		}
	}

	plaindata, err := header.GetUnlockedData(cfg, data)
	if err != nil  {
		return err
	}

	err = corefns.Unzip(cfg, plaindata)
	if err != nil {
		return err
	}

	newZipfilename := strings.TrimSuffix(*outFolder, filepath.Ext(*outFolder))
	*outFolder = newZipfilename

	if cfg.InstructData.Stats {
		stats,err  := os.Stat(cfg.FolderName)
		if err != nil {
			return &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Cannot get stats for the folder : %s`,cfg.FolderName),
				ElapsedTime: time.Now(),
				Provider: `unlocker.Unlocker`,
			}
		}
		foldersize := stats.Size()
		elapsedTime := time.Since(benchmarktimestart)

		msg := fmt.Sprintf(`
				--- --- ---
				Subject Name : %s
				Folder Size : %d B
				Folder Material count : 1
				Elapsed Time : %s
				Average Speed : %.4f B/seconds
				--- --- ---
			`,cfg.FolderName,foldersize,elapsedTime.String(),(float64(foldersize))/elapsedTime.Seconds())
		fmt.Println(msg)
	}

	if cfg.InstructData.DeleteOriginal {
		err = os.Remove(cfg.FolderName)
		if err != nil {
			return &config.FunctionFailError{
				Cause:   err.Error(),
				Message: fmt.Sprintf(`System refused to delete the resedue : %s`, cfg.FolderName),
				ElapsedTime: time.Now(),
				Provider: `unlocker.Unlocker`,
			}
		}
	}
	return nil
}