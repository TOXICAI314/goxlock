package mutex

import "golang.org/x/sys/windows"

// - CloseMutex 
// Will close the given folder mutex -> Normally ending
// Even if crashed app , os will free automatically
func (fMutex *FolderMutex) CloseMutex() {
	if fMutex.handle == 0 {
		// Info : If the handle is already 0 (closed) -> return nothing
		return 
	}
	// else juts close it and assign it as 0 (freeing the handle)
	// In `C` thats `NULL` -> `nil`
	windows.CloseHandle(fMutex.handle)
	fMutex.handle = 0
}