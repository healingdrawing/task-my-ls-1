package file

/*
import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

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

func listFiles(dir string, files *[]fs.DirEntry) {
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			*files = append(*files, d)
		}
		return nil
	})
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

func GetOwnerGroup(path string) (string, string, error) {
    stat, err := os.Stat(path)
    if err != nil {
        return "", "", err
    }

    sysStat := stat.Sys().(*syscall.Stat_t)
    uid := sysStat.Uid
    gid := sysStat.Gid

    user, err := user.LookupId(strconv.Itoa(int(uid)))
    if err != nil {
        return "", "", err
    }

    group, err := LookupGroupId(strconv.Itoa(int(gid)))
    if err != nil {
        return "", "", err
    }

    return user.Username, group.Name, nil
}




/*
func GetOwnerGroup(file string) (string, string, error) {
    info, err := os.Stat(file)
    if err != nil {
        return "", "", err
    }

    owner, group := file.GetOwnerGroup(info)
    return owner, group, nil
}*/

/*
func GetOwnerGroup(info fs.FileInfo) (string, string, error) {
    stat, ok := info.Sys().(*syscall.Stat_t)
    if !ok {
        return "", "", fmt.Errorf("failed to get underlying syscall.Stat_t struct")
    }

    uid := stat.Uid
    gid := stat.Gid

    user, err := user.LookupId(strconv.Itoa(int(uid)))
    if err != nil {
        return "", "", err
    }

    group, err := LookupGroupId(strconv.Itoa(int(gid)))
    if err != nil {
        return "", "", err
    }

    return user.Username, group.Name, nil
}*/

/*

func GetOwnerGroup(info fs.DirEntry) (string, string, error) {
	uid := info.Sys().(*syscall.Stat_t).Uid
	gid := info.Sys().(*syscall.Stat_t).Gid
	user, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return "", "", err
	}

	group, err := user.GroupIds()
	if err != nil {
		return "", "", err
	}

	groupInfo, err := LookupGroupId(group[0])
	if err != nil {
		return "", "", err
	}

	return user.Username, groupInfo.Name, nil
}*/
/*
// LookupGroupId is a wrapper around the user.LookupGroupId function.
// It takes a group ID string as an argument and returns a user.Group object and an error.
func LookupGroupId(gid string) (*user.Group, error) {
    return user.LookupGroupId(gid)
}

// This function takes a file mode and returns a string representation of the file's permissions
// in the form of a 9-character string, where each character represents a permission
// (r, w, or x) for the owner, group, and others.
func GetFilePermissions(mode fs.FileMode) string {
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
func GetFileMode(mode fs.FileMode) string {
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
*/
