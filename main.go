package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
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

type Group struct {
	Gid  string // group ID
	Name string // group name
}

// Define map of file modes to their corresponding string representation.
var fileModeMap = map[os.FileMode]string{
	os.ModeDir:        "d",
	os.ModeSymlink:    "l",
	os.ModeNamedPipe:  "p",
	os.ModeSocket:     "s",
	os.ModeSetuid:     "s",
	os.ModeSetgid:     "s",
	os.ModeCharDevice: "c",
}

// Define map of file permissions to their corresponding string representation.
var filePermMap = map[os.FileMode]string{
	0400: "r",
	0200: "w",
	0100: "x",
	0040: "r",
	0020: "w",
	0010: "x",
	0004: "r",
	0002: "w",
	0001: "x",
}

func main() {
	// Set offset to 1 to skip program name.
	offset := 1

	// Check for command line flags.
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-l":
			longListing = true // set longListing flag to true
		case "-a":
			showHidden = true // set showHidden flag to true
		case "-R":
			recursive = true // set recursive flag to true
		case "-r":
			reverse = true // set reverse flag to true
		case "-t":
			sortByTime = true // set sortByTime flag to true
		default:
			// This is a file or directory argument, so break out of the loop.
			offset = i
			break
		}
	}

	// If there are no files specified, show the current directory.
	files := os.Args[offset:]
	if len(files) == 0 {
		files = []string{"."} // if no files are specified, set files to current directory
	}

	// If the -R flag is set, show long listing recursively.
	if recursive {
		for _, f := range files {
			showLongListingRecursive(f) // call function to show long listing recursively
		}
	} else {
		// Otherwise, show short or long listing depending on the -l flag.
		if longListing {
			showLongListing(files) // call function to show long listing
		} else {
			showShortListing(files) // call function to show short listing
		}
	}
}

// Show short listing function.
func showShortListing(files []string) {
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
		if !showHidden && strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		if !fi.IsDir() { // if the file is not a directory, add it to filesList
			filesList = append(filesList, f)
			continue
		}

		// If it is a directory, get a short listing of its contents.
		dirListing, err = addShortDirListing(dirListing, f)
		if err != nil {
			s := fmt.Sprintf("ls: %v: %v", f, err)
			noFilesList = append(noFilesList, s)
		}
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

}

// addShortDirListing takes a slice of strings representing the directory listing and a directory name and returns an updated slice with the short listing of the directory's contents.
// It reads the directory contents using os.Open and appends each file and directory name to the slice, with directories ending in a slash.
// If showHidden is false, files starting with a dot are skipped.
func addShortDirListing(dirListing []string, dirName string) ([]string, error) {
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
}

// addDirListing takes a slice of strings representing the directory listing, a file path, and a boolean indicating whether to use long listing format.
// It returns an updated slice with the directory listing.
// It opens the directory using os.Open and reads its contents with Readdir.
// It then appends each file or directory name to the slice, either in short or long listing format.
// If showHidden is false, files starting with a dot are skipped.
// If sortByTime is true, the directory contents are sorted by modification time.
func addDirListing(listing []string, f string, longListing bool) []string {
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
			sb.WriteString(getFileMode(fi.Mode()))
			// Gets the file permissions and add it to the builder
			sb.WriteString(getFilePermissions(fi.Mode()))
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
func showLongListingRecursive(dir string) {
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
		}
		// Call the showLongListing function to print the long listing of the file
		showLongListing([]string{path})
		return nil
	})
	// Handle any errors that occur during filepath.Walk
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

// This function takes a file size in bytes as input and returns a string representing
// the file size in a human-readable format (e.g. "10 KB", "2.5 MB", etc.)
func getFileSize(size int64) string {
	const unit = 1024
	// If the file size is less than 1 KB, return the size in bytes
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	// Otherwise, divide the file size by 1024 and calculate the appropriate unit (KB, MB, etc.) to use
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

// This function takes a list of file paths as input and prints
// the long listing of each file (mode, size, modification time, and path)
func showLongListing(files []string) {
	for _, file := range files {
		// Get the file information for the current file
		fileInfo, err := os.Stat(file)
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

// This function takes a system-specific type representing a file and returns the name of the file's owner and group as strings.
func getOwnerAndGroup(sys interface{}) (string, string) {
	// cast the input interface to a pointer to a syscall.Stat_t struct
	stat := sys.(*syscall.Stat_t)
	// get the user ID and group ID as strings
	uid := fmt.Sprint(stat.Uid)
	gid := fmt.Sprint(stat.Gid)
	// lookup the user information for the given user ID
	u, err := user.LookupId(uid)
	// if an error occurred, return the user ID and group ID as strings
	if err != nil {
		return uid, gid
	} // lookup the group information for the given group ID
	g, err := user.LookupGroupId(gid)
	// if an error occurred, return the username and group ID as strings
	if err != nil {
		return u.Username, gid
	}
	// return the username and group name as strings
	return u.Username, g.Name
}

// LookupGroupId is a wrapper around the user.LookupGroupId function.
// It takes a group ID string as an argument and returns a user.Group object and an error.
func LookupGroupId(gid string) (*user.Group, error) {
	return user.LookupGroupId(gid)
}

// This function takes a file mode and returns a string representation of the file's permissions
// in the form of a 9-character string, where each character represents a permission
// (r, w, or x) for the owner, group, and others.
func getFilePermissions(mode fs.FileMode) string {
	// Initialize a string to represent the permission string,
	// with all characters initially set to "-"
	perm := "---------"
	for bit, val := range filePermMap { // Loop through the filePermMap (a map of permission bit masks to permission strings)
		// If the permission bit is set in the file mode, update the permission string with the corresponding permission string from filePermMap
		if mode&bit != 0 {
			perm = perm[:9-len(val)] + val + perm[10-len(val):]
		}
	}
	// Return the permission string
	return perm
}

// This function takes a file mode and returns a string representation of the file's type,
// using one of the following characters: "-", "d", "l", "b", "c", "p", or "s".
func getFileMode(mode fs.FileMode) string {
	switch { // Use a switch statement to determine the file type based on the file mode
	case mode.IsRegular():
		return "-" // regular file
	case mode.IsDir():
		return "d" // directory
	case mode&fs.ModeSymlink != 0:
		return "l" // symbolic link
	case mode&fs.ModeDevice != 0:
		return "b" // block device
	case mode&fs.ModeCharDevice != 0:
		return "c" // character device
	case mode&fs.ModeNamedPipe != 0:
		return "p" // named pipe (FIFO)
	case mode&fs.ModeSocket != 0:
		return "s" // socket
	default:
		return "?" // unknown type
	}
}
