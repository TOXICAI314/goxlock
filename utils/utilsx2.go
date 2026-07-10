package utils

import (
	"fmt"
	"runtime"
)

// - GetBasicVersionDetails
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