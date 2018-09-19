package storagemanagementcli

// LimitFiles returns the files passed in, restricting the amount to the [limit] passed in.
func LimitFiles(files Files, limit int) Files {
	if len(files) > limit {
		return files[:limit]
	}
	return files
}
