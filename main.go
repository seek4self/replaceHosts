package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/robfig/cron/v3"
)

var (
	domain        = "github"
	githubOnce    = false
	newhost       = "0"
	interval      = "2h"
	version       = ""
	disableDomain = false
	printVersion  = false
	daemon        = false
)

func init() {
	flag.StringVar(&domain, "D", domain, "domain in local hosts.")
	flag.StringVar(&newhost, "H", newhost, "the new host ip for the '-D'(input domain) flag.")
	flag.StringVar(&interval, "i", interval, "replace interval. example: '1h30m', 'h' for hour, and 'm' for minute.")
	flag.BoolVar(&githubOnce, "one", githubOnce, "replace github hosts once.")
	flag.BoolVar(&disableDomain, "dd", disableDomain, "disable domain hosts.")
	flag.BoolVar(&printVersion, "v", printVersion, "print version.")
	flag.BoolVar(&daemon, "d", daemon, "run app as a daemon.")
	flag.Usage = func() {
		echoVersion()
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if daemon {
		cmd := exec.Command(os.Args[0], flag.Args()...)
		if err := cmd.Start(); err != nil {
			fmt.Printf("start %s failed, error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
}

func echoVersion() {
	fmt.Fprintln(flag.CommandLine.Output(), "The tool to replace the local hosts.")
	fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n", version)
}

func main() {
	if printVersion {
		echoVersion()
		return
	}
	if domain != "github" {
		replaceDomain()
		return
	}

	replaceGithub()
	if githubOnce || disableDomain {
		return
	}
	// start cron job
	c := cron.New()
	c.AddFunc("@every "+interval, replaceGithub)
	c.Start()
	select {}
}
