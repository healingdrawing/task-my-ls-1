package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

/*
Blue: Directory
Green: Executable or recognized data file
Cyan (Sky Blue): Symbolic link file
Yellow with black background: Device
Magenta (Pink): Graphic image file
Red: Archive file
Red with black background: Broken link
*/
var (
	showHidden  bool   // -a flag
	showAll     bool   // -A flag
	longListing bool   // -l flag
	reverse     bool   // -r flag
	sortByTime  bool   // -t flag
	recursive   bool   // -R flag
	humanize    bool   // -h flag
	dirPath     string // directory path to list

)

var (
	color   = false
	Reset   = "\033[0m"
	Blue    = "\033[34m"
	Green   = "\033[32m"
	Cyan    = "\033[36m"
	Yellow  = "\033[33m"
	Red     = "\033[31m"
	Magenta = "\033[35m"
)

//----------------------------------------------------------------
//colors and ls meaning:
/*
Blue: Directory
Green: Executable or recognized data file
Cyan (Sky Blue): Symbolic link file
Yellow with black background: Device
Magenta (Pink): Graphic image file
Red: Archive file
Red with black background: Broken link
*/
func main() {
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-a":
			showHidden = true
		case "-A":
			showAll = true
		case "-l":
			longListing = true
		case "-r":
			reverse = true
		case "-t":
			sortByTime = true
		case "-R":
			recursive = true
		case "-h":
			humanize = true
		case "-la":
			showHidden = true
			longListing = true
		case "--color=auto":
			color = true
			fmt.Println("color is on")
		default:
			dirPath = os.Args[i]
		}
	}
	if dirPath == "" {
		dirPath = "./"
	}
	dir, err := os.Open(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	dirEntries, err := dir.ReadDir(0)
	if err != nil {
		log.Fatal(err)
	}
	if sortByTime {
		sortByModificationTime(dirEntries)
	} else if reverse {
		sortReverse(dirEntries)
	} else {
		sortByName(dirEntries)
	}

	if longListing {
		var totalSize int64
		for _, entry := range dirEntries {
			if !showAll && !showHidden && entry.Name()[0] == '.' {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				log.Println(err)
				continue
			}
			totalSize += info.Size()
		}
		if humanize {
			fmt.Printf("total %s\n", humanizeBytes(totalSize))
		} else {
			fmt.Printf("total %d\n", totalSize/1024)
		}
		printLongListing(dirEntries)
	} else {
		printShortListing(dirEntries)
	}
	if recursive {
		root := "./"
		prefix := ""
		printDirContents(root, prefix)

	}
	// Use a slice to keep track of all the file paths we need to print
	var files []string

	// Handle any errors that occur during filepath.Walk
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	// Sort the files slice by name, if not already sorted by modification time or reverse sorted
	if !sortByTime && !reverse {
		sort.Strings(files)
	}
	// Print the long listing for each file in the files slice
	for _, file := range files {
		showLongListing([]string{file})
	}
	// Sort the dirEntries slice by modification time, if the -t flag is set
	if sortByTime {
		var fileInfo []os.FileInfo
		for _, entry := range dirEntries {
			// Convert each fs.DirEntry to an os.FileInfo using entry.Info()
			info, err := entry.Info()
			if err != nil {
				log.Println(err)
				continue
			}
			fileInfo = append(fileInfo, info)
		}
		sort.Slice(fileInfo, func(i, j int) bool {
			return fileInfo[i].ModTime().After(fileInfo[j].ModTime())
		})
		if reverse {
			reverseSlice(fileInfo)
		}
		// Print the long listing for each file in the fileInfo slice
		for _, info := range fileInfo {
			showLongListing([]string{info.Name()})
		}

	}

	if err != nil {
		log.Fatal(err)
	}

}

func printFiles(path string, files []os.FileInfo) {
	for _, f := range files {
		// add the path to the filename to get the full path
		//name := filepath.Join(path, f.Name())

		// remove the file extension from the name
		nameWithoutExt := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		// print the full path and name without extension
		fmt.Printf("%s\n", nameWithoutExt)
		break
	}
}

