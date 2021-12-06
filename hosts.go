package main

import (
	"bufio"
	"fmt"
	"io"
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
	if runtime.GOOS == "windows" {
		return winHosts
	}
	return unixHosts
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
	hosts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("read github hosts body err", err)
		return nil
	}
	log.Println("get github hosts ok")
	// fmt.Println(string(hosts))
	// fmt.Println()
	return hosts
}

func exitOnWin(print bool) {
	if runtime.GOOS == "windows" && print {
		fmt.Printf("\nPress ENTER or Ctrl C to exit ...")
		b := make([]byte, 1)
		os.Stdin.Read(b)
		os.Exit(1)
	}
}

func replaceHosts(hosts []byte) {
	printWin := true
	defer func() {
		exitOnWin(printWin)
	}()
	if hosts == nil {
		return
	}

	f, err := os.OpenFile(hostsFile(), os.O_RDWR, 0666)
	if err != nil {
		log.Println("open hosts file err", err)
		return
	}
	defer flushDNS()
	defer f.Close()
	r := bufio.NewReader(f)
	pos, line, action, change := 0, "", "replace", ""
	for {
		pos += len(line)
		line, err = r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("read hosts file err", err)
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
					log.Println("disable github hosts err", err)
					return
				}
				continue
			}
			log.Printf("found '%s' hosts, replacing ...\n", domain)
			if host != "" {
				change = fmt.Sprintf("'%s' -> '%s' ", host, newhost)
			}
			break
		}
	}

	if n, err := f.WriteAt(hosts, int64(pos)); n <= 0 {
		log.Println(action, "hosts err", err)
		return
	}
	log.Printf("%s '%s' hosts %sdone.\n", action, domain, change)
	printWin = false
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
		if len(line) == 0 || line[0] == '#' {
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
	if err := exec.Command(exe, args...).Run(); err != nil {
		log.Fatalln("flush DNS cache err", err)
		return
	}
	log.Println("flush DNS cache done.")
}
