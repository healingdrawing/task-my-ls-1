package listing

/*
import (
	"fmt"
	"io/fs"
	"log"
	"my-ls-1/file"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// global variables to indicate the command line flags
var longListing = false
var showHidden = false
var recursive = false
var reverse = false
var sortByTime = false

func ShowShortListing(files []fs.DirEntry) {
	var noFilesList []string // list of files that could not be listed
	var dirListing []string  // list of directories to be listed
	for _, d := range files {
		// Check if the file or directory exists.
		_, err := d.Info()
		if err != nil {
			s := fmt.Sprintf("ls: %v: no file or directory", d.Name()) // if file not found, add error message to noFilesList
			noFilesList = append(noFilesList, s)
			continue
		}
		dirListing = AddShortDirListing(dirListing, d.Name(), false)
	}
	// Print out the lists.
	for _, s := range noFilesList {
		fmt.Println(s)
	}
	for _, s := range dirListing {
		fmt.Println(s)
	}
}

/*
// Show short listing function.
func ShowShortListing(files []string, recursive bool) {
	var noFilesList []string // list of files that could not be listed
	var filesList []string   // list of files to be listed
	var dirListing []string  // list of directories to be listed
	for _, f := range files {
		// Check if the file or directory exists.
		fi, err := os.Stat(f) // get file info
		if nil != err {
			s := fmt.Sprintf("ls: %v: no file or directory", f) // if file not found, add error message to noFilesList
			noFilesList = append(noFilesList, s)
			continue
		}
		// if file is hidden and showHidden flag is not set, continue to next file
		//if !showHidden && strings.HasPrefix(fi.Name(), ".") {
		if !showHidden && fi.Name()[0] == '.' {
			continue
		}
		if !fi.IsDir() { // if the file is not a directory, add it to filesList
			filesList = append(filesList, f)
			continue
		}

		// If it is a directory, get a short listing of its contents.
		//dirListing, err = AddShortDirListing(dirListing, f)
		//if err != nil {
		//	s := fmt.Sprintf("ls: %v: %v", f, err)
		//	noFilesList = append(noFilesList, s)
		dirListing = AddShortDirListing(dirListing, f, recursive)

	}
	// Print out the lists.
	for _, s := range noFilesList {
		fmt.Println(s)
	}
	for _, s := range filesList {
		fmt.Println(s)
	}
	for _, s := range dirListing {
		fmt.Println(s)
	}

}*/

