package storagemanagementcli

import (
	"os"
	"path/filepath"
	"strings"
)

func walkDir(dir string, ignore string) {
	defer wg.Done()
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
		filesSync.Lock()
		filesSync.files[i] = File{ID: i, Path: path, Size: f.Size()}
		filesSync.Unlock()
		i++
		if f.IsDir() && path != dir {
			wg.Add(1)
			go walkDir(path, ignore)
			return filepath.SkipDir
		}
		return nil
	}

	filepath.Walk(dir, visit)
}
