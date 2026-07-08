package config

import "fmt"

// - EncryptionError 
// EncryptionError : The error which arose due to undesirables in the configuration while encrypting
type EncryptionError struct {
	Cause   string
	Message string
	Fix     string
}

func (ec *EncryptionError) Error() string {
	return fmt.Sprintf("Cause : %s\nMessage : %s\nFix : %s", ec.Cause, ec.Message, ec.Fix)
}

// - ZipError S 
// ZipError : arises due to the formation of broken commitments in zipping process
type ZipError struct {
	Cause   string
	Message string
	Fix     string
}

func (uzp *ZipError) Error() string {
	return fmt.Sprintf("Cause : %s\nMessage : %s\nFix : %s", uzp.Cause, uzp.Message, uzp.Fix)
}

// - UnzipError 
// UnzipError : arises due to the formation of broken commitments in unzipping process
type UnzipError struct {
	Cause   string
	Message string
	Fix     string
}

func (uzp *UnzipError) Error() string {
	return fmt.Sprintf("Cause : %s\nMessage : %s\nFix : %s", uzp.Cause, uzp.Message, uzp.Fix)
}


// - UserSafety 
type UserSafetyError struct {
	Cause   string
	Message string
}

// Info : Making sure that this lies in the error interface
func (usersafe *UserSafetyError) Error() string {
	return fmt.Sprintf("Cause : %s\nMessage : %s", usersafe.Cause, usersafe.Message)
}

// - DecryptionError
// DecryptionError : The error which arose due to undesirables in the configuration while decrypting
type DecryptionError struct {
	Cause   string
	Message string
	Fix     string
}

func (dc *DecryptionError) Error() string {
	return fmt.Sprintf("Cause : %s\nMessage : %s\nFix : %s", dc.Cause, dc.Message, dc.Fix)
}

// - ChnagePasswordError - //
type ChnagePasswordError struct {
	Cause 	string
	Message	string
	Fix		string
}

func (cpe *ChnagePasswordError) Error() string {
	return fmt.Sprintf("Cause : %s\nMessage : %s\nFix : %s", cpe.Cause, cpe.Message, cpe.Fix)
}

// - CheckError 
type CheckError struct {
	Cause 	string
	Message	string
	Fix		string
}

func (e *CheckError) Error() string {
	return fmt.Sprintf("Cause : %s\nMessage : %s\nFix : %s", e.Cause, e.Message, e.Fix)
}