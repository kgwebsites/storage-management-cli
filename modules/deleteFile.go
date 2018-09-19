package storagemanagementcli

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

func deleteFile(filePath string, fileID int, Files *FileMap, SortedFiles *Files) {
	paths := strings.Split(filePath, "/")
	fileName := paths[len(paths)-1]

	// Get user meta data
	usr, err := user.Current()
	checkErr(err)

	// Place the file in the trash
	// Mac OS
	if platform == "darwin" {
		err := os.Rename(filePath, fmt.Sprintf("%v/.Trash/%v", usr.HomeDir, fileName))
		checkErr(err)
	}
	// Linux OS
	if platform == "linux" {
		err := os.Rename(filePath, fmt.Sprintf("%v/.local/share/Trash/%v", usr.HomeDir, fileName))
		checkErr(err)
	}
	// Windows OS - No access to recycle bin so remove file
	if platform == "windows" {
		err := os.Remove(filePath)
		checkErr(err)
	}

	// Remove the file from the file list
	delete(*Files, fileID)

	// Resort file
	*SortedFiles = SortFiles(*Files)
}