func printDirContents(path string, prefix string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != ".git" { // check if the directory is not ".git"
			subPath := filepath.Join(path, file.Name())
			fmt.Printf("%s%s:\n", prefix, subPath)
			printDirContents(subPath, prefix+"  ")
		}
	}
	printFiles(path, files)
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
			log.Println(err)
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

func reverseSlice(files []os.FileInfo) {
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}
}

func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func sortByName(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
}

func sortByModificationTime(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		time1, err := entries[i].Info()
		if err != nil {
			return false
		}
		time2, err := entries[j].Info()
		if err != nil {
			return false
		}
		return time1.ModTime().Before(time2.ModTime())
	})
}

func sortReverse(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})
}
func printShortListing(entries []os.DirEntry) {
	for _, entry := range entries {
		//if !showAll && !showHidden && entry.Name()[0] == '.' {
		if entry.Name() == "." || entry.Name() == ".." || (!showHidden && entry.Name()[0] == '.') {
			continue
		}
		if longListing {
			mode := entry.Type().String()
			info, err := entry.Info()
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Printf("%s %3d %s %s %6d %s %s\n",
				mode, info.Sys().(*syscall.Stat_t).Nlink,
				strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Uid)),
				strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Gid)),
				info.Size(), info.ModTime().Format("%+03d:%+03d" /*"Jan _2 15:04"*/), entry.Name())

		} else {
			fmt.Println(entry.Name())
		}
	}
}

func printLongListing(entries []os.DirEntry) {
	for _, entry := range entries {
		if !showAll && !showHidden && entry.Name()[0] == '.' {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			log.Println(err)
			continue
		}
		name := entry.Name()
		if info.IsDir() {
			if color {
				name = "\x1b[34m" + name + "\x1b[0m"
			}
		} else if info.Mode()&0111 != 0 { // check if file is executable
			if color {
				name = "\x1b[32m" + name + "\x1b[0m"
			}
		} else if info.Mode()&os.ModeSymlink != 0 {
			if color {
				name = "\x1b[36m" + name + "\x1b[0m"
			}
		} else if info.Mode()&os.ModeDevice != 0 {
			if color {
				name = "\x1b[43;30m" + name + "\x1b[0m"
			}
		} else if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") || strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".gif") || strings.HasSuffix(name, ".bmp") || strings.HasSuffix(name, ".tiff") || strings.HasSuffix(name, ".svg") {
			if color {
				name = "\x1b[35m" + name + "\x1b[0m"
			}
		} else if strings.HasSuffix(name, ".zip") || strings.HasSuffix(name, ".rar") || strings.HasSuffix(name, ".tar") || strings.HasSuffix(name, ".gz") || strings.HasSuffix(name, ".7z") || strings.HasSuffix(name, ".bz2") || strings.HasSuffix(name, ".xz") {
			if color {
				name = "\x1b[31m" + name + "\x1b[0m"
			}
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			log.Println("Error getting system-specific file info")
			continue
		}
		mode := info.Mode().String()
		fmt.Printf("%s %3d %s %s %6d %s %s\n",
			mode, stat.Nlink, getUserName(stat.Uid), getGroupName(stat.Gid),
			info.Size(), info.ModTime().Format("Jan _2 15:04"), name)
	}
}

func getUserName(uid uint32) string {
	user, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return fmt.Sprintf("%d", uid)
	}
	return user.Username
}

func getGroupName(gid uint32) string {
	group, err := user.LookupGroupId(strconv.Itoa(int(gid)))
	if err != nil {
		return fmt.Sprintf("%d", gid)
	}
	return group.Name
}

