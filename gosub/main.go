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

	// Round 1:

	// item.runCommand("assetfinder -subs-only "+domain+" | anew > "+outdir+"/assetfinder", "assetfinder", signal)
	// item.runCommand("subfinder -d "+domain+" -o "+outdir+"/subfinder", "subfinder", signal)
	// item.runCommand("amass enum -passive -d "+domain+" > "+outdir+"/amass", "amass", signal)

	// item.runCommand("cat "+outdir+"/assetfinder "+outdir+"/subfinder "+outdir+"/amass | deduplicate --sort > "+outdir+"/round1", "", signal)
	// item.runCommand("rm -f "+outdir+"/assetfinder "+outdir+"/subfinder "+outdir+"/amass 2>/dev/null", "", signal)

	// time.Sleep(time.Second * 2)

	// Round 2

	// runCommand("cp " + wordlist + " " + outdir + "/shuffle")

	// runCommand("sed -e \"s/$/.${" + domain + "##*\\/}/\"  -i " + outdir + "/shuffle")

	// runCommand("dnsx -list " + outdir + "/shuffle -r " + resolver + "-silent -o " + outdir + "/step1")

	// runCommand("cat " + outdir + "/step1 | anew -q " + outdir + "/round1")

	// runCommand("gotator -silent -sub " + outdir + "/round1 -depth 2 -mindup > " + outdir + "/gotator")

	// runCommand("dnsx -list " + outdir + "/gotator -r " + resolver + " -silent -o " + outdir + "/step2")

	// runCommand("cat " + outdir + "/step1 " + outdir + "/step2 | deduplicate --sort > " + outdir + "/final")

	// results <- item

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

func initialCommands(outdir string) map[int]string {
	commands := map[int]string{
		0: "assetfinder -subs-only  %[1]s | anew > " + outdir + "/%[1]s" + "/assetfinder",
		1: "subfinder -d %[1]s -o " + outdir + "/%[1]s" + "/subfinder",
		2: "amass enum -passive -d %[1]s > " + outdir + "/%[1]s" + "/amass",
		3: "cat " + outdir + "/%[1]s" + "/assetfinder " + outdir + "/%[1]s" + "/subfinder " + outdir + "/%[1]s" + "/amass | deduplicate --sort > " + outdir + "/%[1]s" + "/round1",
		4: "rm -f " + outdir + "/%[1]s" + "/assetfinder " + outdir + "/%[1]s" + "/subfinder " + outdir + "/%[1]s" + "/amass 2>/dev/null",
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

	commands := initialCommands(output)

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
