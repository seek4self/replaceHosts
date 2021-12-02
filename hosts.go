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

func getHosts() []byte {
	resp, err := http.Get("https://gitee.com/ineo6/hosts/raw/master/hosts")
	if err != nil {
		fmt.Println("get github hosts err", err)
		return nil
	}
	defer resp.Body.Close()
	hosts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read github hosts body err", err)
		return nil
	}
	fmt.Println("get github hosts ok:")
	fmt.Println(string(hosts))
	fmt.Println()
	return hosts
}

func replaceHosts(github []byte) {
	if github == nil {
		return
	}
	hostsfile := unixHosts
	if runtime.GOOS == "windows" {
		hostsfile = winHosts
	}
	f, err := os.OpenFile(hostsfile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("open hosts file err", err)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	pos := 0
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("not found github hosts, writing ...")
				break
			}
			fmt.Println("read hosts err", err)
			return
		}
		if strings.Contains(line, githubStart) {
			fmt.Println("found github hosts, replacing ...")
			break
		}
		pos += len(line)
	}
	n, err := f.WriteAt(github, int64(pos))
	if n <= 0 {
		fmt.Println("write hosts err", err)
		return
	}
	fmt.Println("replace github hosts done.")
	flushDNS()
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
