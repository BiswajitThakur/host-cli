package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
)

var version string = "0.0.1-06042023"

// var hostPath string = "test_hosts"
var hostPath string = HostPath()
var blockPath string = ".block.txt"
var allowPath string = ".allow.txt"
var sourcePath string = ".sources.txt"

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
var updateReg = regexp.MustCompile(`(?i)^\s*\-?\-?update\s*$`)

func main() {
	if help1Regx.MatchString(args) || help2Regx.MatchString(args) {
		fmt.Println(help())
	} else if versionReg1.MatchString(args) || versionReg2.MatchString(args) {
		fmt.Println(version)
	} else if blockReg.MatchString(args) {
		block()
	} else if unBlockReg.MatchString(args) {
		unblock()
	} else if addAllowListReg.MatchString(args) {
		addAllowList()
	} else if addBlockListReg.MatchString(args) {
		addBlockList()
	} else if updateReg.MatchString(args) {
		block()
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
	return `Commands :
    --help                                          This message
    --version                                       Print Version
    --block                                         Block ads
    --block host_name1 host_name2 ...               Add host_name1 host_name2 ... to block list.
    --unblock                                       Unblock ads and all host blocked by you
    --unblock host_name1 host_name2 ...             Unblock host_name1 host_name2  ...
    --update                                        Update Ads Hostname List
`
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
}

func addBlockList() {
	fmt.Println("------------------Please wait------------------")
	s := addBlockListReg.FindStringSubmatch(args)[1]
	t := regexp.MustCompile(`\s+`).Split(s, -1)
	bl := ReadLocalHosts(hostPath)
	al := NewSet()
	bll := NewSet()
	// fmt.Println(8888)
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
