package doctor

import "goxlock/config"

// - Const 
const (
	// secondarypasswordname -> Is a substitute for the main password which is made while checking up for all the actions
	secondarypasswordname string = `goxlock2`
)

// - Doctor 
// The struct that conatains the data to run all the tests
type Doctor struct {
	config  	config.Config
}