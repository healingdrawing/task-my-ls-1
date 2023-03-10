package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"

	"my-ls-1/data"
)

const (
	format = "\033[38;5;%dm%s\033[39;49m"
)

var (
	is_l       = false
	is_r       = false
	is_a       = false
	is_A       = false
	is_R       = false
	is_t       = false
	total      int64
	mx_blocks  = 0
	mx_bytes   = 0
	mix        = 0
	answer_cnt = 0
	ok2        = 0
	mx_user    = 0
	mx_group   = 0
	mx_mode    = 10
	mx_minor   = 0
	mx_major   = 0
)

/*___________________________________________________________*/
//Use this link to know better - https://golang.org/pkg/os/#FileInfo
//To see more information about table on ls -l use this link - https://stackoverflow.com/questions/43181806/result-of-linux-ls-lisa-command

func printLsL(tmp os.FileInfo, Name, path string) {
	mode := tmp.Mode().String()  // <- ex: drxrxx
	mode = strings.ToLower(mode) // <- fix small bug: on links, because on link file first rune is L, not l
	if mode[1] == 'c' {
		mode = mode[1:]
	}
	if mode[1] == 't' {
		mode = mode[:1] + mode[2:]
	}
	if mode == "crw-rw-rw-" || mode == "crw-rw-r--" {
		mode += "+"
	}

	stat, _ := tmp.Sys().(*syscall.Stat_t) // <- how to parse interface - https://stackoverflow.com/questions/28339240/get-file-inode-in-go
	// fmt.Println(stat.Rdev / 256)
	// fmt.Println(stat.Rdev % 256)
	user, group := data.GetUserName(stat.Uid, stat.Gid)
	size := tmp.Size()                                 //<- File size in byte
	date := tmp.ModTime()                              // <- Last modifed time
	month := date.Month().String()[:3]                 // <- Month
	day := date.Day()                                  // <- Day
	h, m, _ := date.Clock()                            // <- hours minutes
	s := fmt.Sprintf(format, data.GetColor(tmp), Name) // <- our output

	if mode[0] == 'l' { // origin file of link file
		origin, err := os.Readlink(path + "/" + tmp.Name())
		path_or := origin
		if origin[0] != '/' {
			path_or = path + "/" + origin
		}
		f2, err2 := os.Open(path_or)

		if err != nil || err2 != nil { // <- if link file is broken
			arr := strings.Split(origin, "/")
			s = fmt.Sprintf(format, 9, tmp.Name()) + " -> " + fmt.Sprintf(format, 9, arr[len(arr)-1])
		} else {
			stat, _ := f2.Stat()
			s += " -> " + fmt.Sprintf(format, data.GetColor(stat), origin)
		}
		defer f2.Close()
	}

	now := time.Now()
	now = now.AddDate(0, -6, 0)
	/* Adding spaces for output format using max length of each row */
	for i := len(user); i < mx_user; i++ {
		user += " "
	}
	for i := len(group); i < mx_group; i++ {
		group += " "
	}
	for i := len(mode); i <= mx_mode; i++ {
		mode += " "
	}
	/* Pinting our date with special formal like in ls */
	if now.Before(date) {
		if mode[0] != 'c' {
			printFormat := "%s %" + strconv.Itoa(mx_blocks) + "d %s %s %" + strconv.Itoa(mx_bytes) + "d %s %2d %02d:%02d %s\n"
			fmt.Printf(printFormat, mode, stat.Nlink, user, group, size, month, day, h, m, s)
		} else {
			printFormat := "%s %" + strconv.Itoa(mx_blocks) + "d %s %s %" + strconv.Itoa(mx_minor) + "d, %" + strconv.Itoa(mx_major) + "d %s %2d  %d %s\n"
			fmt.Printf(printFormat, mode, stat.Nlink, user, group, data.CountDivision(uint64(stat.Rdev)), data.CountMod(uint64(stat.Rdev)), month, day, date.Year(), s)
		}
	} else {
		if mode[0] != 'c' {
			printFormat := "%s %" + strconv.Itoa(mx_blocks) + "d %s %s %" + strconv.Itoa(mx_bytes) + "d %s %2d  %d %s\n"
			fmt.Printf(printFormat, mode, stat.Nlink, user, group, size, month, day, date.Year(), s)
		} else {
			printFormat := "%s %" + strconv.Itoa(mx_blocks) + "d %s %s %" + strconv.Itoa(mx_minor) + "d, %" + strconv.Itoa(mx_major) + "d %s %2d  %d %s\n"

			fmt.Printf(printFormat, mode, stat.Nlink, user, group, data.CountDivision(uint64(stat.Rdev)), data.CountMod(uint64(stat.Rdev)), month, day, date.Year(), s)
		}
	}
}

