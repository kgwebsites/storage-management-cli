package storagemanagementcli

var fileMap = FileMap{}

// GetFiles returns a map of the files in a path
var GetFiles = func(
	root string,
	ignore string) FileMap {
	walkDir(root, ignore)
	return fileMap
}
