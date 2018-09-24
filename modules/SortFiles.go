package storagemanagementcli

import "sort"

// SortFiles sorts a file map into a slice of files ordered by byte size
var SortFiles = func(files FileMap) Files {
	sortedFiles := Files{}
	for f := range files {
		sortedFiles = append(sortedFiles, files[f])
	}
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Size > sortedFiles[j].Size
	})
	return sortedFiles
}