//----------------------------------------------------------------
//To print the output of the ls -R command using a custom implementation,
// you need to traverse the directory tree and print the names of files and directories recursively.
// The -R flag indicates that you need to print the contents of subdirectories as well.
/*
func main() {
    args := os.Args[1:]
    if len(args) == 0 {
        args = []string{"."}
    }
    for _, arg := range args {
        err := listFiles(arg, "")
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}*/
/*
func listFiles(path string, indent string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		fmt.Println(indent + path + ":")
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			if file.IsDir() {
				err = listFiles(filepath.Join(path, file.Name()), indent+"  ")
			} else {
				fmt.Printf("%s%-20s %10d %s\n", indent+"  ", file.Name(), file.Size(), file.ModTime().Format("Jan 02 15:04"))
			}
		}
	} else {
		fmt.Printf("%-20s %10d %s\n", path, info.Size(), info.ModTime().Format("Jan 02 15:04"))
	}

	return nil
}
*/
/*----------------------------------------------------------------
Not cleaned up
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

/*
Blue: Directory
Green: Executable or recognized data file
Cyan (Sky Blue): Symbolic link file
Yellow with black background: Device
Magenta (Pink): Graphic image file
Red: Archive file
Red with black background: Broken link
*/
/*
var (
	showHidden  bool   // -a flag
	showAll     bool   // -A flag
	longListing bool   // -l flag
	reverse     bool   // -r flag
	sortByTime  bool   // -t flag
	recursive   bool   // -R flag
	humanize    bool   // -h flag
	dirPath     string // directory path to list

)

var (
	color   = false
	Reset   = "\033[0m"
	Blue    = "\033[34m"
	Green   = "\033[32m"
	Cyan    = "\033[36m"
	Yellow  = "\033[33m"
	Red     = "\033[31m"
	Magenta = "\033[35m"
)
*/
//----------------------------------------------------------------
//colors and ls meaning:
/*
Blue: Directory
Green: Executable or recognized data file
Cyan (Sky Blue): Symbolic link file
Yellow with black background: Device
Magenta (Pink): Graphic image file
Red: Archive file
Red with black background: Broken link
*/
/*----------------------------------------------------------------
func main() {
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-a":
			showHidden = true
		case "-A":
			showAll = true
		case "-l":
			longListing = true
		case "-r":
			reverse = true
		case "-t":
			sortByTime = true
		case "-R":
			recursive = true
		case "-h":
			humanize = true
		case "-la":
			showHidden = true
			longListing = true
		case "--color=auto":
			color = true
			fmt.Println("color is on")
		default:
			dirPath = os.Args[i]
		}
	}
	if dirPath == "" {
		dirPath = "./"
	}
	dir, err := os.Open(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	dirEntries, err := dir.ReadDir(0)
	if err != nil {
		log.Fatal(err)
	}
	if sortByTime {
		sortByModificationTime(dirEntries)
	} else if reverse {
		sortReverse(dirEntries)
	} else {
		sortByName(dirEntries)
	}

	if longListing {
		var totalSize int64
		for _, entry := range dirEntries {
			if !showAll && !showHidden && entry.Name()[0] == '.' {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				log.Println(err)
				continue
			}
			totalSize += info.Size()
		}
		if humanize {
			fmt.Printf("total %s\n", humanizeBytes(totalSize))
		} else {
			fmt.Printf("total %d\n", totalSize/1024)
		}
		printLongListing(dirEntries)
	} else {
		printShortListing(dirEntries)
	}
	if recursive {
		root := "./"
		prefix := ""
		printDirContents(root, prefix)

		//print("hello")
		//printDirContents(".", "")
		//printDirContents(".", 0, 2, "")
		//showLongListingRecursive(dirPath)
		//return
	}
	// Use a slice to keep track of all the file paths we need to print
	var files []string
	/*
		// Check if the file is a hidden file and if the showHidden flag is false, then skip it
		for _, entry := range dirEntries {
			info, err := entry.Info()
			fmt.Println("hello")
			if err != nil {
				log.Println(err)
				continue
			}
			if !showHidden && strings.HasPrefix(info.Name(), ".") {
				if info.IsDir() {
					fmt.Println("hello1")
					continue // skip the directory if it's hidden
				}
				continue // skip the file if it's hidden
			}
			if info.IsDir() {
				// If the file is a directory, add it to the files slice to be printed later
				files = append(files, entry.Name())
				continue // don't print the directory now, wait until we iterate over its contents
			}
			// If the file is not a directory, add it to the files slice to be printed later
			files = append(files, entry.Name())
		}*
	// Handle any errors that occur during filepath.Walk
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	// Sort the files slice by name, if not already sorted by modification time or reverse sorted
	if !sortByTime && !reverse {
		sort.Strings(files)
	}
	// Print the long listing for each file in the files slice
	for _, file := range files {
		showLongListing([]string{file})
	}
	// Sort the dirEntries slice by modification time, if the -t flag is set
	if sortByTime {
		var fileInfo []os.FileInfo
		for _, entry := range dirEntries {
			// Convert each fs.DirEntry to an os.FileInfo using entry.Info()
			info, err := entry.Info()
			if err != nil {
				log.Println(err)
				continue
			}
			fileInfo = append(fileInfo, info)
		}
		sort.Slice(fileInfo, func(i, j int) bool {
			return fileInfo[i].ModTime().After(fileInfo[j].ModTime())
		})
		if reverse {
			reverseSlice(fileInfo)
		}
		// Print the long listing for each file in the fileInfo slice
		for _, info := range fileInfo {
			showLongListing([]string{info.Name()})
		}

	}----------------------------------------------------------------
*/
/*
		// If the -R flag is set, show long listing recursively.
		if recursive {
			var files []string
			for _, f := range files {
				showLongListingRecursive(f) // call function to show long listing recursively
			}
			dir, err := os.Open(dirPath)
			if err != nil {
				log.Fatal(err)
			}
			defer dir.Close()
			dirEntries, err := dir.ReadDir(0)
			if err != nil {
				log.Fatal(err)
			}

			if sortByTime {
				sort.Slice(dirEntries, func(i, j int) bool {
					return dirEntries[i].ModTime().After(dirEntries[j].ModTime())
				})
				if reverse {
					reverseSlice(dirEntries)
				}
			} else if reverse {
				sortReverse(dirEntries)
			} else {
				sortByName(dirEntries)
			}

			if longListing {
				var totalSize int64
				for _, entry := range dirEntries {
					if !showAll && !showHidden && entry.Name()[0] == '.' {
						continue
					}
					info, err := entry.Info()
					if err != nil {
						log.Println(err)
						continue
					}
					totalSize += info.Size()
				}
				if humanize {
					fmt.Printf("total %s\n", humanizeBytes(totalSize))
				} else {
					fmt.Printf("total %d\n", totalSize/1024)
				}
				printLongListing(dirEntries)
			} else {
				printShortListing(dirEntries)
			}
		}*
	if err != nil {
		log.Fatal(err)
	}
	here
----------------------------------------------------------------
}*/

