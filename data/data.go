package data

import (
	"os"
	"strings"
)

type Files struct {
	Stat os.FileInfo
	Name string
}

// Lets remove point (.) from hidden files
func delPoint(s string) string {
	if s[0] == '.' {
		return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(s[1:]), "-", ""), ".", "")
	}
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(s), "-", ""), ".", "")
}

// Reverse order
func Reverse(arr []Files) []Files {
	tmp := []Files{}
	for i := len(arr) - 1; i >= 0; i-- {
		tmp = append(tmp, arr[i])
	}
	return tmp
}

// Sort For sorting let's use quick sort algorithm from wiki
func Sort(arr []Files, Type int) {
	n := len(arr)
	for i := 0; i < n; i++ {
		for j := n/2 - 1 - i/2; j > -1; j-- {
			if 2*j+2 <= n-1-i {
				if delPoint(arr[2*j+1].Name) > delPoint(arr[2*j+2].Name) && Type == 1 || Type == 2 && !arr[2*j+1].Stat.ModTime().After(arr[2*j+2].Stat.ModTime()) {
					if delPoint(arr[j].Name) < delPoint(arr[2*j+1].Name) && Type == 1 || Type == 2 && !arr[2*j+1].Stat.ModTime().After(arr[j].Stat.ModTime()) {
						arr[j], arr[j*2+1] = arr[j*2+1], arr[j]
					}
				} else if delPoint(arr[j].Name) < delPoint(arr[j*2+2].Name) && Type == 1 || Type == 2 && !arr[2*j+2].Stat.ModTime().After(arr[j].Stat.ModTime()) {
					arr[j], arr[j*2+2] = arr[j*2+2], arr[j]
				}
			} else if 2*j+1 <= n-1-i {
				if delPoint(arr[j].Name) < delPoint(arr[2*j+1].Name) && Type == 1 || Type == 2 && !arr[2*j+1].Stat.ModTime().After(arr[j].Stat.ModTime()) {
					arr[j], arr[j*2+1] = arr[j*2+1], arr[j]
				}
			}
		}
		arr[0], arr[n-1-i] = arr[n-1-i], arr[0]
	}
}
