package header

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"goxlock/config"
	corefns "goxlock/core-fns"
	"os"
	"path/filepath"
	"time"
)

// Gives the raw header of the `g-lock` file
func Header(file string) (*config.Header,error) {
	// - Pre Safety
	if ext := filepath.Ext(file);ext != config.LockExt {
		return nil,&config.FunctionFailError{
			Cause: `Unwanted extension`,
			Message: fmt.Sprintf(`Wanted - %s ; Given - %s`,config.LockExt,ext),
			ElapsedTime: time.Now(),
			Provider: `header.Header`,
		}
	}

	O_file,err := os.Open(file)
	if err != nil {
		return nil,&config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot read from the file : %s`,file),
			ElapsedTime: time.Now(),
			Provider: `header.Header`,
		}
	}

	defer O_file.Close()
	var data [40]byte = [40]byte{}
	// Info : This reduces the overhead to read the whole file and then go forward
	_,err = O_file.ReadAt(data[:],0)
	if err != nil {
		return nil,&config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot Read from the data buffer to fill the header`,
			ElapsedTime: time.Now(),
			Provider: `header.Header`,
		}
	}

	var header *config.Header = &config.Header{}

	dataBuffer := bytes.NewBuffer(data[:])
	err = binary.Read(dataBuffer,binary.BigEndian,header)
	if err != nil {
		return nil,&config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot read to get the data into the header`,
			ElapsedTime: time.Now(),
			Provider: `header.Header`,
		}
	}

	return header,nil
}

// will verify the unlocked file and will give the data of the file that is needed
func GetUnlockedData(cfg *config.Config, rawData []byte) ([]byte, error) {
	// Pre Safety
	switch {
	case cfg == nil: 
		return nil,&config.FunctionCancelError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
			ElapsedTime: time.Now(),
			Provider: `header.GetUnlockedData`,
		}
	case  len(rawData) == 0:
		return nil,&config.FunctionCancelError{
			Cause: `Empty Data set`,
			Message: `A 0 lenght data slice is provided`,
			ElapsedTime: time.Now(),
			Provider: `header.GetUnlockedData`,
		}
	}
	
	header, encodeddata, err := config.ReadHeaderAndRest(rawData)
	if err != nil  {
		return nil,err
	}

	err = config.ValidateHeader(header)
	if err != nil {
		return nil, err
	}

	var sec *config.SharedEncryptionData = &config.SharedEncryptionData{
		Salt:          header.Salt,
		Nonce:         header.Nonce,
		EncryptedData: encodeddata,
	}

	plaindata, err := corefns.Decrypt(sec, cfg)
	if err != nil {
		return nil,err
	}

	return plaindata, nil
}