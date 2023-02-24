package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sync"
)

func worker(domain string, results chan result, output string) {

	var item result
	outdir := output + "/" + domain
	os.MkdirAll(outdir, os.ModePerm)
	item.domain = domain
	item.assetfinder = runCommand("assetfinder -subs-only " + domain + " | anew > " + outdir + "/assetfinder")
	item.subfinder = runCommand("subfinder -d " + domain + " -o " + outdir + "/subfinder")

	results <- item

}

func runCommand(command string) bool {
	// fmt.Println(command)
	// com := exec.Command(command, args...)
	com := exec.Command("bash", "-c", command)
	com.Stdout = os.Stdout
	com.Stderr = os.Stderr
	if err := com.Run(); err != nil {
		return false
	}

	return true

}

type empty struct{}

type result struct {
	domain      string
	assetfinder bool
	subfinder   bool
	amass       bool
}

var wg sync.WaitGroup

func main() {

	var (
		output = flag.String("o", "/home/dav00d/BugBounty/programs/output", "output directory")
		_      = flag.String("w", "~/BugBounty/wordlist/sort_subs12.txt", "wordlist path")
		_      = flag.String("r", "~/BugBounty/wordlist/resolvers.txt", "resolver path")
	)
	var gather []result
	var counter int = 0
	flag.Parse()

	results := make(chan result)
	// domains := make(chan string)

	// go func() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		r, _ := regexp.Compile("^\\*\\.")

		if r.MatchString(sc.Text()) {
			// domains <- r.ReplaceAllString(sc.Text(), "")
			domain := r.ReplaceAllString(sc.Text(), "")
			counter++
			go worker(domain, results, *output)
		}

	}

	for i := 0; i < counter; i++ {
		gather = append(gather, <-results)
	}

	for i, v := range gather {
		fmt.Println(i+1, ":", v)
	}
	// }()

}
