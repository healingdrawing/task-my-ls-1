package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
)

var (
	showHidden  bool   // -a flag
	longListing bool   // -l flag
	reverse     bool   // -r flag
	sortByTime  bool   // -t flag
	recursive   bool   // -R flag
	dirPath     string // directory path to list
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-a":
			showHidden = true
		case "-l":
			longListing = true
		case "-r":
			reverse = true
		case "-t":
			sortByTime = true
		case "-R":
			recursive = true
		case "-la":
			showHidden = true
			longListing = true
		default:
			dirPath = os.Args[i]
		}
	}
	if dirPath == "" {
		dirPath = "."
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
	var total int64
	if longListing {
		for _, entry := range dirEntries {
			if !showHidden && entry.Name()[0] == '.' {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				log.Println(err)
				continue
			}
			total += info.Size()
		}
		fmt.Printf("total %d\n", total/1024)
		printLongListing(dirEntries)
	} else {
		printShortListing(dirEntries)
	}
	if recursive {
		err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error walking path %q: %v\n", path, err)
				return nil
			}
			if info.IsDir() {
				if info.Name() != "." && info.Name() != ".." {
					fmt.Printf("\n%s:\n", path)
					dir, err := os.Open(path)
					if err != nil {
						log.Printf("Error opening directory %q: %v\n", path, err)
						return nil
					}
					defer dir.Close()
					dirEntries, err := dir.ReadDir(0)
					if err != nil {
						log.Printf("Error reading directory %q: %v\n", path, err)
						return nil
					}
					if sortByTime {
						sortByModificationTime(dirEntries)
					} else if reverse {
						sortReverse(dirEntries)
					} else {
						sortByName(dirEntries)
					}
					if longListing {
						printLongListing(dirEntries)
					} else {
						printShortListing(dirEntries)
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
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
		if !showHidden && entry.Name()[0] == '.' {
			continue
		}
		if longListing {
			fmt.Printf("-l %s\n", entry.Name())
		} else {
			fmt.Println(entry.Name())
		}
	}
}

func printLongListing(entries []os.DirEntry) {
	for _, entry := range entries {
		if !showHidden && entry.Name()[0] == '.' {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			log.Println(err)
			continue
		}
		var stat syscall.Stat_t
		err = syscall.Stat(entry.Name(), &stat)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		mode := info.Mode().String()
		size := info.Size()
		uid := int(stat.Uid)
		gid := int(stat.Gid)
		userInfo, err := user.LookupId(strconv.Itoa(uid))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		groupInfo, err := user.LookupGroupId(strconv.Itoa(gid))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		timestamp := info.ModTime().Format("Jan _2 15:04")
		fmt.Printf("%s %3d %s %s %6d %s %s\n", mode, info.Sys().(*syscall.Stat_t).Nlink, userInfo.Username, groupInfo.Name, size, timestamp, entry.Name())
	}
}
