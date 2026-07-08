package config

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
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
		return nil,&UserSafetyError{
			Cause: `Nil pointer reference`,
			Message: `The given header pointer is pointing to nill`,
		}
	case end == nil:
		return nil,&UserSafetyError{
			Cause: `Nil pointer dereference`,
			Message: `A nil pointer of passed instead of a config pointer`,
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

// - ReadHeader() 
// ReadHeader : Reads the valid `glock` header
func ReadHeaderAndRest(data []byte) (*Header,[]byte,error) {
	if data == nil {
		return nil,nil,&UserSafetyError{
			Cause: `Nil pointer or empty data reference`,
			Message: `The given data slice cannot give any usefull thing to read header`,
		}
	}

	mainbuffer := bytes.NewBuffer(data)
	
	// - Reading 
	name := make([]byte,7)
	_,err := io.ReadFull(mainbuffer,name)
	if err != nil {
		return nil,nil,err
	}

	version := make([]byte,4)
	_,err = io.ReadFull(mainbuffer,version)
	if err != nil {
		return nil,nil,err
	}

	salt := make([]byte,16)
	_,err = io.ReadFull(mainbuffer,salt)
	if err != nil {
		return nil,nil,err
	}

	nonce := make([]byte,12)
	_,err = io.ReadFull(mainbuffer,nonce)
	if err != nil {
		return nil,nil,err
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
		return &DecryptionError{
			Cause: `Invalid Header Naming`,
			Message: fmt.Sprintf(`Naming of the '%s' file is not as intented`,ZipExt),
			Fix: fmt.Sprintf(`Make sure that the file is of current version
			Provided ->
			{
			name : %s
			}
			Needed ->
			{
			name : %s
			}
			`,name,Name),
		}
	}

	// - Version 
	// Info : must be a valid version number to actually run
	v := math.Float32frombits(binary.LittleEndian.Uint32(hd.Version[:]))
	if v < 1.0 {
		return &DecryptionError{
			Cause: `Invalid Version Number`,
			Message: fmt.Sprintf(`Version is not whats intended : %.2f != (>= 1.0)`,v),
			Fix: `Download the Best current version from my github 'TOXICAI314'`,
		}
	}

	return nil
}