package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"time"
)

func worker(domain string, results chan result, wg *sync.WaitGroup) {
	wg.Add(1)

	var item result
	outdir := output + "/" + domain
	os.MkdirAll(outdir, os.ModePerm)
	item.domain = domain

	// Round 1:

	item.assetfinder = runCommand("assetfinder -subs-only " + domain + " | anew > " + outdir + "/assetfinder")
	item.subfinder = runCommand("subfinder -d " + domain + " -o " + outdir + "/subfinder")
	item.amass = runCommand("amass enum -passive -d " + domain + " > " + outdir + "/amass")

	results <- item

	runCommand("cat " + outdir + "/assetfinder " + outdir + "/subfinder " + outdir + "/amass | deduplicate --sort > " + outdir + "/round1")
	runCommand("rm -f " + outdir + "/assetfinder " + outdir + "/subfinder " + outdir + "/amass 2>/dev/null")

	time.Sleep(time.Second * 2)

	// Round 2

	runCommand("cp " + wordlist + " " + outdir + "/shuffle")

	runCommand("sed -e \"s/$/.${" + domain + "##*\\/}/\"  -i " + outdir + "/shuffle")
	runCommand("dnsx -list " + outdir + "/shuffle -r " + resolver + "-silent -o " + outdir + "/step1")

	runCommand("cat " + outdir + "/step1 | anew -q " + outdir + "/round1")

	runCommand("gotator -silent -sub " + outdir + "/round1 -depth 2 -mindup > " + outdir + "/gotator")

	runCommand("dnsx -list " + outdir + "/gotator -r " + resolver + " -silent -o " + outdir + "/step2")

	final := runCommand("cat " + outdir + "/step1 " + outdir + "/step2 | deduplicate --sort > " + outdir + "/final")

	if final {
		fmt.Println("final done!")
	}

	time.Sleep(time.Second * 2)
	wg.Done()
}

func runCommand(command string) bool {
	com := exec.Command("bash", "-c", command)
	// com.Stdout = os.Stdout
	// com.Stderr = os.Stderr
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

var wordlist string
var resolver string
var output string

func main() {

	var wg sync.WaitGroup

	flag.StringVar(&output, "o", "/home/dav00d/BugBounty/programs/output", "output directory")
	flag.StringVar(&wordlist, "w", "/home/dav00d/BugBounty/wordlist/sort_subs12.txt", "wordlist path")

	flag.StringVar(&resolver, "r", "/home/dav00d/BugBounty/wordlist/resolvers.txt", "resolver path")

	var gather []result
	var counter int = 0
	flag.Parse()

	results := make(chan result)
	// domains := make(chan string)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		r, _ := regexp.Compile("^\\*\\.")

		if r.MatchString(sc.Text()) {
			// domains <- r.ReplaceAllString(sc.Text(), "")
			domain := r.ReplaceAllString(sc.Text(), "")
			counter++
			go worker(domain, results, &wg)
		}

	}

	for i := 0; i < counter; i++ {
		gather = append(gather, <-results)
	}

	for i, v := range gather {
		fmt.Println(i+1, ":", v)
	}

	wg.Wait()

}
