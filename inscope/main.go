package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func scopeFile() string {

	pwd, err := filepath.Abs(".")
	if err != nil {
		os.Exit(0)
	}

	path := strings.Split(pwd, "/")

	for index := range path {
		root := strings.Join(path[:len(path)-index], "/") + "/.scope"

		_, err := os.Stat(root)
		if err == nil {
			return root
		} else {
			continue
		}
	}

	return ""
}

func runVim() {
	fmt.Println("add .scope file here!!!")
	fpath, _ := filepath.Abs(".")
	fpath += "/.scope"

	cmd := exec.Command("bash", "-c", "vim fuck")

	cmd.Start()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("could not run command :", err)
	}
}

func main() {

	// reading flags
	var current_sub bool
	var quite_mode bool
	var scope_file_path string
	flag.BoolVar(&current_sub, "cs", false, "current sub only")
	flag.BoolVar(&quite_mode, "q", false, "quite mode")
	flag.StringVar(&scope_file_path, "f", "", "path of scope file")

	flag.Parse()

	// findign .inscope file
	scope := scopeFile()
	if scope != "" {
		fmt.Println("file is: ", scope)
	} else {
		runVim()
	}

	if current_sub {
		fmt.Print("fuck\n")
	}

	// reading from standard input:
	// stdin := bufio.NewScanner(os.Stdin)
	// for stdin.Scan() {
	// 	fmt.Println(stdin.Text())
	// }
}
