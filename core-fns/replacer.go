package corefns

import (
	"goxlock/config"
	"os"
	"strings"
)

// - ReplaceZipwithGLock
// Replaces the zip file to a secure `.lock` structure
func ReplaceZipwithGLock(cfg *config.Config) error {
	
	if cfg == nil {
		return &config.UserSafetyError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
		}
	}
	
	zipfile := &cfg.OutputName
	target := &cfg.FolderName
	// - Pre Safety 

	if _,err := os.Stat(*zipfile);err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot find the zipped file`,
		}
	}
	if _,err := os.Stat(*target);err != nil {
		return &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot find the target file`,
		}
	}

	// - Replace
	if cfg.InstructData.DeleteOriginal {
		err := os.RemoveAll(*target)
		if err != nil {
			return err
		}
	}
	newZipfilename := strings.TrimSuffix(*zipfile,config.ZipExt)+ config.LockExt
	err := os.Rename(*zipfile,newZipfilename)
	if err != nil {
		return err
	}
	*zipfile = newZipfilename
	return nil
}