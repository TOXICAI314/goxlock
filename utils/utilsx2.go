package utils

import (
	"fmt"
	"regexp"
	"runtime"
)

// - Const
const (
	// DeleteUpperCursorData -> The ANSII keyword that deletes the upper written line from the buffer of the writter
	DeleteUpperCursorData string = "\033[A\033[K\033[A\033[K"
)

// GetBasicVersionDetails
// Gives the version from the system along with golang version and application version
func GetBasicVersionDetails(appVersion string) string {
	return fmt.Sprintf(
		`
		goxlock : %s
		Golang version : %s
		OS : %s
		Architecture : %s
		`,appVersion,runtime.Version(),runtime.GOOS,runtime.GOARCH,
	)
}

// ValidateCharacters
// Validates the caharcter set to a limited number of char like Aa-Zz,0-9,_,-
func ValidateCharacters(name string) (isValid bool) {
	var validProfile = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	return validProfile.MatchString(name) 
}