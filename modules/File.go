package storagemanagementcli

// File contains a unique ID, a byte size, and the path where it came from
type File struct {
	ID   int
	Size int64
	Path string
}

// Files contains a slice of File
type Files []File

// FileMap is a map of indexes to files
type FileMap map[int]File