/*
func printDirContents(path string, depth int, maxDepth int, prefix string) {
	if depth > maxDepth {
		return
	}

	fmt.Println(prefix + path + ":")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			printDirContents(filepath.Join(path, file.Name()), depth+1, maxDepth, prefix+"  ")
		} else {
			fmt.Printf("%s%-20s%10d %s\n", prefix+"  ", file.Mode(), file.Size(), file.Name())
		}
	}
}*/

/*
func printDirContents(path string, depth int, maxDepth int) {
	if depth > maxDepth {
		return
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println(strings.Repeat(" ", depth*2) + file.Name() + ":")
			printDirContents(filepath.Join(path, file.Name()), depth+1, maxDepth)
		} else {
			fmt.Println(strings.Repeat(" ", depth*2) + file.Name())
		}
	}
}

*/
/*
// to print recursive listing of directories
func printDirContents(path string, prefix string) {
	printFiles(prefix+path+":", []os.FileInfo)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			fmt.Println(path+"/"+file.Name(), prefix+"  ")
		}
	}
}
*/
/*
func printFiles(path string, files []os.FileInfo) {
	for _, f := range files {
		// add the path to the filename to get the full path
		name := filepath.Join(path, f.Name())
		// print the full path and name
		//fmt.Println(name)

		// if this is a directory, print its contents recursively
		if f.IsDir() {
			_, err := os.ReadDir(name)
			if err == nil {

				//fmt.Println(name)
				//printFiles(name, subFiles)

			}
		}
	}
}*/
/*----------------------------------------------------------------
from here
func printFiles(path string, files []os.FileInfo) {
	for _, f := range files {
		// add the path to the filename to get the full path
		//name := filepath.Join(path, f.Name())

		// remove the file extension from the name
		nameWithoutExt := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		// print the full path and name without extension
		fmt.Printf("%s\n", nameWithoutExt)
		break
		/*
			// if this is a directory, print its contents recursively
			if f.IsDir() {
				subFiles, err := os.ReadDir(name)
				if err == nil {
					var subFilesInfo []os.FileInfo
					for _, subFile := range subFiles {
						// convert fs.DirEntry to os.FileInfo
						subFileInfo, err := subFile.Info()
						if err == nil {
							subFilesInfo = append(subFilesInfo, subFileInfo)
						}
					}
					printFiles(name, subFilesInfo)
				}
			}*
	}
}

func printDirContents(path string, prefix string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != ".git" { // check if the directory is not ".git"
			subPath := filepath.Join(path, file.Name())
			fmt.Printf("%s%s:\n", prefix, subPath)
			printDirContents(subPath, prefix+"  ")
		}
	}
	printFiles(path, files)
}
here----------------------------------------------------------------
*/
/*
func printFiles(path string, files []os.FileInfo) {
	for _, f := range files {
		// add the path to the filename to get the full path
		name := filepath.Join(path, f.Name())
		// print the full path and name
		fmt.Printf("%s:%s\n", f.Name(), filepath.Ext(f.Name()))

		// if this is a directory, print its contents recursively
		if f.IsDir() {
			subFiles, err := os.ReadDir(name)
			if err == nil {
				var subFilesInfo []os.FileInfo
				for _, subFile := range subFiles {
					// convert fs.DirEntry to os.FileInfo
					subFileInfo, err := subFile.Info()
					if err == nil {
						subFilesInfo = append(subFilesInfo, subFileInfo)
					}
				}
				printFiles(name, subFilesInfo)
			}
		}
	}
}
func printDirContents(path string, prefix string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != ".git" { // check if the directory is not ".git"
			subPath := filepath.Join(path, file.Name())
			fmt.Printf("%s%s:\n", prefix, subPath)
		}
	}

	printFiles(path, files)
}*/

