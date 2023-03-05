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

	"github.com/jedib0t/go-pretty/v6/progress"
)

func worker(domain string, commands map[int]string, wg *sync.WaitGroup, signal chan<- int, gather map[string][]int, pw progress.Writer) {

	var item result

	tracker := progress.Tracker{Message: domain, Total: 12, Units: progress.UnitsDefault}
	tracker.Reset()
	pw.AppendTracker(&tracker)
	//mkdir folder for each domain
	outdir := output + "/" + domain
	os.MkdirAll(outdir, os.ModePerm)

	item.domain = domain

	for i := 0; i < len(commands); i++ {
		cmd := fmt.Sprintf(commands[i], domain)
		item.runCommand(cmd)
		gather[domain] = append(gather[domain], 1)
		// signal <- 1
		time.Sleep(time.Millisecond * 150)
		tracker.Increment(1)
	}
	time.Sleep(time.Second * 2)
	wg.Done()
}

func (i *result) runCommand(command string) {
	com := exec.Command("bash", "-c", command)
	if err := com.Run(); err != nil {
		fmt.Println("fuck you have an error")
		os.Exit(1)
	}

}

type result struct {
	domain string
}

func initialCommands(outdir string, wordlist string) map[int]string {
	commands := map[int]string{
		0:  "assetfinder -subs-only  %[1]s | anew > " + outdir + "/%[1]s" + "/assetfinder",
		1:  "subfinder -d %[1]s -o " + outdir + "/%[1]s" + "/subfinder",
		2:  "amass enum -passive -d %[1]s > " + outdir + "/%[1]s" + "/amass",
		3:  "cat " + outdir + "/%[1]s" + "/assetfinder " + outdir + "/%[1]s" + "/subfinder " + outdir + "/%[1]s" + "/amass | deduplicate --sort > " + outdir + "/%[1]s" + "/round1",
		4:  "rm -f " + outdir + "/%[1]s" + "/assetfinder " + outdir + "/%[1]s" + "/subfinder " + outdir + "/%[1]s" + "/amass 2>/dev/null",
		5:  "cp " + wordlist + " " + outdir + "/%[1]s" + "/shuffle",
		6:  "sed -e \"s/$/.${%[1]s##*\\/}/\"  -i " + outdir + "/%[1]s" + "/shuffle",
		7:  "dnsx -list " + outdir + "/%[1]s" + "/shuffle -silent -o " + outdir + "/%[1]s" + "/step1",
		8:  "cat " + outdir + "/%[1]s" + "/step1 | anew -q " + outdir + "/%[1]s" + "/round1",
		9:  "gotator -silent -sub " + outdir + "/%[1]s" + "/round1 -depth 2 -mindup > " + outdir + "/%[1]s" + "/gotator",
		10: "dnsx -list " + outdir + "/%[1]s" + "/gotator -r " + resolver + " -silent -o " + outdir + "/%[1]s" + "/step2",
		11: "cat " + outdir + "/%[1]s" + "/step1 " + outdir + "/%[1]s" + "/step2 | deduplicate --sort > " + outdir + "/%[1]s" + "/final",
	}

	return commands
}

var wordlist string
var resolver string
var output string

func main() {

	gather := make(map[string][]int)

	var wg sync.WaitGroup

	flag.StringVar(&output, "o", "/home/dav00d/BugBounty/programs/output", "output directory")
	flag.StringVar(&wordlist, "w", "/home/dav00d/BugBounty/wordlist/sort_subs12.txt", "wordlist path")
	flag.StringVar(&resolver, "r", "/home/dav00d/BugBounty/wordlist/resolvers.txt", "resolver path")

	flag.Parse()

	pw := progress.NewWriter()
	pw.SetOutputWriter(os.Stdout)
	pw.SetAutoStop(true)

	signal := make(chan int, 1)

	commands := initialCommands(output, wordlist)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		r, _ := regexp.Compile(`^\*\.`)

		if r.MatchString(sc.Text()) {
			domain := r.ReplaceAllString(sc.Text(), "")
			wg.Add(1)
			go worker(domain, commands, &wg, signal, gather, pw)
		}

	}
	time.Sleep(time.Millisecond * 200)
	pw.Render()
	wg.Wait()

}
