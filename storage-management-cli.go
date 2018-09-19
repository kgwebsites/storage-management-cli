package main

import (
	"flag"
	"log"
	"os"

	"github.com/kgwebsites/storagemanagementcli/modules"
)

var files = storagemanagementcli.FileMap{}
var sortedFiles = storagemanagementcli.Files{}

var root string
var ignore string
var resultCount int
var cliActive = true
var cliStatus = ""

func init() {
	rt, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&root, "d", rt, "The directory to search for files.")
	flag.StringVar(&ignore, "i", "", "Comma seperated list of paths of path keywords to ignore")
	flag.IntVar(&resultCount, "c", 10, "The number of files to show")
}

func main() {
	flag.Parse()
	// Walk path concurrently recursively and retrieve file meta data
	files = storagemanagementcli.GetFiles(root, ignore)

	// Sort files by size
	sortedFiles = storagemanagementcli.SortFiles(files)

	// CLI Loop
	for cliActive {
		// Get the CLI selection
		selection := storagemanagementcli.LimitFiles(sortedFiles, resultCount)
		// Configure CLI
		var config = storagemanagementcli.CLIConfig{
			selection,
			&files,
			&sortedFiles,
			&resultCount,
			&cliStatus,
			&cliActive,
			&ignore,
		}
		// Generate CLI
		storagemanagementcli.GenerateCLI(config)
	}
}
