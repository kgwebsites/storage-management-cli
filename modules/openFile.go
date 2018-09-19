package storagemanagementcli

import (
	"os"
	"path/filepath"

	"github.com/skratchdot/open-golang/open"
)

func openFile(path string) {
	// If path is a directory open it, otherwise open containing directory
	f, err := os.Stat(path)
	file := f.Mode()
	checkErr(err)
	if file.IsDir() {
		open.Run(path)
	} else {
		open.Run(filepath.Dir(path))
	}
}
