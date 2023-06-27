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

var Version string = "1.1.0"
var PkgName string = "host-cli"
var HomeDir, _ = os.UserHomeDir()

var HostPath string = HostPath1()

// var HostPath string = "hosts"
// var BlockPath string = "block.txt"
// var AllowPath string = "allow.txt"
// var SourcePath string = "sources.txt"
// var RedirectPath string = "redirect.txt"

var BlockPath string = getPathIfExist("block.txt", []string{
	fmt.Sprintf("/usr/share/%s/block.txt", PkgName),
	fmt.Sprintf("/data/data/com.termux/files/usr/share/%s/block.txt", PkgName),
})
var AllowPath string = getPathIfExist("allow.txt", []string{
	fmt.Sprintf("/usr/share/%s/allow.txt", PkgName),
	fmt.Sprintf("/data/data/com.termux/files/usr/share/%s/allow.txt", PkgName),
})
var SourcePath string = getSourcePathIfExist("sources.txt", []string{
	"https://adaway.org/hosts.txt",
	"https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&showintro=0&mimetype=plaintext",
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
}, []string{
	fmt.Sprintf("/usr/share/%s/sources.txt", PkgName),
	fmt.Sprintf("/data/data/com.termux/files/usr/share/%s/sources.txt", PkgName),
})

var RedirectPath string = getPathIfExist("redirect.txt", []string{
	fmt.Sprintf("/usr/share/%s/redirect.txt", PkgName),
	fmt.Sprintf("/data/data/com.termux/files/usr/share/%s/redirect.txt", PkgName),
})

var blockList *Set = NewSet()
var args string = strings.Join(os.Args[1:], " ")
var waitStr string = "------------- Please Wait ------------"
var succStr string = "------------- Success ------------"
var validHost = regexp.MustCompile(`(?i)^[^\.][a-z\d\-]+\.[a-z\d\-\.]+[^\.]$`)

var Z = options.NewOptions(os.Args)

func main() {

	Z.Add(options.Option{
		SortName: "-h",
		LongName: "--help",
		Callback: Help,
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
			fmt.Println(waitStr)
			v := NewSet()
			if len(Z.Args()) == 0 {
				errs := OptionBlock(v, true)
				if len(errs) == 0 {
					fmt.Println(succStr)
					return
				} else {
					for _, i := range errs {
						fmt.Println(i)
					}
					return
				}
			}
			for _, i := range Z.Args() {
				v.Add(i)
			}
			errs := OptionBlock(v, false)
			if len(errs) == 0 {
				fmt.Println(succStr)
				return
			}
			for _, i := range errs {
				fmt.Println(i)
			}
		},
	})

	Z.Add(options.Option{
		SortName: "-upsl",
		LongName: "--updatesourcelist",
		Callback: func() {
			errs := OptionBlock(NewSet(), true)
			for _, i := range errs {
				fmt.Println(i)
			}
		},
	})

	Z.Add(options.Option{
		SortName: "-ub",
		LongName: "--unblock",
		Callback: func() {
			fmt.Println(waitStr)
			v := NewSet()
			for _, i := range Z.Args() {
				v.Add(i)
			}
			errs := OptionAllowFn(v)
			if len(errs) == 0 {
				fmt.Println(succStr)
				return
			}
			for _, i := range errs {
				fmt.Println(i)
			}
		},
	})

	Z.Add(options.Option{
		SortName: "-http",
		LongName: "--http",
		Callback: func() { Server(Z.Args()[0]) },
	})

	Z.Start()
}

func HostPath1() string {
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

func Help() {
	s1 := fmt.Sprintf("\t\t%s Version %s\nUse: %s [OPTIONS] ...", PkgName, Version, os.Args[0])
	s2 := `    --help                                          This message
    --version                                       Print Version
    --block                                         Block ads
    --block host_name1 host_name2 ...               Add host_name1 host_name2 ... to block list.
    --unblock                                       Unblock ads and all host blocked by you
    --unblock                                       Unblock host_name1 host_name2  ...
    --updateSourceList | --upsl                     Update Ads Hostname List
	--http 3999 or --http 127.0.0.1:3999            http server start at port 3999
Note : You can write options in any case.
`
	fmt.Printf("%s\n%s\n", s1, s2)
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
	pth := filepath.Join(HomeDir, PkgName)
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
	pth := filepath.Join(HomeDir, PkgName)
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
