package main

import (
	"goxlock/cmd"
)

/*
- `goxlock` is a tool used for a safe encryption on the targetted model of windows -> 11
- As windows 11 dont have a native safe folder option : This app will create that itself
- It will encrypt your folder so that no one can see them , and someone has tempered them then you will know
- If you have any issues regarding the uses of this app -> Report to me : `TOXICAI314` is my user on github
*/

func main() {
	// Execution Starts 
	// Info : This execution includes the calling of root command done in `cmd/root.go`
	// Any Error passed through it will be counted here and will act as Fatal error
	cmd.Execute()
}