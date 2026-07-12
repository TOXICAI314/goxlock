package dpapi

import (
	"goxlock/config"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Info : Point to point explanation in `dpapi/protect.go`

// - Dlls - //
var (
	// prcoUnprotectData -> Using the same dll from `protect.go` to uncover the data
	procUnprotectData = crypt32.NewProc("CryptUnprotectData")
)

// - Unprotect 
// Makes the data in raw byte format (readable string bytes)
func Unprotect(data []byte) ([]byte, error) {

	// - Pre safety - //
	if len(data) == 0 {
		return nil, nil
	}

	in := windows.DataBlob{
		Size: uint32(len(data)),
	}

	in.Data = &data[0]

	var out windows.DataBlob

	// - Decryting 
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

	defer procLocalFree.Call(uintptr(unsafe.Pointer(out.Data)))

	result := make([]byte, out.Size)

	copy(result, unsafe.Slice(out.Data, out.Size))

	return result, nil
}