/*
func printDirContents(path string, prefix string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// add the path to the filename to get the full path
		name := filepath.Join(path, file.Name())

		// if this is a directory, print its contents recursively
		if file.IsDir() {
			fmt.Println(prefix + name + ":")
			//printDirContents(name, prefix+"  ")
		} /* else {
			// print the full path and name
			fmt.Println(name)
		}
	}
}*/
/*
func printDirContents(path string, prefix string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	/*
		// print files in the current directory
		for _, f := range files {
			// add the path to the filename to get the full path
			name := filepath.Join(path, f.Name())
			// print the full path and name
			fmt.Println(name)
		}*/

// recursively print contents of directories in the current directory
/*for _, file := range files {
		if file.IsDir() {
			subDirPath := filepath.Join(path, file.Name())
			fmt.Println(prefix + subDirPath + ":" + filepath.Base(file.Name()) + "\n")
			//printDirContents(subDirPath, prefix+"  ")
		}
	}
}*/
/*
func printDirContents(path string, prefix string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			//fmt.Println(prefix + path + "/" + file.Name() + ":")
			_, err := os.ReadDir(path + "/")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(prefix + path + "/" + file.Name() + ":")
			//printFiles(file.Name(), subFiles)

		} /*
			// add the path to the filename to get the full path
			//name := filepath.Join(path, file.Name())

			// if this is a directory, print its contents recursively
			if file.IsDir() {
				//fmt.Println(prefix + name + ":")

				// recursively call printDirContents for the subdirectory
				printDirContents(path+"/"+file.Name(), prefix+"  ")
				//fmt.Println(name, prefix+"  ")
			}*
	}

	// call printFiles to print the files in this directory
	printFiles(path, files)

}*/

