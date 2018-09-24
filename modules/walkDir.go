package storagemanagementcli

import (
	"os"
	"path/filepath"
	"strings"
)

func walkDir(dir string, ignore string) {
	i := 1
	visit := func(path string, f os.FileInfo, err error) error {
		if len(ignore) > 0 {
			ignoreList := strings.Split(ignore, ",")
			for _, i := range ignoreList {
				if strings.Contains(path, strings.TrimSpace(i)) {
					return nil
				}
			}
		}
		fileMap[i] = File{ID: i, Path: path, Size: f.Size()}
		i++
		return nil
	}

	filepath.Walk(dir, visit)
}