/* Parse - parsing  files on path */
func Parse(path string, cut_path int) {

	total = 0
	mx_blocks = 0
	mx_bytes = 0
	var arr []data.Files
	f, err := os.Open(path)

	if err != nil { // if file not exist
		fmt.Printf("ls: cannot access '%s': No such file or directory\n", path[cut_path+1:])
		return
	}
	//defer f.Close()

	tmp_arr, err := f.Readdir(0) // <- Read all files
	if err != nil {
		//fmt.Println("no such file")
		//return
		//log.Fatal(err)
	}
	dir, _ := f.Stat()
	_, err2 := os.Readlink(path)
	dd2, _ := os.Lstat(path)
	if dir.IsDir() && err2 == nil {
		pp := path
		if len(pp) > 0 && pp[len(pp)-1] == '/' { // <- removing useless '/'
			pp = pp[:len(pp)-1]
		}
		tmp := strings.Split(pp, "/")
		pp = strings.Join(tmp[:len(tmp)-1], "/")

		printLsL(dd2, dd2.Name(), pp)
		return
	}
	if !dir.IsDir() { // if file is not dir
		if ok2 == 1 { // <- special output format in ls
			fmt.Println()
			ok2 = 0
		}
		if answer_cnt != 0 && is_l {
			fmt.Println()
		}
		if is_l {
			if err2 == nil {
				dir = dd2
			}
			pp := path
			if len(pp) > 0 && pp[len(pp)-1] == '/' { // <- removing useless '/'
				pp = pp[:len(pp)-1]
			}
			tmp := strings.Split(pp, "/")
			pp = strings.Join(tmp[:len(tmp)-1], "/")
			printLsL(dir, dir.Name(), pp)
		} else {
			s := fmt.Sprintf(format, data.GetColor(dir), path[cut_path+1:])
			if ok2 == 2 {
				fmt.Print(" ")
			}
			fmt.Printf("%s", s)
			ok2 = 2
		}
		answer_cnt++
		return
	}

	if ok2 == 2 { // <- special output format in ls
		ok2 = 0
		fmt.Println()
	}

	if is_R && mix == 0 { // <- Special output for only R tag on ls
		if len(path) != cut_path {
			fmt.Println()
			fmt.Printf("%s:\n", "."+path[cut_path:])
		} else {
			fmt.Printf("%s:\n", "./")
		}

	}
	if mix > 1 || mix != 0 && is_R { // <- Special output on R tag with arguments in ls
		if answer_cnt != 0 {
			fmt.Println()
		}
		fmt.Println(path[cut_path+1:] + ":")
		answer_cnt++
	}

	pp := path
	if len(pp) > 0 && pp[len(pp)-1] == '/' { // <- removing useless '/'
		pp = pp[:len(pp)-1]
	}

	tmp := strings.Split(pp, "/")
	pp = strings.Join(tmp[:len(tmp)-1], "/")

	if len(pp) == 0 {
		pp = "/"
	}

	tmp_f, err := os.Open(pp) // <- parent directory of the current path
	if err != nil {
		log.Fatal()
	}

	dir2, _ := tmp_f.Stat()
	arr = append(arr, data.Files{Stat: dir, Name: "."})   // <- parent directory
	arr = append(arr, data.Files{Stat: dir2, Name: ".."}) // <- out directory

	n := len(tmp_arr)

	for i := 0; i < n; i++ { // <- Using our structure for sorting files by name or date
		arr = append(arr, data.Files{Stat: tmp_arr[i], Name: tmp_arr[i].Name()})
	}

	tmp_arr2 := []data.Files{}
	for i := range arr { // <- Count total
		if !is_a && arr[i].Name[0] == '.' && (is_A && (arr[i].Name == "." || arr[i].Name == "..") || !is_A) {
			continue
		}
		stat, _ := arr[i].Stat.Sys().(*syscall.Stat_t) // <- Parsing stats like above
		user, group := data.GetUserName(stat.Uid, stat.Gid)

		data.Max(&mx_user, len(user))
		data.Max(&mx_group, len(group))

		if arr[i].Stat.Mode().String() == "dcrw-rw-rw-" || arr[i].Stat.Mode().String() == "dcrw-rw-r--" { // <- Special output on ls
			mx_mode = 11
		}

		data.Max(&mx_blocks, data.NumberLen(int(stat.Nlink)))
		data.Max(&mx_bytes, data.NumberLen(int(stat.Size)))

		total += stat.Blocks

		data.Max(&mx_minor, data.NumberLen(int(data.CountDivision(uint64(stat.Rdev)))))
		data.Max(&mx_major, data.NumberLen(int(data.CountMod(uint64(stat.Rdev)))))

		if arr[i].Stat.Mode().String()[1] == 'c' {
			data.Max(&mx_bytes, mx_major+mx_minor+2)
		}
		tmp_arr2 = append(tmp_arr2, arr[i])
	}

	total++
	total /= 2
	arr = tmp_arr2
	/* checking flags on our ls */

	if is_t {
		data.Sort(arr, 2)
	} else {
		data.Sort(arr, 1)
	}
	if is_r {
		arr = data.Reverse(arr)
	}
	if is_l {
		fmt.Printf("total %d\n", total)
		for i := range arr {
			if !is_a && arr[i].Name[0] == '.' && (is_A && (arr[i].Name == "." || arr[i].Name == "..") || !is_A) {
				continue
			}
			printLsL(arr[i].Stat, arr[i].Name, path)
		}
	} else {
		ok := false
		for i := range arr {
			if !is_a && arr[i].Name[0] == '.' && (is_A && (arr[i].Name == "." || arr[i].Name == "..") || !is_A) {
				continue
			}
			s := fmt.Sprintf(format, data.GetColor(arr[i].Stat), arr[i].Name)
			if ok {
				fmt.Print(" ")
			}
			fmt.Printf("%s", s)

			ok = true

		}
		if ok {
			fmt.Println()
			ok = false
		}
	}
	if is_R {
		for i := range arr {
			if arr[i].Name != "." && arr[i].Name != ".." {
				if arr[i].Stat.IsDir() {
					add := "/"
					if len(path) > 0 && path[len(path)-1] == '/' {
						add = ""
					}
					Parse(path+add+arr[i].Name, cut_path)
				}
			}
		}
	}
}

