package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

// - GetSomeorNoneString
// Given string -> if Something returns else `None`
func GetSomeorNoneString(str string) string {
	switch str {
	case ``:
		return `None`
	default:
		return str
	}
}

// - SecretInput 
// As for password and for sesitive fields the `echo` of the terminal need to be either stopped or done in abrivated char
func SecretInput(prompt string) (string,error) {
	// Info : The code currently turns off the `echo` made by the terminal by each iteration of char
	// Just like Ubuntu password login or any thing else
	// Steps :
	// 1. File Descriptor of `os.Stdin` -> Yes its acts as a file and other two too
	// 2. Giving it to `term.ReadPassword` to make a secure byte password
	fmt.Print(prompt)

	fileDes := os.Stdin.Fd()
	bytePass,err := term.ReadPassword(int(fileDes))
	fmt.Print("\n")
	if err != nil {
		return ``,err
	}
	if string(bytePass) == `` {
		return ``,fmt.Errorf(`Empty password input is not allowed`)
	}
	return string(bytePass),nil
}

// - ConfirmSecurePassword 
// Make a 2 step verification for the password so that data wont be locked away by mistake
func ConfirmSecretPassword(prompt string) (string,error) {
	// - Pass input - //
	p1,err := SecretInput(prompt)
	if err != nil {
		return ``,err
	}

	// - Check input 
	p2,err := SecretInput(`Confirm the password: `)
	if err != nil {
		return ``,err
	}

	if p1 != p2 {
		return ``,fmt.Errorf(`Password integrity failed, Wrong password trial`)
	}
	return p1,nil
}

// - CreateSessionID 
// Is a random string generator which will genrate a session id that can be used in configuration
func CreateSessionID() string {
	randomSessionSeed := make([]byte,12)
	rand.Read(randomSessionSeed)
	return hex.EncodeToString(randomSessionSeed)
}

// - Input 
// A python like easy input statement given returning a string
func Input(prompt string) string {
	// Info : 
	// Steps goes like ->
	// 1. Print prompt
	// 2. Scanner to scan `os.Stdin`
	// 3. Returns the string from Scanner
	fmt.Print(prompt)
	var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

// - ClearAllTempFolderJunk 
// This clears the junk of the temp folder made by goxlock during proccess
func ClearAllTempFolderJunk(name string) error {
	temdir := os.TempDir()
	dirdata,err := os.ReadDir(temdir)
	if err != nil {
		return err
	}
	pattern := fmt.Sprintf(`%s-*`,name)
	for _,entries := range dirdata {
		name = entries.Name()
		// Info : Now filepath.Match will be used to match any junk folder
		match,err := filepath.Match(pattern,name)
		if err != nil {
			continue
		}
		if match {
			fullPath := filepath.Join(temdir, name)
			err := os.RemoveAll(fullPath)
			if err != nil {
				fmt.Println(`Cannot delete `,name,` in temdir`)
			}
		}
	}
	return nil
}

// - GetYesORNo
// Gets the input answer for yes or no from the user
func GetYesORNo(prompt string) bool {
	ans := Input(prompt)
	if strings.ToLower(strings.TrimSpace(ans)) == `y` {
		return true
	}
	return false
}

// - Where
// Tells the user where the running executable is located
func Where() (string,error) {
	return os.Executable()
} 