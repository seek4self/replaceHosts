package main

import (
	"flag"
	"fmt"

	"github.com/robfig/cron/v3"
)

var (
	domain     = "github"
	githubOnce = false
	newhost    = "0"
	interval   = "2h"
)

func main() {
	flag.StringVar(&domain, "d", domain, "domain in hosts")
	flag.StringVar(&newhost, "H", newhost, "domain host")
	flag.BoolVar(&githubOnce, "one", githubOnce, "replace github hosts once")
	flag.StringVar(&interval, "i", interval, "replace interval, 'h' is hour, 'm' is minute")
	flag.Parse()

	if githubOnce {
		replaceHosts(getHosts())
		return
	}

	if domain != "github" {
		if newhost == "0" {
			fmt.Println("Non-github domain requires '-H' host parameter ")
			return
		}
		replaceHosts([]byte(newhost + " " + domain + "\n"))
		return
	}

	c := cron.New()
	c.AddFunc("@every "+interval, func() { replaceHosts(getHosts()) })
	c.Start()
	select {}
}
