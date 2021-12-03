package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
)

var (
	domain     = "github"
	githubOnce = false
	newhost    = "0"
	interval   = "2h"
	version    = ""
)

func init() {
	flag.StringVar(&domain, "d", domain, "domain in local hosts.")
	flag.StringVar(&newhost, "H", newhost, "the new host ip for the '-d'(input domain) flag.")
	flag.BoolVar(&githubOnce, "one", githubOnce, "replace github hosts once.")
	flag.StringVar(&interval, "i", interval, "replace interval. example: '1h30m', 'h' for hour, and 'm' for minute.")
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "replacer is a tool to replace the local hosts.")
		fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n", version)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	if domain != "github" {
		replaceDomain()
		return
	}

	replaceGithub()
	if githubOnce {
		return
	}
	// start cron job
	c := cron.New()
	c.AddFunc("@every "+interval, replaceGithub)
	c.Start()
	select {}
}
