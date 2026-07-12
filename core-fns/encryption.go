package corefns

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"goxlock/config"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/argon2"
)

// EncryptFileWithHeader : Encrypts the file info the undetectable form for user privacy that can only be breached by giving correct password
func EncryptFileWithHeader(cfg *config.Config) (err error) {
	if cfg == nil {
		return &config.FunctionCancelError{
			Cause: `Nil pointer reference`,
			Message: `A nil pointer to the configuration has been passed`,
			ElapsedTime: time.Now(),
			Provider: `corefns.EncryptFileWithHeader`,
		}
	}
	
	file := &cfg.OutputName

	openfile, err := os.OpenFile(*file, os.O_RDWR, 0700)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf("Cannot open given file : %s", filepath.Base(*file)),
			ElapsedTime: time.Now(),
			Provider: `corefns.EncryptFileWithHeader`,
		}
	}
	defer func() {
		closeErr := openfile.Close()
		// Info : This focuses on the more priotized error message of the original returner instead of the closing
		if closeErr != nil && err == nil {
			err = &config.FunctionFailError{
				Cause: closeErr.Error(),
				Message: fmt.Sprintf(`Error while closing the opened file - %s`,*file),
				ElapsedTime: time.Now(),
				Provider: `corefns.EncryptFileWithHeader`,
			}
		}
	}()
	data, err := io.ReadAll(openfile)
	if err != nil {
		return &config.FunctionFailError{
			Cause:   err.Error(),
			Message: fmt.Sprintf("Cannot read data from the given file : %s", filepath.Base(*file)),
			ElapsedTime: time.Now(),
			Provider: `corefns.EncryptFileWithHeader`,
		}
	}

	cipherdata, err := Encrypt(data, cfg)
	if err != nil {
		return err
	}

	// Header and Writing 
	header := config.CreateHeader(cipherdata)
	packet,err := config.CreatePacket(header,cipherdata)
	if err != nil {
		return err
	}

	_,err = openfile.Seek(0,0)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Cannot seek to the value of 0,0 of the opened file - %s`,*file),
			ElapsedTime: time.Now(),
			Provider: `corefns.EncryptFileWithHeader`,
		}
	}
	_,err = openfile.Write(packet)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`The file write to the %s has failed`,*file),
			ElapsedTime: time.Now(),
			Provider: `croefns.EncryptFileWithHeader`,
		}
	}
	return nil
}

// Encrypt : Given the bytes , it will encrypt with `aes` with valid protection of `gcm`
func Encrypt(data []byte, cfg *config.Config) (*config.SharedEncryptionData, error) {
	
	// Pre Safety
	switch {
	case cfg == nil:
		return nil,&config.FunctionCancelError{
			Cause: `Nil pointer reference`,
			Message: `A nil pointer to the configuration has been passed`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Encrypt`,
		}
	case len(data) == 0 :
		return nil,&config.FunctionCancelError{
			Cause: `Empty data set`,
			Message: `A 0 length data slice is given to encrypt`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Encrypt`,
		}
	}
	// Info :
	// salt -> a 16 bit byte array that acts as a parametre of extra entropy in encrypting
	// `rand.Read` makes the salt even more random
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot read random into the salt`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Encrypt`,
		}
	}

	// key -> the key used for encrypting the data (must be learnt by the user)
	// resultant : a perfect 32 bit code
	key := argon2.IDKey(
		[]byte(cfg.Password),
		salt,
		1,
		64*1024,
		4,
		32,
	)

	// block -> a new cipher for a secure encryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, &config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot create a new cipher`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Encrypt`,
		}
	}

	// gcm -> creates a bullet proof data structure which get denied in two conditions :
	// 1. Wrong password
	// 2. Alternation in the original data
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil,  &config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot create a new gcm`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Encrypt`,
		}
	}

	// nonce -> a one time used random byte array which will give random encrypt data on the same data
	// e1 != e2 for any encryption on same data
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil,  &config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot create a new nonce`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Encrypt`,
		}
	}

	// Encryption 
	// From here the encryption of the data starts
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	return &config.SharedEncryptionData{
		Salt:          [16]byte(salt),
		Nonce:         [12]byte(nonce),
		EncryptedData: ciphertext,
	}, nil
}

// decrypts the file and make it an array of bytes that can be seen
func Decrypt(sec *config.SharedEncryptionData,cfg *config.Config) ([]byte,error) {
	// Pre Safety 
	switch {
	case cfg == nil:
		return nil,&config.FunctionCancelError{
			Cause: `Nil pointer reference`,
			Message: `A nil pointer to the configuration has been passed`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Decrypt`,
		}
	case sec == nil :
		return nil,&config.FunctionCancelError{
			Cause: `Nil pointer reference`,
			Message: `A nil pointer to the Shared Encryption data has been passed`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Decrypt`,
		}
	}

	key := argon2.IDKey(
		[]byte(cfg.Password),
		sec.Salt[:],
		1,
		64*1024,
		4,
		32,
	)

	// AES Gcm open 
	block,err := aes.NewCipher(key)
	if err != nil {
		return nil,&config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot create a new cipher`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Decrypt`,
		}
	}

	gcm,err := cipher.NewGCM(block)
	if err != nil {
		return nil,&config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot create a new gcm`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Decrypt`,
		}
	}

	plaindata,err := gcm.Open(nil,sec.Nonce[:],sec.EncryptedData,nil)
	if err != nil {
		return nil,&config.FunctionFailError{
			Cause: err.Error(),
			Message: `Cannot make decrypt the data because of the given reason`,
			ElapsedTime: time.Now(),
			Provider: `corefns.Decrypt`,
		}
	} 

	return plaindata,nil
}