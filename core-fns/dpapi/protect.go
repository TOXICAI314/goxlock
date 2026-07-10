package dpapi

import (
	"goxlock/config"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// - Dlls
// Dlls will be used to call the windows encryption agents -> Encryption is unique per user
var (
	// crypt32 -> A cryptography dll for the windows system
	crypt32  = windows.NewLazyDLL(`Crypt32.dll`)
	// kernel32 -> The main which provokes the handle sharing and closing of them
	kernel32 = windows.NewLazyDLL(`Kernel32.dll`)

	// - Process Calls 
	// Precess for dll -> Internal function that are usable

	// procPrtotectData -> The function for the encryption of the dll
	procProtectData = crypt32.NewProc(`CryptProtectData`)
	// procLocalFree -> Frees the space made by those native `C` functions
	procLocalFree   = kernel32.NewProc(`LocalFree`)
)

// - Flags 
// These flags are constant and are used for custom behaviours in the dll
const (
	// CryptoProtectUIForbidden -> Prevnts the ui for the cryptography from opening
	CRYPTPROTECT_UI_FORBIDDEN = 0x1
)

// - Protect 
// Protect : Will encrypt the data that is given to it by the os level encryption
func Protect(data []byte) ([]byte, error) {
	// - Pre Safety Check 
	if data == nil {
		return nil, &config.FunctionCancelError{
			Cause:   `Data is nil`,
			Message: `Cannot encypt the data because of its length`,
			ElapsedTime: time.Now(),
			Provider: `dpapi.Protect`,
		}
	}

	// - Data blob
	// The needed fundamental for the sharing of the data b/w windows api
	in := windows.DataBlob{
		Size: uint32(len(data)),
	}

	// Info : As in `C` the array degrades to a pointer
	// The must be a reference to that starting pointer
	// in.Size will handle the rest of the size matter
	in.Data = &data[0]

	var out windows.DataBlob

	// - Function 
	// `Call` makes that function get called with the desired parametres
	// See the web for the parametres or just breifly see them here
	r1, _, err := procProtectData.Call(
		// Info : Every data here has to be a uintptr -> a pointer
		// As these heavenly depends on the pointer data passage

		// pDataIn : Where the Datablob will be taken
		uintptr(unsafe.Pointer(&in)),
		// sZDataDescr
		0,
		// pOptionalEntropy
		0,
		// pvReserved
		0,
		// pPromptStruct
		0,
		// dwFlags : Custom behaviour flags
		CRYPTPROTECT_UI_FORBIDDEN,
		// pDataOut : The pointer where the data will be given out
		uintptr(unsafe.Pointer(&out)),
	)

	// Info : Any error tackled during the program || The given out is zero
	// As err != nil always even if its not an error 
	// As error is an interface it depends on two things : {
	// 		1. The underlying value -> Thats why an array can be a nil pointer
	//		2. Concrete data type	-> Thats why a pointer can be nil itelself
	//}
	// But here even if the work is done it always returns a value
	// Low level function always run `GetLastError` immediatly to catch their errors
	// That error is then wrapped in a error struct -> `syscall.Errno(<num>)` 
	// That num can be 0 -> success(true) , 1 -> Failure(false)
	// So even if the program is successfuly completed it will always get an error as `syscall.Errno(0)`
	// Which have its own value making : err != nil
	// Therefore for low level work : always use r1 == 0
	if r1 == 0 {
		return nil, &config.FunctionFailError{
			Cause:   err.Error(),
			Message: `Encryption Failed due to local dll failures`,
			ElapsedTime: time.Now(),
			Provider: `dpapi.Portect`,
		}
	}

	// Info : As data is from `dll` that is made on `C` and `C++`
	// The memory has to be freed from our side (via a strict order) by another dll
	// And as out.Data already contains the `datablob`its size is known by the compiler
	defer procLocalFree.Call(uintptr(unsafe.Pointer(out.Data)))

	// - Extraction 
	result := make([]byte, out.Size)
	// Info : As out data is not of this programm
	// The programm cant use it directly -> alternative .. copy
	copy(result, unsafe.Slice(out.Data, out.Size))

	return result, nil
}
