package config

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// CreateHeader : makes a header based on the Shared Encrypted data
func CreateHeader(end *SharedEncryptionData) *Header {
	return &Header{
		Magic: [7]byte([]byte(Name)),
		Version: [5]byte([]byte(Version)),
		Salt: [16]byte(end.Salt),
		Nonce: [12]byte(end.Nonce),
	}
}

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

	// Info : 40 is for the header -> See the data structure
	totalSize := 40	+ len(end.EncryptedData)
	packet := make([]byte, 0, totalSize)

	packet = append(packet, hd.Magic[:]...)
	packet = append(packet, hd.Version[:]...)
	packet = append(packet, hd.Salt[:]...)
	packet = append(packet, hd.Nonce[:]...)
	packet = append(packet, end.EncryptedData...)

	return packet,nil
}

// ReadHeader : Reads the valid `glock` header
func ReadHeaderAndRest(data []byte) (*Header,[]byte,error) {
	if len(data) == 0 {
		return nil,nil,&FunctionFailError{
			Cause: `Empty Data`,
			Message: `The given data length is zero`,
			Provider: `config.ReadHeaderAndRest`,
			ElapsedTime: time.Now(),
		}
	}

	dataReader := bytes.NewBuffer(data)
	var header *Header = &Header{}

	// Info : binary.Read reads the binary into already defied structure (as the structure of the Header is already defined in bytes)
	err := binary.Read(dataReader,binary.BigEndian,header)
	if err != nil {
		return nil,nil,&FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot store the header data into the Header struct`,
			ElapsedTime: time.Now(),
			Provider: `config.ReadHeaderAndRest`,
		}
	}

	// Info : And the rest of the data will be read (as the Reader pointer has been shifted)
	return header,dataReader.Bytes(),nil
}
 
// ValidateHeader : Gives a forward walking rights to the header if it is correct
func ValidateHeader(hd *Header) error {
	// PreSafety
	if hd == nil {
		return &FunctionFailError{
			Cause: `Nil pointer reference`,
			Message: `The given header pointer is pointing to nil`,
			Provider: `config.ValidateHeader`,
			ElapsedTime: time.Now(),
		}
	}
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
	return nil
}