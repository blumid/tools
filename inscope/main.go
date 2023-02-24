package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	// "os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type scopeChecker struct {
	patterns []*regexp.Regexp
	// antipatterns []*regexp.Regexp
}

func newChecker(r io.Reader) (*scopeChecker, error) {
	line := bufio.NewScanner(r)
	s := &scopeChecker{
		patterns: make([]*regexp.Regexp, 0),
	}
	for line.Scan() {
		p := strings.TrimSpace(line.Text())
		if p == "" {
			continue
		}

		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		s.patterns = append(s.patterns, re)

	}
	return s, nil
}

func scopeFile() (io.ReadCloser, error) {

	pwd, err := filepath.Abs(".")
	if err != nil {
		os.Exit(0)
	}

	path := strings.Split(pwd, "/")

	for index := range path {
		fpath := strings.Join(path[:len(path)-index], "/") + "/.scope"

		_, err := os.Stat(fpath)
		if err == nil {
			f, _ := os.Open(fpath)
			return f, nil
		} else {
			continue
		}
	}

	return nil, errors.New("unable to find .scope file in current directory or any parent directory")
}

/*
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
*/

func (s *scopeChecker) inScope(domain string) bool {

	for _, p := range s.patterns {
		if p.MatchString(domain) {
			return true
		} else {
			return false
		}
	}
	fmt.Println(domain)
	return false
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

	// findign .scope file
	sf, err := scopeFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening scope file: %s\n", err)
		return
	}
	ch, _ := newChecker(sf)
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		domain := strings.TrimSpace(stdin.Text())
		if ch.inScope(domain) {
			fmt.Println(domain)
		}
	}

	if current_sub {
		fmt.Print("fuck\n")
	}

}
