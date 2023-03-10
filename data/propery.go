package data

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

const (
	tar   = "tar, tgz, arj, taz, lzh, lzma, tlz, txz, zip, z, Z, dz, gz, lz, xz, bz2, bz, tbz, tbz2, tz, deb, rpm, jar, rar, ace, zoo, cpio, 7z, rz"
	image = "jpg, jpeg, gif, bmp, pbm, pgm, ppm, tga, xbm, xpm, tif, tiff, png, svg, svgz, mng, pcx, mov, mpg, mpeg, m2v, mkv, ogm, mp4, m4v, mp4v, vob, qt, nuv, wmv, asf, rm, rmvb, flc, avi, fli, flv, gl, dl, xcf, xwd, yuv, cgm, emf, axv, anx, ogv, ogx"
)

/*GetUserName - see here - https://golang.org/pkg/os/user/#Current - about Current user */
func GetUserName(a, b uint32) (string, string) {
	x, err := user.LookupId(strconv.Itoa(int(a)))
	if err != nil {
		fmt.Println(string(a))
		fmt.Println("You are ghost wtf man")
		os.Exit(0)
	}
	//return in which group user in a member
	y, err2 := x.GroupIds()
	if err2 != nil || len(y) == 0 {
		fmt.Println("Strange things happends you are not member of any group WTF?")
		os.Exit(0)
	}
	//Return first group where are you member
	GroupName, _ := user.LookupGroupId(strconv.Itoa(int(b)))
	return x.Username, GroupName.Name
}

/*GetPathFile => Use this link to see better example - https://yourbasic.org/golang/current-directory/ */
func GetPathFile() string {
	x, err := os.Getwd() // <- os Getwd() - return path of current directory
	if err != nil {
		fmt.Println("This directory is prohibiten")
		os.Exit(0)
	}
	return x
}

// GetColor - Decide which color should be on output in terminal
func GetColor(tmp os.FileInfo) int {
	mode := tmp.Mode().String() // <- ex: drxrxx
	mode = strings.ToLower(mode)
	if mode[0] == 'd' {
		return 69 // <- Blue
	}
	if mode[0] == 'l' {
		return 37 // <- Light Blue
	}
	if mode[3] == 'x' {
		return 70 // <- Green
	}

	for _, i := range strings.Split(tar, ", ") {
		if strings.HasSuffix(tmp.Name(), i) {
			return 9 // <- red
		}
	}
	for _, i := range strings.Split(image, ", ") {
		if strings.HasSuffix(tmp.Name(), i) {
			return 97 // <- purple
		}
	}
	return 15 // <- white
}