/* Priority in sorting elements */
func priority(x, y, z string) bool {
	fy, err := os.Open(z + "/" + y)
	if len(y) > 0 && y[0] == '/' {
		fy, err = os.Open(y)
		y = y[1:]
	}
	statY, _ := fy.Stat()
	fx, err2 := os.Open(z + "/" + x)
	if len(x) > 0 && x[0] == '/' {
		fx, err2 = os.Open(x)
		x = x[1:]
	}
	statX, _ := fx.Stat()
	ok := false
	if err != nil || err2 != nil && err == nil {
		ok = false
	} else if !statY.IsDir() && statX.IsDir() {
		ok = true
	} else if x > y && err == nil && err2 == nil && statX.IsDir() == statY.IsDir() {
		ok = true
	}
	fy.Close()
	fx.Close()
	return ok
}

func main() {
	arg := os.Args[1:]
	path := data.GetPathFile()
	arr := []string{}

	/* Get flags */
	for _, i := range arg {
		if len(i) != 0 && i[0] == '-' {
			if len(i) == 1 {
				arr = append(arr, i)
				mix++
			}
			for j := 1; j < len(i); j++ {
				if i[j] == 'R' {
					is_R = true
				}
				if i[j] == 'a' {
					is_a = true
					is_A = false
				}
				if i[j] == 'A' {
					is_A = true
					is_a = false
				}
				if i[j] == 'r' {
					is_r = true
				}
				if i[j] == 't' {
					is_t = true
				}
				if i[j] == 'l' {
					is_l = true
				}
			}
		} else {
			arr = append(arr, i)
			mix++
		}
	}
	/* Specials signature for ~ */
	for i := 0; i < len(arr); i++ {
		if arr[i] == "~" {
			user_name, _ := user.Current()
			arr[i] = "/home/" + user_name.Name
		}
	}

	/* sort in priority of processing ls on dirs or files */
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if priority(arr[i], arr[j], path) {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	/*if the arguments in ls */
	for _, i := range arr {
		if len(i) > 0 && i[0] == '/' {
			Parse(i, -1)
		} else {
			if i == "" {
				fmt.Printf("ls: cannot access '': No such file or directory\n")
				answer_cnt++
			} else {
				Parse(path+"/"+i, len(path))
			}
		}
	}
	if mix == 0 {
		Parse(path, len(path))
	}
}
