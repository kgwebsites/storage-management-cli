package storagemanagementcli

func changeDirectory(dir string, SortedFiles *Files, Ignore *string) {
	Files := GetFiles(dir, *Ignore)
	*SortedFiles = SortFiles(Files)
}
