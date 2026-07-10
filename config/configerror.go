package config

import (
	"fmt"
	"time"
)

// - FunctionFailError
// Function Fail error will be shown when a function has failed due to its internal structure or underlying function
type FunctionFailError struct {
	Cause       string
	Message     string
	ElapsedTime time.Time
	Provider    string
}

func (ff *FunctionFailError) Error() string {
	return fmt.Sprintf("Cause (The function Failed because): %s\nMessage : %s\nCrashed Time : %s\nCrash Function : %s", ff.Cause, ff.Message, ff.ElapsedTime.Format("15:04:05"), ff.Provider)
}

// - FunctionCancelError
// The error which is formed when the function returns abnormally due to wrong configuration
type FunctionCancelError struct {
	Cause 		string
	Message 	string
	ElapsedTime	time.Time
	Provider 	string
}

func (fc *FunctionCancelError) Error() string {
	return fmt.Sprintf("Cause (The function abnormally ended): %s\nMessage : %s\nCrashed Time : %s\nCrash Function : %s",fc.Cause,fc.Message,fc.ElapsedTime.Format(`15:04:05`),fc.Provider)
}