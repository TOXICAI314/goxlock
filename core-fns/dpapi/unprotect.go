package dpapi

import (
	"goxlock/config"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Info : Point to point explanation in `dpapi/protect.go`

var (
	// prcoUnprotectData -> Using the same dll from `protect.go` to uncover the data
	procUnprotectData = crypt32.NewProc("CryptUnprotectData")
)

// Makes the data in raw byte format (readable string bytes)
func Unprotect(data []byte) (retData []byte, err error) {

	// Pre safety 
	if len(data) == 0 {
		return nil, nil
	}

	in := windows.DataBlob{
		Size: uint32(len(data)),
	}

	in.Data = &data[0]

	var out windows.DataBlob

	// Decryting 
	// Using the same format as done in protect
	r1, _, err := procUnprotectData.Call(
		uintptr(unsafe.Pointer(&in)),
		0,
		0,
		0,
		0,
		CRYPTPROTECT_UI_FORBIDDEN,
		uintptr(unsafe.Pointer(&out)),
	)

	if r1 == 0 {
		return nil, &config.FunctionFailError{
			Cause:   err.Error(),
			Message: `Decryption failed due to internal dll erros`,
			Provider: `dpapi.Unprotect`,
			ElapsedTime: time.Now(),
		}
	}

	defer func() {
		r1,_,freeError := procLocalFree.Call(uintptr(unsafe.Pointer(out.Data)))
		if r1 != 0 && err == nil {
			err = &config.FunctionFailError{
				Cause: freeError.Error(),
				Message: `The pointer to the out Data cant be cleared by the application`,
				ElapsedTime: time.Now(),
				Provider: `dpapi.Unprotect`,
			}
		}
	}()

	result := make([]byte, out.Size)

	copy(result, unsafe.Slice(out.Data, out.Size))

	return result, nil
}
