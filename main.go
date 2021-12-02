package main

import (
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	c.AddFunc("@every 2h", func() { replaceHosts(getHosts()) })
	c.Start()
}
