package main

import (
	"io"
	"os"
	"io/ioutil"
	"fmt"
	"strconv"
)

func countChildren(printFiles bool, files []os.FileInfo) int {
	// count children
	var totalChildren int
	if printFiles {
		// count all children
		totalChildren = len(files)
	} else {
		// count only directories
		for _, file := range files {
			if file.IsDir() {
				totalChildren += 1
			}
		}
	}
	return totalChildren
}

func printer(out io.Writer, fileInfo os.FileInfo, printFiles bool, path string, prefix string, isLast bool) {
	if fileInfo.IsDir() || printFiles {

		size := ""
		if printFiles && !fileInfo.IsDir() {
			if fileInfo.Size() > 0 {
				size = " ("+strconv.FormatInt(fileInfo.Size(), 10)+"b)"
			} else {
				size = " (empty)"
			}
		}

		if isLast {
			fmt.Fprintln(out, prefix+"└───"+fileInfo.Name()+size)
		} else {
			fmt.Fprintln(out, prefix+"├───"+fileInfo.Name()+size)
		}
	}

	if fileInfo.IsDir() {
		// increase depth
		if isLast {
			prefix += "	"
		} else {
			prefix += "│	"
		}

		// get children
		path += "/" + fileInfo.Name()
		files, _ := ioutil.ReadDir(path)
		totalChildren := countChildren(printFiles, files)
		printedFiles := 0
		isFileLast := false
		// recursively print children
		for _, file := range files {
			if !file.IsDir() && !printFiles {
				// don't even call printer if it's a file
				// and printFiles is false
				continue
			}
			printedFiles += 1
			isFileLast = printedFiles == totalChildren
			printer(out, file, printFiles, path, prefix, isFileLast)
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	files, _ := ioutil.ReadDir(path)
	totalChildren := countChildren(printFiles, files)
	printedFiles := 0
	isFileLast := false
	for _, file := range files {
		if !file.IsDir() && !printFiles {
			// don't even call printer if it's a file
			// and printFiles is false
			continue
		}
		printedFiles += 1
		isFileLast = printedFiles == totalChildren
		printer(out, file, printFiles, path, "", isFileLast)
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
