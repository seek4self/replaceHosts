package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	if !disableDomain && newhost == "0" {
		fmt.Println("Non-github domain requires '-H' host parameter ")
		return
	}
	replaceHosts([]byte(fmt.Sprintf("%-18s %s\n", newhost, domain)))
}

func getHosts() []byte {
	if disableDomain {
		return []byte{}
	}
	resp, err := http.Get("https://gitee.com/ineo6/hosts/raw/master/hosts")
	if err != nil {
		log.Println("get github hosts from gitee err", err)
		log.Fatalln("please check gitee: https://gitee.com/ineo6/hosts")
		return nil
	}
	defer resp.Body.Close()
	hosts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("read github hosts body err", err)
		return nil
	}
	log.Println("get github hosts ok")
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
		log.Fatalln("open hosts file err", err)
		return
	}
	defer flushDNS()
	defer f.Close()
	r := bufio.NewReader(f)
	pos, action, change := 0, "replace", ""
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatalln("read hosts file err", err)
				return
			}
			if disableDomain {
				log.Println("disable github hosts done.")
				return
			}
			fmt.Printf("not found '%s' hosts, writing ...\n", domain)
			action = "write"
			break
		}

		if ok, host := findTarget(line); ok {
			if disableDomain {
				log.Printf("disable github host: %s", line)
				if n, err := f.WriteAt(commentDomain(line), int64(pos)); n <= 0 {
					log.Fatalln("disable github hosts err", err)
					return
				}
				pos += len(line)
				continue
			}
			log.Printf("found '%s' hosts, replacing ...\n", domain)
			if host != "" {
				change = fmt.Sprintf("'%s' -> '%s' ", host, newhost)
			}
			break
		}
		pos += len(line)
	}

	if n, err := f.WriteAt(hosts, int64(pos)); n <= 0 {
		log.Fatalln(action, "hosts err", err)
		return
	}
	log.Printf("%s '%s' hosts %sdone.\n", action, domain, change)
}

func commentDomain(line string) []byte {
	pos := strings.IndexByte(line, byte(' '))
	data := []byte(line)
	for i := pos + 1; i > 1; i-- {
		data[i] = data[i-2]
	}
	data[0] = byte('#')
	data[1] = byte(' ')
	return data
}

func findTarget(line string) (bool, string) {
	if disableDomain {
		if strings.HasPrefix(line, "#") {
			return false, ""
		}
		return strings.Contains(line, domain), ""
	}
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
		log.Fatalln("flush DNS cache err", err)
		return
	}
	log.Println("flush DNS cache done.")
}
