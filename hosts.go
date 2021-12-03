package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	winHosts  = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	unixHosts = "/etc/hosts"

	githubStart = "# GitHub Host Start"
)

func hostsFile() string {
	hostsfile := unixHosts
	if runtime.GOOS == "windows" {
		hostsfile = winHosts
	}
	return hostsfile
}

func replaceGithub() {
	replaceHosts(getHosts())
}

func replaceDomain() {
	if newhost == "0" {
		fmt.Println("Non-github domain requires '-H' host parameter ")
		return
	}
	replaceHosts([]byte(newhost + " " + domain + "\n"))
}

func getHosts() []byte {
	resp, err := http.Get("https://gitee.com/ineo6/hosts/raw/master/hosts")
	if err != nil {
		fmt.Println("get github hosts from gitee err", err)
		fmt.Println("please check gitee: https://gitee.com/ineo6/hosts")
		return nil
	}
	defer resp.Body.Close()
	hosts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read github hosts body err", err)
		return nil
	}
	fmt.Println("get github hosts ok")
	// fmt.Println(string(hosts))
	// fmt.Println()
	return hosts
}

func replaceHosts(hosts []byte) {
	if hosts == nil {
		return
	}
	f, err := os.OpenFile(hostsFile(), os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("open hosts file err", err)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	pos, action, change := 0, "replace", ""
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("not found '%s' hosts, writing ...\n", domain)
				action = "write"
				break
			}
			fmt.Println("read hosts file err", err)
			return
		}
		if ok, host := findTarget(line); ok {
			fmt.Printf("found '%s' hosts, replacing ...\n", domain)
			if host != "" {
				change = fmt.Sprintf("'%s' -> '%s' ", host, newhost)
			}
			break
		}
		pos += len(line)
	}
	n, err := f.WriteAt(hosts, int64(pos))
	if n <= 0 {
		fmt.Println("write hosts err", err)
		return
	}
	fmt.Printf("%s '%s' hosts %sdone.\n", action, domain, change)
	flushDNS()
}

func findTarget(line string) (bool, string) {
	target, host := githubStart, ""
	if domain != "github" {
		target = domain
		host = strings.Split(line, " ")[0]
	}
	return strings.Contains(line, target), host
}

func flushDNS() {
	exe, args := "ls", []string{}
	switch runtime.GOOS {
	case "windows":
		exe = "ipconfig"
		args = []string{"/flushdns"}
	case "linux":
		exe = "/etc/init.d/networking"
		args = []string{"restart"}
	case "darwin":
		exe = "killall"
		args = []string{"-HUP", "mDNSResponder"}
	}
	cmd := exec.Command(exe, args...)
	if err := cmd.Run(); err != nil {
		fmt.Println("flush dns cache err", err)
		return
	}
	fmt.Println("flush dns cache done.")
}
