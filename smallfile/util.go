package smallfile

import "unsafe"

func GlobalPageIDToBytes(id GlobalPageID) *recordUint64 {
	return (*recordUint64)(unsafe.Pointer(&id))
}
