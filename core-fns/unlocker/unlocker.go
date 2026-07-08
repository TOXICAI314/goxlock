package unlocker

import (
	"fmt"
	"goxlock/config"
	corefns "goxlock/core-fns"
	"goxlock/core-fns/dpapi"
	"goxlock/core-fns/mutex"
	"goxlock/core-fns/scheduler"
	"goxlock/core-fns/sessions"
	"goxlock/usersafety"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// - Unlocker
// unlocks the locked extension file -> only if the header is correct && file has not been altered
func Unlocker(cfg *config.Config) error {
	if cfg == nil {
		return &config.UserSafetyError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
		}
	}
	targetglockfile := &cfg.FolderName
	outFolder := &cfg.OutputName

	// Info : Protected psasword can be used globally to spread across multiple sessions
	protectedPassword, err := dpapi.Protect([]byte(cfg.Password))
	if err != nil {
		return err
	}

	// - Pre Safety 
	if ext := filepath.Ext(*targetglockfile); ext != config.LockExt {
		return &config.DecryptionError{
			Cause:   `Invalid extension`,
			Message: fmt.Sprintf("goxlock cant decrypt the data that is not in its native extension - %s", config.LockExt),
			Fix: fmt.Sprintf(
				`
			Provided ->
			{
			ext : %s
			}
			Needed ->
			{
			ext : %s
			}`, ext, config.LockExt),
		}
	}
	if _,err := os.Stat(*targetglockfile);err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Invalid target file which stats cant be fetched out`,
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
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot create the folder %s for the unlocking data to be stored`,selfDir),
		}
	} 

	// - Mutex Locking 
	mut, alrexist, err := mutex.NewMutex(*targetglockfile)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: `An internal error has occured while Creating the Mutex for the given folder`,
		}
	}
	defer mut.CloseMutex()
	if alrexist {
		return &config.UserSafetyError{
			Cause:   `Mutex already exist`,
			Message: fmt.Sprintf(`The mutex of the given folder %s is already there`, cfg.FolderName),
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
	// - Timeout Saefty 
	if cfg.InstructData.Timeout > 0 {
		if cfg.SessionID == `` {
				return &config.UserSafetyError{
					Cause:   `Cannot create a session out of an empty string`,
					Message: `The given config Session is invalid as it is empty`,
				}
			}

		scheduleID := cfg.SessionID

		s := &sessions.Session{
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
		return &config.UnzipError{
			Cause:   err.Error(),
			Message: `Cannot read from the file`,
			Fix: `
			Make sure the file is not:
			1. Opened by too many application
			2. Os restricted file
			3. Self laid retrictions
			`,
		}
	}

	plaindata, err := GetUnlockedData(cfg, data)
	if err != nil  {
		return err
	}

	err = corefns.Unzip(cfg, plaindata)
	if err != nil {
		return err
	}

	if cfg.InstructData.DeleteOriginal {
		err = os.Remove(cfg.FolderName)
		if err != nil {
			return &config.UserSafetyError{
				Cause:   err.Error(),
				Message: fmt.Sprintf(`System refused to delete the resedue : %s`, cfg.FolderName),
			}
		}
	}

	newZipfilename := strings.TrimSuffix(*outFolder, filepath.Ext(*outFolder))
	*outFolder = newZipfilename

	if cfg.InstructData.Stats {
		stats,_ := os.Stat(cfg.FolderName)
		foldersize := stats.Size()
		elapsedTime := time.Since(benchmarktimestart)

		msg := fmt.Sprintf(`
				--- --- ---
				Subject Name : %s
				Folder Size : %d
				Folder Material count : 1
				Elapsed Time : %s
				Average Speed : %.4f MB/seconds
				--- --- ---
			`,cfg.FolderName,foldersize,elapsedTime.String(),(float64(foldersize)/(1024 * 1024))/elapsedTime.Seconds())
		fmt.Println(msg)
	}
	return nil
}

// - GetUnlockedData 
// will verify the unlocked file and will give the data of the file that is needed
func GetUnlockedData(cfg *config.Config, rawData []byte) ([]byte, error) {
	// - Pre Safety
	switch {
	case cfg == nil: 
		return nil,&config.UserSafetyError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
		}
	case rawData == nil :
		return nil,&config.UserSafetyError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a raw data pointer`,
		}
	}
	
	header, encodeddata, err := config.ReadHeaderAndRest(rawData)
	if err != nil  {
		return nil, &config.UnzipError{
			Cause:   err.Error(),
			Message: `Cannot get header from the file`,
			Fix: `
			Make sure your file is not:
			1. Altered by anyone
			2. Made private or not allowing unlocking
			`,
		}
	}

	err = config.ValidateHeader(header)
	if err != nil {
		return nil, err
	}

	var sec *config.SharedEncryptionData = &config.SharedEncryptionData{
		Salt:          header.Salt[:],
		Nonce:         header.Nonce[:],
		EncryptedData: encodeddata,
	}

	plaindata, err := corefns.Decrypt(sec, cfg)
	if err != nil {
		return nil,err
	}

	return plaindata, nil
}

// - VerifyUnlock
// Will verify the given password and will provide error if there are any
func VerifyUnlock(cfg *config.Config) error {
	// - Pre Safety
	if cfg == nil {
		return &config.UserSafetyError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
		}
	} 
	if ext := filepath.Ext(cfg.FolderName); ext != config.LockExt {
		return &config.UserSafetyError{
			Cause:   fmt.Sprintf(`Cannot Verify a non '%s' folder`, config.LockExt),
			Message: fmt.Sprintf(`The given extension is %s ; needed -> %s`, ext, config.LockExt),
		}
	}

	data, err := os.ReadFile(cfg.FolderName)
	if err != nil {
		return &config.UserSafetyError{
			Cause:   err.Error(),
			Message: fmt.Sprintf(`Cannot read from the given file - %s`, cfg.FolderName),
		}
	}
	_, err = GetUnlockedData(cfg, data)
	if err != nil {
		return err
	}

	return nil
}

// - Header
// Gives the raw header of the `g-lock` file
func Header(file string) (*config.Header,error) {
	// - Pre Safety
	if ext := filepath.Ext(file);ext != config.LockExt {
		return nil,&config.UserSafetyError {
			Cause: `Unwanted extension`,
			Message: fmt.Sprintf(`Wanted - %s ; Given - %s`,config.LockExt,ext),
		}
	}

	data,err := os.ReadFile(file)
	if err != nil {
		return nil,&config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot read from the file`,
		}
	}

	header,_,err := config.ReadHeaderAndRest(data)
	if err != nil {
		return nil,err
	}

	return header,nil
}