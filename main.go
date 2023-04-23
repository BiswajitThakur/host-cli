package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var Version string = "1.0.0"
var PkgName string = "host-cli"
var homeDir, _ = os.UserHomeDir()

var hostPath string = HostPath()

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

var versionReg1 = regexp.MustCompile(`(?i)^\s*\-\-?v\s*$`)
var versionReg2 = regexp.MustCompile(`(?i)^\s*\-?\-?version\s*$`)
var validHost = regexp.MustCompile(`(?i)^[^\.][a-z\d\-]+\.[a-z\d\-\.]+[^\.]$`)
var help1Regx = regexp.MustCompile(`(?i)^\s*\-?\-?help\s*$`)
var help2Regx = regexp.MustCompile(`^\s*$`)
var blockReg = regexp.MustCompile(`(?i)^\s*\-?\-?block\s*$`)
var unBlockReg = regexp.MustCompile(`(?i)^\s*\-?\-?unblock\s*$`)
var addAllowListReg = regexp.MustCompile(`(?i)^\s*\-?\-?unblock\s*([^\s].+)$`)
var addBlockListReg = regexp.MustCompile(`(?i)^\s*\-?\-?block\s*([^\s].+)$`)
var upslReg1 = regexp.MustCompile(`(?i)^\s*\-?\-?updatesourcelist\s*$`)
var upslReg2 = regexp.MustCompile(`(?i)^\s*\-?\-?upsl\s*$`)

func main() {
	if help1Regx.MatchString(args) || help2Regx.MatchString(args) {
		fmt.Println(help())
	} else if versionReg1.MatchString(args) || versionReg2.MatchString(args) {
		fmt.Println(Version)
	} else if blockReg.MatchString(args) {
		os.Exit(block())
	} else if unBlockReg.MatchString(args) {
		os.Exit(unblock())
	} else if addAllowListReg.MatchString(args) {
		os.Exit(addAllowList())
	} else if addBlockListReg.MatchString(args) {
		os.Exit(addBlockList())
	} else if upslReg1.MatchString(args) || upslReg2.MatchString(args) {
		os.Exit(block())
	} else {
		fmt.Printf("Invalid Option : %s\n", args)
	}
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

func help() string {
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
	return fmt.Sprintf("%s\n%s\n", s1, s2)
}

func block() int {
	for _, i := range GetList(sourcePath) {
		c, err := GetContent(i)
		if err != nil {
			fmt.Println(err)
			return 1
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
	return 0
}

func unblock() int {
	fmt.Println("------------------Please wait------------------")
	WriteHosts(hostPath, NewSet())
	fmt.Println("------------------Success------------------")
	return 0
}

func addAllowList() int {
	fmt.Println("------------------Please wait------------------")
	s := addAllowListReg.FindStringSubmatch(args)[1]
	t := regexp.MustCompile(`\s+`).Split(s, -1)
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
	return 0
}

func addBlockList() int {
	fmt.Println("------------------Please wait------------------")
	s := addBlockListReg.FindStringSubmatch(args)[1]
	t := regexp.MustCompile(`\s+`).Split(s, -1)
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
	return 0
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