/*
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
*/
/*----------------------------------------------------------------
// This function takes a list of file paths as input and prints
// the long listing of each file (mode, size, modification time, and path)
func showLongListing(files []string) {
	for _, file := range files {
		// Get the file information for the current file
		fileInfo, err := os.Stat(file)
		if err != nil {
			// Handle any errors that occur while reading the file info
			fmt.Fprintf(os.Stderr, "Error reading file info: %v\n", err)
			log.Println(err)
			continue
		}
		// Print the long listing for the current file (mode, size, modification time, and path)
		fmt.Printf("%s\t%d\t%s\t%s\n",
			fileInfo.Mode(),
			fileInfo.Size(),
			fileInfo.ModTime().Format("Jan 02 15:04"),
			file)
	}
}----------------------------------------------------------------
*/
/*
func reverseSlice(entries []os.DirEntry) {
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
}*/
/*----------------------------------------------------------------
here
func reverseSlice(files []os.FileInfo) {
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}
}

func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func sortByName(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
}

func sortByModificationTime(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		time1, err := entries[i].Info()
		if err != nil {
			return false
		}
		time2, err := entries[j].Info()
		if err != nil {
			return false
		}
		return time1.ModTime().Before(time2.ModTime())
	})
}

func sortReverse(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})
}
func printShortListing(entries []os.DirEntry) {
	for _, entry := range entries {
		//if !showAll && !showHidden && entry.Name()[0] == '.' {
		if entry.Name() == "." || entry.Name() == ".." || (!showHidden && entry.Name()[0] == '.') {
			continue
		}
		if longListing {
			mode := entry.Type().String()
			info, err := entry.Info()
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Printf("%s %3d %s %s %6d %s %s\n",
				mode, info.Sys().(*syscall.Stat_t).Nlink,
				strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Uid)),
				strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Gid)),
				info.Size(), info.ModTime().Format("%+03d:%+03d" /*"Jan _2 15:04"*), entry.Name())

		} else {
			fmt.Println(entry.Name())
		}
	}
}

func printLongListing(entries []os.DirEntry) {
	for _, entry := range entries {
		if !showAll && !showHidden && entry.Name()[0] == '.' {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			log.Println(err)
			continue
		}
		name := entry.Name()
		if info.IsDir() {
			if color {
				name = "\x1b[34m" + name + "\x1b[0m"
			}
		} else if info.Mode()&0111 != 0 { // check if file is executable
			if color {
				name = "\x1b[32m" + name + "\x1b[0m"
			}
		} else if info.Mode()&os.ModeSymlink != 0 {
			if color {
				name = "\x1b[36m" + name + "\x1b[0m"
			}
		} else if info.Mode()&os.ModeDevice != 0 {
			if color {
				name = "\x1b[43;30m" + name + "\x1b[0m"
			}
		} else if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") || strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".gif") || strings.HasSuffix(name, ".bmp") || strings.HasSuffix(name, ".tiff") || strings.HasSuffix(name, ".svg") {
			if color {
				name = "\x1b[35m" + name + "\x1b[0m"
			}
		} else if strings.HasSuffix(name, ".zip") || strings.HasSuffix(name, ".rar") || strings.HasSuffix(name, ".tar") || strings.HasSuffix(name, ".gz") || strings.HasSuffix(name, ".7z") || strings.HasSuffix(name, ".bz2") || strings.HasSuffix(name, ".xz") {
			if color {
				name = "\x1b[31m" + name + "\x1b[0m"
			}
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			log.Println("Error getting system-specific file info")
			continue
		}
		mode := info.Mode().String()
		fmt.Printf("%s %3d %s %s %6d %s %s\n",
			mode, stat.Nlink, getUserName(stat.Uid), getGroupName(stat.Gid),
			info.Size(), info.ModTime().Format("Jan _2 15:04"), name)
	}
}

func getUserName(uid uint32) string {
	user, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return fmt.Sprintf("%d", uid)
	}
	return user.Username
}

func getGroupName(gid uint32) string {
	group, err := user.LookupGroupId(strconv.Itoa(int(gid)))
	if err != nil {
		return fmt.Sprintf("%d", gid)
	}
	return group.Name
}

*/
