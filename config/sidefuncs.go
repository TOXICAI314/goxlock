package config

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"
)

// - CreateHeader()
// CreateHeader : makes a header based on the Shared Encrypted data
func CreateHeader(end *SharedEncryptionData) *Header {

	// -- Version -- //
	versionbytes := [4]byte{}
	binary.LittleEndian.PutUint32(versionbytes[:],math.Float32bits(Version))

	return &Header{
		Magic: [7]byte([]byte(Name)),
		Version: versionbytes,
		Salt: [16]byte(end.Salt),
		Nonce: [12]byte(end.Nonce),
	}
}

// - CreatePacket
// CreatePacket : Creates a header package that can get written into the file directly
func CreatePacket(hd *Header,end *SharedEncryptionData) ([]byte,error) {
	switch {
	case hd == nil:
		return nil,&FunctionCancelError{
			Cause: `Nil pointer reference`,
			Message: `The given header pointer is pointing to nil`,
			Provider: `config.CreatePacket`,
			ElapsedTime: time.Now(),
		}
	case end == nil:
		return nil,&FunctionCancelError{
			Cause: `Nil pointer reference`,
			Message: `The given config pointer is pointing to nil`,
			Provider: `config.CreatePacket`,
			ElapsedTime: time.Now(),
		}
	}

	// Info : 39 is for the header -> See the data structure
	totalSize := 39	+ len(end.EncryptedData)
	packet := make([]byte, 0, totalSize)

	packet = append(packet, hd.Magic[:]...)
	packet = append(packet, hd.Version[:]...)
	packet = append(packet, hd.Salt[:]...)
	packet = append(packet, hd.Nonce[:]...)
	packet = append(packet, end.EncryptedData...)

	return packet,nil
}

// - ReadHeader
// ReadHeader : Reads the valid `glock` header
func ReadHeaderAndRest(data []byte) (*Header,[]byte,error) {
	if data == nil {
		return nil,nil,&FunctionFailError{
			Cause: `Nil pointer reference`,
			Message: `The given data pointer is pointing to nil`,
			Provider: `config.ReadHeaderAndRest`,
			ElapsedTime: time.Now(),
		}
	}

	mainbuffer := bytes.NewBuffer(data)
	
	// - Reading 
	name := make([]byte,7)
	_,err := io.ReadFull(mainbuffer,name)
	if err != nil {
		return nil,nil,&FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot read from the main buffer to the name storing buffer`,
			ElapsedTime: time.Now(),
			Provider: `config.ReadHeaderAndRest`,
		}
	}

	version := make([]byte,4)
	_,err = io.ReadFull(mainbuffer,version)
	if err != nil {
		return nil,nil,&FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot read from the main buffer to the version storing buffer`,
			ElapsedTime: time.Now(),
			Provider: `config.ReadHeaderAndRest`,
		}
	}

	salt := make([]byte,16)
	_,err = io.ReadFull(mainbuffer,salt)
	if err != nil {
		return nil,nil,&FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot read from the main buffer to the salt storing buffer`,
			ElapsedTime: time.Now(),
			Provider: `config.ReadHeaderAndRest`,
		}
	}

	nonce := make([]byte,12)
	_,err = io.ReadFull(mainbuffer,nonce)
	if err != nil {
		return nil,nil,&FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot read from the main buffer to the nonce storing buffer`,
			ElapsedTime: time.Now(),
			Provider: `config.ReadHeaderAndRest`,
		}
	}

	return &Header{
		Magic: [7]byte(name),
		Version: [4]byte(version),
		Salt: [16]byte(salt),
		Nonce: [12]byte(nonce),
	},mainbuffer.Bytes(),nil
}

// - ValidateHeader 
// ValidateHeader : Gives a forward walking rights to the header if it is correct
func ValidateHeader(hd *Header) error {
	// - Name 
	name := string(hd.Magic[:])
	if name != Name {
		return &FunctionCancelError{
			Cause: `Invalid Header Naming`,
			Message: fmt.Sprintf(`Naming of the '%s' file is not as intented`,ZipExt),
			ElapsedTime: time.Now(),
			Provider: `config.ValidateHeader`,
		}
	}

	// - Version 
	// Info : must be a valid version number to actually run
	v := math.Float32frombits(binary.LittleEndian.Uint32(hd.Version[:]))
	if v < 1.0 {
		return &FunctionCancelError{
			Cause: `Invalid Version Number`,
			Message: fmt.Sprintf(`Version is not whats intended : %.2f != (>= 1.0)`,v),
			ElapsedTime: time.Now(),
			Provider: `config.ValidateHeader`,
		}
	}

	return nil
}