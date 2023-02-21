package main

import (
	"fmt"
	"os"
)

var longListing = false

func main() {
	offset := 1
	if 2 <= len(os.Args) && "-1" == os.Args[offset] {
		longListing = true
		offset++
	}

	files := os.Args[offset:]
	if 0 == len(files) {
		files = []string{"."}
		// show current directory
	}

	if longListing {
		showLongListing(os.Args[offset:])
		return
	}
	showShortListing(os.Args[offset:])
}

func showShortListing(files []string) {
	var noFilesList []string
	var filesList []string
	var dirListing []string
	fmt.Println("showing short listing for: ", files)
	for _, f := range files {
		fi, err := os.Stat(f)
		if nil != err {
			s := fmt.Sprintf("ls: %v: no file or directory\n", f)
			noFilesList = append(noFilesList, s)
			continue
		}
		if !fi.IsDir() {
			filesList = append(filesList, f)
			continue
			// get files in directory
		}
		dirListing = addDirListing(dirListing, f)
	}
	for _, s := range noFilesList {
		fmt.Println(s)
	}
	for _, s := range filesList {
		fmt.Println(s)
	}
	fmt.Println("")
	for _, s := range dirListing {
		fmt.Println(s)
	}
}

func addDirListing(listing []string, f string) []string {
	dir, err := os.Open(f)
	if nil != err {
		return listing
	}
	fileNames, err := dir.Readdirnames(0)
	if nil != err {
		return listing
	}
	listing = append(listing, "\n"+f+":")
	for _, d := range fileNames {
		listing = append(listing, d)
	}
	return listing
}
func showLongListing(files []string) {
	fmt.Println("showing long listing for: ", files)
	var noFilesList []string
	var filesList []string
	var dirListing []string
	for _, f := range files {
		fi, err := os.Stat(f)
		if nil != err {
			s := fmt.Sprintf("ls: %v: no file or directiory", f)
			noFilesList = append(filesList, s)
			continue
		}

		if !fi.IsDir() {
			size := getSize(fi.Size())
			perm := getPermString(fi.Mode())
			s := fmt.Sprintf("-rw-r--r-- | 0B May 18 09:23 %s", f)
			filesList = append(filesList, s)
			continue
		}
		dirListing = addDirListing(dirListing, f)
	}
}
