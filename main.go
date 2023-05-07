package main

import (
	"bufio"
	"fmt"
	"host-cli/options"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var Version string = "1.0.5"
var PkgName string = "host-cli"
var homeDir, _ = os.UserHomeDir()

var hostPath string = HostPath()

// var hostPath string = "hosts"

var blockPath string = getPathIfExist("block.txt", []string{
	fmt.Sprintf("/usr/share/%s/block.txt", PkgName),
	fmt.Sprintf("/data/data/com.termux/files/usr/share/%s/block.txt", PkgName),
})
var allowPath string = getPathIfExist("allow.txt", []string{
	fmt.Sprintf("/usr/share/%s/allow.txt", PkgName),
	fmt.Sprintf("/data/data/com.termux/files/usr/share/%s/allow.txt", PkgName),
})
var sourcePath string = getSourcePathIfExist("sources.txt", []string{
	"https://adaway.org/hosts.txt",
	"https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&showintro=0&mimetype=plaintext",
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
}, []string{
	fmt.Sprintf("/usr/share/%s/sources.txt", PkgName),
	fmt.Sprintf("/data/data/com.termux/files/usr/share/%s/sources.txt", PkgName),
})

var blockList *Set = NewSet()
var args string = strings.Join(os.Args[1:], " ")

var validHost = regexp.MustCompile(`(?i)^[^\.][a-z\d\-]+\.[a-z\d\-\.]+[^\.]$`)

var Z = options.NewOptions(os.Args)

func main() {

	Z.Add(options.Option{
		SortName: "-h",
		LongName: "--help",
		Callback: help,
	})

	Z.Add(options.Option{
		SortName: "-v",
		LongName: "--version",
		Callback: func() { fmt.Println(Version) },
	})

	Z.Add(options.Option{
		SortName: "-b",
		LongName: "--block",
		Callback: func() {
			if len(Z.Args()) == 0 {
				block()
				return
			}
			addBlockList()
		},
	})

	Z.Add(options.Option{
		SortName: "-upsl",
		LongName: "--updatesourcelist",
		Callback: block,
	})

	Z.Add(options.Option{
		SortName: "-u",
		LongName: "--unblock",
		Callback: func() {
			if len(Z.Args()) == 0 {
				unblock()
				return
			}
			addAllowList()
		},
	})

	Z.Start()
}

func HostPath() string {
	win := `c:\Windows\System32\Drivers\etc\hosts`
	linux := "/etc/hosts"
	mac := "/private/etc/hosts"
	os := runtime.GOOS
	switch os {
	case "windows":
		return win
	case "darwin":
		return mac
	case "linux":
		return linux
	default:
		return linux
	}
}

func help() {
	s1 := fmt.Sprintf("\t\t%s Version %s\nUse: %s [OPTIONS] ...", PkgName, Version, os.Args[0])
	s2 := `    --help                                          This message
    --version                                       Print Version
    --block                                         Block ads
    --block host_name1 host_name2 ...               Add host_name1 host_name2 ... to block list.
    --unblock                                       Unblock ads and all host blocked by you
    --unblock                                       Unblock host_name1 host_name2  ...
    --updateSourceList | --upsl                     Update Ads Hostname List
Note : You can write options in any case.
`
	fmt.Printf("%s\n%s\n", s1, s2)
}

func block() {
	for _, i := range GetList(sourcePath) {
		c, err := GetContent(i)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, j := range FilterHosts(c) {
			blockList.Add(j)
		}
	}
	fmt.Println("------------------Please wait------------------")
	for _, i := range GetList(blockPath) {
		blockList.Add(i)
	}
	for _, i := range GetList(allowPath) {
		blockList.Remove(i)
	}
	WriteHosts(hostPath, blockList)
	fmt.Println("------------------Success------------------")
	fmt.Println("You may reboot your system to apply changes.")
}

func unblock() {
	fmt.Println("------------------Please wait------------------")
	WriteHosts(hostPath, NewSet())
	fmt.Println("------------------Success------------------")
}

func addAllowList() {
	fmt.Println("------------------Please wait------------------")
	t := Z.Args()
	al := NewSet()
	bll := NewSet()
	bl := ReadLocalHosts(hostPath)
	for _, i := range GetList(blockPath) {
		if validHost.MatchString(i) {
			bll.Add(i)
		} else {
			fmt.Printf("Invalid Hostname : %s\n", i)
		}
	}
	for _, i := range t {
		if validHost.MatchString(i) {
			al.Add(i)
			bll.Remove(i)
		} else {
			fmt.Printf("Invalid Hostname : %s\n", i)
		}
	}
	for _, i := range GetList(allowPath) {
		if validHost.MatchString(i) {
			al.Add(i)
			bll.Remove(i)
		} else {
			fmt.Printf("Invalid Hostname : %s\n", i)
		}
	}
	WriteList(allowPath, al)
	WriteList(blockPath, bll)
	for _, i := range bl {
		blockList.Add(i)
	}
	for _, i := range al.GetAll() {
		blockList.Remove(i)
	}
	WriteHosts(hostPath, blockList)
	fmt.Println("------------------Success------------------")
}

func addBlockList() {
	fmt.Println("------------------Please wait------------------")
	t := Z.Args()
	bl := ReadLocalHosts(hostPath)
	al := NewSet()
	bll := NewSet()
	for _, i := range GetList(allowPath) {
		if validHost.MatchString(i) {
			al.Add(i)
		} else {
			fmt.Printf("Invalid Hostname : %s\n", i)
		}
	}
	for _, i := range GetList(blockPath) {
		if validHost.MatchString(i) {
			bll.Add(i)
			al.Remove(i)
		} else {
			fmt.Printf("Invalid Hostname : %s\n", i)
		}
	}
	for _, i := range t {
		if validHost.MatchString(i) {
			bll.Add(i)
			al.Remove(i)
		} else {
			fmt.Printf("Invalid Hostname : %s\n", i)
		}
	}
	WriteList(allowPath, al)
	WriteList(blockPath, bll)
	for _, i := range bl {
		bll.Add(i)
	}
	for _, i := range al.GetAll() {
		bll.Remove(i)
	}
	WriteHosts(hostPath, bll)
	fmt.Println("------------------Success------------------")
}

func creteEmptyFile(f string) {
	if !IsExist(f) {
		g, err := os.Create(f)
		defer g.Close()
		fmt.Println(err)
	}
}

func getPathIfExist(p string, paths []string) string {
	for _, i := range paths {
		if IsExist(i) {
			return i
		}
	}
	pth := filepath.Join(homeDir, PkgName)
	g := filepath.Join(pth, p)
	if IsExist(g) {
		return g
	}
	_, err := os.Stat(pth)
	if os.IsNotExist(err) {
		err1 := os.MkdirAll(pth, 0700)
		if err1 != nil {
			fmt.Println(err1)
		}
	}
	creteEmptyFile(g)
	return g
}

func getSourcePathIfExist(p string, urls []string, paths []string) string {
	for _, i := range paths {
		if IsExist(i) {
			return i
		}
	}
	pth := filepath.Join(homeDir, PkgName)
	g := filepath.Join(pth, p)
	if IsExist(g) {
		return g
	}
	_, err := os.Stat(pth)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(pth, 0700)
	}
	f, _ := os.Create(g)
	defer f.Close()
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(strings.Join(urls, "\n"))
	if err != nil {
		fmt.Println(err)
	}
	writer.Flush()
	return g
}