/*
func AddShortDirListing(listing []string, f string, recursive bool) []string {
	dir, err := os.Open(f)
	if nil != err {
		return listing
	}
	fileInfos, err := dir.Readdir(0)
	if nil != err {
		return listing
	}
	listing = append(listing, "\n"+f+":")
	for _, fi := range fileInfos {
		if !showHidden && fi.Name()[0] == '.' {
			continue
		}
		if !recursive && fi.IsDir() {
			continue
		}
		listing = append(listing, fi.Name())
	}
	return listing
}

/*
// addShortDirListing takes a slice of strings representing the directory listing and a directory name and returns an updated slice with the short listing of the directory's contents.
// It reads the directory contents using os.Open and appends each file and directory name to the slice, with directories ending in a slash.
// If showHidden is false, files starting with a dot are skipped.
func AddShortDirListing(dirListing []string, dirName string) ([]string, error) {
	f, err := os.Open(dirName) // open the directory
	if err != nil {

		return nil, err // if there's an error opening the directory, return nil and the error
	}

	defer f.Close() // close the file when the function returns

	fis, err := f.Readdir(-1) // read the directory contents
	if err != nil {

		return nil, err // if there's an error reading the directory, return nil and the error
	}
	// for each file/directory in the directory
	for _, fi := range fis {
		// if showHidden is false and the file/directory starts with a dot, skip it
		if !showHidden && strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		s := fi.Name()  // get the file/directory name
		if fi.IsDir() { // if it's a directory, add a slash to the end
			s += "/"
		}
		dirListing = append(dirListing, s) // add the file/directory name to the slice
	}
	return dirListing, nil // return the updated slice and no error
}*/
/*
// addDirListing takes a slice of strings representing the directory listing, a file path, and a boolean indicating whether to use long listing format.
// It returns an updated slice with the directory listing.
// It opens the directory using os.Open and reads its contents with Readdir.
// It then appends each file or directory name to the slice, either in short or long listing format.
// If showHidden is false, files starting with a dot are skipped.
// If sortByTime is true, the directory contents are sorted by modification time.
func AddDirListing(listing []string, f string, longListing bool) []string {
	dir, err := os.Open(f) // open the directory
	if err != nil {
		log.Printf("Error opening directory: %v", err)
		return listing // if there's an error opening the directory, log the error and return the original slice
	}
	defer dir.Close() // close the directory when the function returns

	fis, err := dir.Readdir(0) // read the directory contents
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		// if there's an error reading the directory, log the error and return the original slice
		return listing
	}

	if sortByTime { // if we need to sort by modification time
		sort.Slice(fis, func(i, j int) bool {
			// sort the files by modification time in descending order
			return fis[i].ModTime().After(fis[j].ModTime())
		})
	}

	listing = append(listing, "\n"+f+":") // add a new line and the directory path to the slice

	for _, fi := range fis { // for each file/directory in the directory
		// if showHidden is false and the file/directory starts with a dot, skip it
		if !showHidden && fi.Name()[0] == '.' {
			continue
		}

		if longListing { // if we're using long listing format
			var sb strings.Builder // Create a new strings builder to store the long listing format
			// Gets the file mode and add it to the builder
			sb.WriteString(file.GetFileMode(fi.Mode()))
			// Gets the file permissions and add it to the builder
			sb.WriteString(file.GetFilePermissions(fi.Mode()))
			// Get the user ID and add it to the builder
			sb.WriteString(fmt.Sprintf(" %d ", fi.Sys().(*syscall.Stat_t).Uid))
			// Get the file size and add it to the builder
			sb.WriteString(fmt.Sprintf("%10d ", fi.Size()))
			// Get the modification time and add it to the builder
			sb.WriteString(getTimeString(fi.ModTime()))
			// Get the file name and add it to the builder
			sb.WriteString(fmt.Sprintf(" %s", fi.Name()))
			// Append the long listing format to the listing slice
			listing = append(listing, sb.String())
		} else {
			listing = append(listing, fi.Name()) // Append the file name to the listing slice
		}
	}
	return listing // Return the listing slice
}

// This function takes a directory path as input and prints the directory path itself
// and the long listing of all the files and directories inside it recursively
func ShowLongListingRecursive(dir string) {
	// Print the current directory path
	fmt.Println(dir)
	// Use filepath.Walk function to recursively iterate over all files and directories in the specified directory and its subdirectories
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil { // Handle any errors that occur during filepath.Walk
			return err
		}
		// Check if the file is a hidden file and if the showHidden flag is false, then skip it
		if !showHidden && strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		// If the file is a directory, print the directory path
		if info.IsDir() {
			fmt.Printf("%s:\n", path)
		} else {
			// If the file is a regular file, print its information
			fmt.Printf("%s %d %s %s\n", info.Mode(), info.Size(), info.ModTime().Format("Jan _2 15:04"), info.Name())
		}
		return nil
	})
	// Handle any errors that occur during filepath.Walk
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return
}

func addLongDirListing(listing []string, d fs.DirEntry) []string {
	if !showHidden && strings.HasPrefix(d.Name(), ".") {
		return listing
	}
	info, err := d.Info()
	if err != nil {
		return listing
	}
	perm := info.Mode().String()
	size := fmt.Sprintf("%d", info.Size())
	modTime := info.ModTime().Format("Jan _2 15:04")
	listing = append(listing, fmt.Sprintf("%s %s %s %s", perm, size, modTime, d.Name()))
	return listing
}

// This function takes a list of file paths as input and prints
// the long listing of each file (mode, size, modification time, and path)
func ShowLongListing(files []fs.DirEntry) {
	for _, file := range files {
		// Get the file information for the current file
		fileInfo, err := file.Info()
		if err != nil {
			// Handle any errors that occur while reading the file info
			fmt.Fprintf(os.Stderr, "Error reading file info: %v\n", err)
			continue
		}
		// Print the long listing for the current file (mode, size, modification time, and path)
		fmt.Printf("%s\t%d\t%s\t%s\n",
			fileInfo.Mode(),
			fileInfo.Size(),
			fileInfo.ModTime().Format("Jan 02 15:04"),
			file)
	}
}

func getLongListing(path string, info os.FileInfo) string {
    mode := info.Mode().String()
    owner, group, err := GetOwnerGroup(info)
    if err != nil {
        owner = strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Uid))
        group = strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Gid))
    }
    size := info.Size()
    modTime := info.ModTime().Format("Jan _2 15:04")

    return fmt.Sprintf("%s %s %s %6d %s %s", mode, owner, group, size, modTime, path)
}

/*
func getLongListing(path string, info fs.FileInfo) string {
	mode := info.Mode().String()
	fileInfo, err := os.Stat(path)
	if err != nil {
		return ""
	}
	owner, group, err := file.GetOwnerGroup(fileInfo)
	if err != nil {
		owner = strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Uid))
		group = strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Gid))
	}
	size := info.Size()
	modTime := info.ModTime().Format("Jan _2 15:04")

	return fmt.Sprintf("%s %s %s %6d %s %s", mode, owner, group, size, modTime, path)
}*/
/*
// This function takes a time.Time object as input and
// returns a formatted string representing the time in the format "2006-01-02 15:04:05"
func getTimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Function formatTime formats the time based on whether the duration between the input time and the current time is more than a year or not.
// If it is more than a year, it formats the time in "Jan 2 2006" format, else in
func formatTime(t time.Time) string {
	// define a default date/time layout format
	layout := "Jan 2 15:04"
	// if the duration between the current time and the specified time is more than a year,
	// change the layout to include the year
	if time.Now().Sub(t).Hours() > 24*365 {
		layout = "Jan 2 2006"
	}
	// format the specified time using the layout and return the formatted time string
	return t.Format(layout)
}
*/
