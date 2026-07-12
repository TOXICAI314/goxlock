package corefns

import (
	"archive/zip"
	"bytes"
	"fmt"
	"goxlock/config"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// commit zip onto the given folder but preserving the tree structure
func Zip(cfg *config.Config) error {
	if cfg == nil {
		return &config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Zip`,
		}
	}
	
	zipName := &cfg.OutputName
	folder := &cfg.FolderName

	// Pre Safety 
	if _,err := os.Stat(*folder);err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`The folder given %s is not a valid path`,err),
			ElapsedTime: time.Now(),
			Provider: `corefns.Zip`,
		}
	}

	// Zipping
	zipFile,err := os.OpenFile(*zipName,os.O_CREATE|os.O_RDWR,0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: `The file failed to zip due to internal errors`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Zip`,
		}
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	// Walk :
	// - Walking down the whole tree structure and counting on every file or folder beneathe
	// - Memory intensive payload
	err = filepath.Walk(*folder,func(path string, info os.FileInfo, err error) error {
		relPath,err := filepath.Rel(*folder,path)
		// Info : `return nil` for continuing the recursion even if something fails
		// - For folder : `relPath` will handle the tree structuring
		if err != nil {
			err = &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Defect in walk by %s`,relPath),
				ElapsedTime: time.Now(),
				Provider: `corefns.Zip`,
			}
			fmt.Println(err.Error())
			return nil
		}
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		// Info : This will now will never attach to a symlink.
		// As symlink can cause recursion : To infite || to sensitive files 
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		// Exclusion 
		for _,exdExt := range cfg.InstructData.Exclusion {
			if ok,_ := filepath.Match(exdExt,filepath.Base(relPath));ok {
				return nil
			}
		}
		header,err := zip.FileInfoHeader(info)
		if err != nil {
			err = &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Defect in walk by the unredalble info of %s`,relPath),
				ElapsedTime: time.Now(),
				Provider: `corefns.Zip`,
			}
			fmt.Println(err.Error())
			return nil
		}
		header.Name = relPath

		openFile,err := os.Open(path)
		if err != nil {
			err = &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Cannot open the file : %s`,path),
				ElapsedTime: time.Now(),
				Provider: `corefns.Zip`,
			}
			fmt.Println(err.Error())
			return nil
		}
		w,err := zipWriter.CreateHeader(header)
		if err != nil {
			err = &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Cannot create the header for %s`,relPath),
				ElapsedTime: time.Now(),
				Provider: `corefns.Zip`,
			}
			fmt.Println(err.Error())
			openFile.Close()
			return nil
		}
		_,err = io.Copy(w,openFile)
		if err != nil {
			err = &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Cannot copy the data from %s`,path),
				ElapsedTime: time.Now(),
				Provider: `corefns.Zip`,
			}
			fmt.Println(err.Error())
			openFile.Close()
			return nil
		}
		openFile.Close()
		return nil
	})
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message : `Cannot zip the given folder structure`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Zip`,
		}
	}

	return nil
}

// Unzips the provided zip folder and makes a final folder out of that
func Unzip(cfg *config.Config,data []byte) error {
	// Pre Saefety
	switch {
	case len(data) == 0:
		return &config.FunctionCancelError{
			Cause: `Empty Data set`,
			Message: `A 0 length data slice is passed instead of actual data`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Unzip`,
		}
	case cfg == nil :
		return &config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Unzip`,
		}
	}
	// Info : This makes a new reader from the zip data provided
	zipPath := &cfg.OutputName
	outputDir := strings.TrimSuffix(
		*zipPath,
		filepath.Ext(*zipPath),
	)

	// Info : Then this zipp data get fetched into the reader to read from 
	reader, err := zip.NewReader(bytes.NewReader(data),int64(len(data)))
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: `The file failed to unzip due to internal errors`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Unzip`,
		}
	}

	// Main Unzipping
	// This unzips the file and secures its data into a folder structure
	
	for _, file := range reader.File {

		targetPath := filepath.Join(
			outputDir,
			file.Name,
		)

		err := os.MkdirAll(
			filepath.Dir(targetPath),
			0755,
		)
		if err != nil {
			err = &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Cannot make directory for %s`,targetPath),
				ElapsedTime: time.Now(),
				Provider: `corefns.Unzip`,
			}
			fmt.Println(err.Error())
			continue
		}

		src, err := file.Open()
		if err != nil {
			err = &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Cannot open the file %s`,file.Name),
				ElapsedTime: time.Now(),
				Provider: `corefns.Unzip`,
			}
			fmt.Println(err.Error())
			continue
		}

		dst, err := os.OpenFile(
			targetPath,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			file.Mode(),
		)
		if err != nil {
			src.Close()
			err = &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Cannot open the file to write the data : %s`,targetPath),
				ElapsedTime: time.Now(),
				Provider: `corefns.Unzip`,
			}
			fmt.Println(err.Error())
			continue
		}

		_, err = io.Copy(dst, src)

		dst.Close()
		src.Close()

		if err != nil {
			err = &config.FunctionFailError{
				Cause: err.Error(),
				Message: fmt.Sprintf(`Cannot copy from the Source %s to the destination %s`,file.Name,targetPath),
				ElapsedTime: time.Now(),
				Provider: `corefns.Unzip`,
			}
			fmt.Println(err.Error())
			continue
		}
	}

	return nil
}