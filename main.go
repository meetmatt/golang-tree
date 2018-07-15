package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

/*
dir			=>	prefix = prefix+(isLastDir ? "	└" : "	├")
hasFiles		=>	prefix = prefix+(isLastDir ? "	└" : "	├")
*/

func printer(out io.Writer, fileInfo os.FileInfo, printFiles bool, depth int, currentPath string, prefix string, isLast bool) {
	if fileInfo.IsDir() {
		// print directory name
		fmt.Fprintln(out, prefix+"───"+fileInfo.Name())
		// go through children
		files, error := ioutil.ReadDir(currentPath + "/" + fileInfo.Name())
		if error != nil {
			fmt.Fprintln(out, error)
		}
		totalChildren := len(files)
		newPrefix := prefix
		for index, file := range files {
			printer(out, file, printFiles, depth+1, currentPath+"/"+fileInfo.Name(), newPrefix, index == totalChildren)
			if file.IsDir() {
				if index == totalChildren-1 {
					newPrefix = prefix + "	└"
				} else {
					newPrefix = prefix + "	|"
				}
			} else {
				if index == totalChildren-1 {
					newPrefix = prefix + "	└"
				} else {
					newPrefix = prefix + "	|"
				}
			}
		}
		prefix += "	"
	} else {
		if printFiles {
			fmt.Fprintln(out, prefix+"───"+fileInfo.Name())
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	files, error := ioutil.ReadDir(path)
	if error != nil {
		return fmt.Errorf("Error reading directory", error)
	}
	for _, file := range files {
		printer(out, file, printFiles, 0, path, "|", false)
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
