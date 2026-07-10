package corefns

import (
	"fmt"
	"goxlock/config"
	"os"
	"strings"
	"time"
)

// - ReplaceZipwithGLock
// Replaces the zip file to a secure `.lock` structure
func ReplaceZipwithGLock(cfg *config.Config) error {
	
	if cfg == nil {
		return &config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			ElapsedTime: time.Now(),
			Provider: `corefns.ReplaceZipWithGlock`,
		}
	}
	
	zipfile := &cfg.OutputName
	target := &cfg.FolderName
	// - Pre Safety 

	if _,err := os.Stat(*zipfile);err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot find the zipped file : %s`,*zipfile),
			ElapsedTime: time.Now(),
			Provider: `corefns.ReplaceZipWithGlock`,
		}
	}
	if _,err := os.Stat(*target);err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot find the target file : %s`,*target),
			ElapsedTime: time.Now(),
			Provider: `corefns.ReplaceZipWithGlock`,
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