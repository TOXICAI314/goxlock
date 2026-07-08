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
	"golang.org/x/crypto/argon2"
)

// - EncryptFileWithHeader
// EncryptFileWithHeader : Encrypts the file info the undetectable form for user privacy that can only be breached by giving correct password
func EncryptFileWithHeader(cfg *config.Config) error {
	if cfg == nil {
		return &config.UserSafetyError{
			Cause: `Nil pointer reference`,
			Message: `A nil pointer to the configuration has been passed`,
		}
	}
	
	file := &cfg.OutputName

	openfile, err := os.OpenFile(*file, os.O_RDWR, 0700)
	if err != nil {
		return &config.EncryptionError{
			Cause:   err.Error(),
			Message: fmt.Sprintf("Cannot open given file : %s", filepath.Base(*file)),
			Fix: `
			Make sure that the given file is not:
			1. Deleted
			2. Locked by a system mutex
			3. Opened by too many application
			`,
		}
	}
	defer openfile.Close()
	data, err := io.ReadAll(openfile)
	if err != nil {
		return &config.EncryptionError{
			Cause:   err.Error(),
			Message: fmt.Sprintf("Cannot read data from the given file : %s", filepath.Base(*file)),
			Fix: `
			Make sure that the given file is not:
			1. Deleted
			2. Being written over
			3. Locked by a system mutex
			4. Too much reads
			`,
		}
	}

	cipherdata, err := Encrypt(data, cfg)
	if err != nil {
		return err
	}

	// - Header and Writing 
	header := config.CreateHeader(cipherdata)
	packet,err := config.CreatePacket(header,cipherdata)
	if err != nil {
		return err
	}
	// Info : File pointer exhausted by read -> already EOF , hence direct writeFile
	os.WriteFile(*file,packet,0700)
	return nil
}

// - Encrypt
// Encrypt : Given the bytes , it will encrypt with `aes` with valid protection of `gcm`
func Encrypt(data []byte, cfg *config.Config) (*config.SharedEncryptionData, error) {
	
	// - Pre Safety
	switch {
	case cfg == nil:
		return nil,&config.UserSafetyError{
			Cause: `Nil pointer reference`,
			Message: `A nil pointer to the configuration has been passed`,
		}
	case data == nil :
		return nil,&config.UserSafetyError{
			Cause: `Nil pointer reference`,
			Message: `A nil pointer to the byte data has been passed`,
		}
	}
	// Info :
	// salt -> a 16 bit byte array that acts as a parametre of extra entropy in encrypting
	// `rand.Read` makes the salt even more random
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot read random into the salt`,
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
		return nil, &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot create a new cipher`,
		}
	}

	// gcm -> creates a bullet proof data structure which get denied in two conditions :
	// 1. Wrong password
	// 2. Alternation in the original data
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil,  &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot create a new gcm`,
		}
	}

	// nonce -> a one time used random byte array which will give random encrypt data on the same data
	// e1 != e2 for any encryption on same data
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil,  &config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot create a new nonce`,
		}
	}

	// - Encryption 
	// From here the encryption of the data starts
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	return &config.SharedEncryptionData{
		Salt:          salt,
		Nonce:         nonce,
		EncryptedData: ciphertext,
	}, nil
}

// - Decrypt
// decrypts the file and make it an array of bytes that can be seen
func Decrypt(sec *config.SharedEncryptionData,cfg *config.Config) ([]byte,error) {
	// Pre Safety 
	switch {
	case cfg == nil:
		return nil,&config.UserSafetyError{
			Cause: `Nil pointer reference`,
			Message: `A nil pointer to the configuration has been passed`,
		}
	case sec == nil :
		return nil,&config.UserSafetyError{
			Cause: `Nil pointer reference`,
			Message: `A nil pointer to the Shared Encryption data has been passed`,
		}
	}

	key := argon2.IDKey(
		[]byte(cfg.Password),
		sec.Salt,
		1,
		64*1024,
		4,
		32,
	)

	// - AES Gcm open 
	block,err := aes.NewCipher(key)
	if err != nil {
		return nil,&config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot create a new cipher`,
		}
	}

	gcm,err := cipher.NewGCM(block)
	if err != nil {
		return nil,&config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot create a new gcm`,
		}
	}

	plaindata,err := gcm.Open(nil,sec.Nonce,sec.EncryptedData,nil)
	if err != nil {
		return nil,&config.UserSafetyError{
			Cause: err.Error(),
			Message: `Cannot make decrypt the data because of the given reason`,
		}
	} 

	return plaindata,nil
}