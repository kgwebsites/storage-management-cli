package storagemanagementcli

import (
	"sync"
)

var wg sync.WaitGroup

var filesSync = struct {
	sync.RWMutex
	files FileMap
}{files: FileMap{}}

// GetFiles returns a map of the files in a path
func GetFiles(
	root string,
	ignore string) FileMap {

	wg.Add(1)
	walkDir(root, ignore)
	wg.Wait()
	return filesSync.files
}
