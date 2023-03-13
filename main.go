package main

import (
	"os"
	"syscall"
)

// version with extra "exit status 1" string at the end of the output in comparing with the built-in ls command. But mystically, today it works without exit status 1 string. And i do not remember that i changed something in the code. Perhaps the go build command can be the reason. Before i tested it with "go run . " .
func main() {
	syscall.Exec("/bin/ls", append([]string{"ls"}, os.Args[1:]...), os.Environ())
}

/* version without extra "exit status 1" string at the end of the output in comparing with the built-in ls command
all this shit requires to cut the "exit status 1" string from the output of the command
*/

// package main

// import (
// 	"fmt"
// 	"os"
// 	"strings"
// 	"syscall"
// )

// func main() {
// 	// Create a pipe to capture the output of the command
// 	r, w, err := os.Pipe()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	// Fork the current process
// 	pid, err := syscall.ForkExec("/bin/ls", os.Args, &syscall.ProcAttr{
// 		Files: []uintptr{0, w.Fd(), 2},
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	// Wait for the command to finish and close the write end of the pipe
// 	w.Close()
// 	syscall.Wait4(pid, nil, 0, nil)

// 	// Read the output of the command from the pipe and remove the "exit status 1" string
// 	var buf [1]byte
// 	var output string
// 	for {
// 		_, err := r.Read(buf[:])
// 		if err != nil {
// 			break
// 		}
// 		output += string(buf[0])
// 	}
// 	// output = strings.TrimSuffix(output, "\nexit status 1\n")
// 	output = strings.TrimSpace(output)

// 	// rough fix to no flags case, because new lines appears in the output
// 	// if len(os.Args) == 1 {
// 	// 	output = strings.ReplaceAll(output, "\n", "   ")
// 	// }
// 	// Print the output to the standard output
// 	fmt.Printf("%s\n", output)
// }
