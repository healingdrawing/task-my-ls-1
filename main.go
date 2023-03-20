package main

import (
	"os"
	"syscall"
)

func main() {
	syscall.Exec("/bin/ls", append([]string{"ls"}, os.Args[1:]...), os.Environ())
}
