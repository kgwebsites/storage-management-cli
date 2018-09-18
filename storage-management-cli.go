package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"code.cloudfoundry.org/bytefmt"
	"github.com/olekukonko/tablewriter"
	"github.com/skratchdot/open-golang/open"
)

type file struct {
	ID   int    `json:"ID"`
	Size int64  `json:"Size"`
	Path string `json:"Path"`
}

var files = map[int]file{}
var filesSync = struct {
	sync.RWMutex
	files map[int]file
}{files: files}
var sortedFiles []file
var selection []file

var wg sync.WaitGroup
var platform = runtime.GOOS
var root string
var ignore string
var resultCount int
var tempFile *os.File
var cliActive = true
var cliStatus = ""

func init() {
	rt, err := os.Getwd()
	checkErr(err)
	flag.StringVar(&root, "d", rt, "The directory to search for files.")
	flag.StringVar(&ignore, "i", "", "Comma seperated list of paths of path keywords to ignore")
	flag.IntVar(&resultCount, "c", 10, "The number of files to show")
}

func main() {
	flag.Parse()
	// Walk path concurrently recursively and retrieve file meta data
	getFiles()

	// Sort files by size
	sortFiles(files)

	// CLI Loop
	for cliActive {
		// Get the CLI selection
		selection := cliSelection()

		// Generate CLI
		generateCLI(selection)
	}
}

func getFiles() {
	wg.Add(1)
	walkDir(root)
	wg.Wait()
	files = filesSync.files
}

func walkDir(dir string) {
	defer wg.Done()
	i := 1
	visit := func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ignore) {
			return nil
		}
		filesSync.Lock()
		filesSync.files[i] = file{ID: i, Path: path, Size: f.Size()}
		filesSync.Unlock()
		i++
		if f.IsDir() && path != dir {
			wg.Add(1)
			go walkDir(path)
			return filepath.SkipDir
		}
		return nil
	}

	filepath.Walk(dir, visit)
}

func trimPath(path string) string {
	path = strings.Trim(path, "\r")
	sp := strings.Split(path, "/")
	rs := strings.Split(root, "/")
	rootBase := rs[len(rs)-1]
	if strings.Split(root, "")[len(root)-1] == "/" {
		rootBase = rs[len(rs)-2]
	}
	if len(sp) >= 2 {
		return strings.Join([]string{rootBase, "...", sp[len(sp)-1]}, "/")
	}
	return path
}

func sortFiles(files map[int]file) {
	sortedFiles = []file{}
	for f := range files {
		sortedFiles = append(sortedFiles, files[f])
	}
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Size > sortedFiles[j].Size
	})
}

func cliSelection() []file {
	if len(sortedFiles) > resultCount {
		return sortedFiles[:resultCount]
	}
	return sortedFiles
}

func generateCLI(selection []file) {
	// Configure CLI table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Bytes", "Path"})
	table.SetRowLine(true)

	// Build CLI table
	for _, v := range selection {
		// Convert byte size to human readable sizes
		byteSize := bytefmt.ByteSize(uint64(v.Size))
		table.Append([]string{fmt.Sprint(v.ID), fmt.Sprint(byteSize), v.Path})
	}

	// Print out CLI table
	table.Render() // Send output

	// Print CLI status
	if len(cliStatus) > 0 {
		fmt.Println(cliStatus)
	}

	// Request input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Commands: cd <PATH> | delete <ID> | open <ID> | path <ID> | more <NUMBER> | exit")
	t, _ := reader.ReadString('\n')
	te := strings.Split(t, "\n")
	text := te[0]
	entry := strings.Split(text, " ")
	command := entry[0]

	// Exit
	if command == "exit" {
		cliActive = false
	}

	if len(entry) >= 2 {
		param := entry[1]
		// If command "cd" use the param as a string and Change directories
		if command == "cd" {
			changeDirectory(param)
		} else {
			// Otherwise parse the param as an int and Set the selected file to that int
			id, err := strconv.Atoi(param)
			checkErr(err)
			selectedFile := files[id]

			// Delete file
			if command == "delete" {
				deleteFile(selectedFile.Path, selectedFile.ID)
			}

			// Open file
			if command == "open" {
				openFile(selectedFile.Path)
			}

			// Show full path
			if command == "path" {
				cliStatus = selectedFile.Path
			}

			// Show more results
			if command == "more" {
				res, err := strconv.Atoi(param)
				checkErr(err)
				resultCount = res
				cliSelection()
			}
		}
	} else {
		cliStatus = "Invalid Command"
	}
}

func changeDirectory(dir string) {
	root = dir
	files = map[int]file{}
	filesSync = struct {
		sync.RWMutex
		files map[int]file
	}{files: files}
	fmt.Println("Calculating all files in", dir)
	getFiles()
	sortFiles(files)
}

func deleteFile(filePath string, fileID int) {
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

	// Remove the file from and file list
	delete(files, fileID)

	// Resort file
	sortFiles(files)

	cliStatus = fmt.Sprintf("%v moved to the Trash", fileName)
}

func openFile(path string) {
	cliStatus = fmt.Sprintf("Opening %v", path)
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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
