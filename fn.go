package main

import (
    "fmt"
    "os"
    "strings"
    "net/http"
    "io/ioutil"
    "regexp"
)

//var reg = 

func GetContent(url string) (string, error) {
    fmt.Printf("Downloading... --> %s\n", url)
    response, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

func ReadFile(filename string) (string, error) {
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return "", err
    }
    return string(content), nil
}

func WriteFile(path string, d string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
	    return err
	}
	defer file.Close()
	data := []byte(d)
	_, err = file.Write(data)
	if err != nil {
	    return err
	}
	return nil
}

func IsComment(s string) bool {
	r1 := regexp.MustCompile(`^\s*#+.*`)
	r2 := regexp.MustCompile(`^\s*$`)
	return r1.MatchString(s) || r2.MatchString(s)
}

func FindStr(s []string, v string) int {
	for i, j := range s {
		if j == v {
			return i
		}
	}
	return -1
}

// This function filter hostname from string.
func FilterHosts(s string) []string {
    var v []string = []string{}
	r := regexp.MustCompile(`^\s*\d\d?\d?\.\d\d?\d?\.\d\d?\d?\.\d\d?\d?\s+([A-Za-z0-9-\.]+)\s*`)
	for _, i := range strings.Split(s, "\n") {
		if !IsComment(i) && r.MatchString(i) {
			v = append(v, r.FindStringSubmatch(i)[1])
		}
	}
	return v
}

func ReadLocalHosts(path string) []string {
	r, err := ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	reg := regexp.MustCompile(`#\s*BT\-start\s*#?([\d\D])*#\s*BT\-end\s*#?`)
	return FilterHosts(reg.FindString(r))
}

func WriteHosts(path string, s *Set) {
	v, err := ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	r := regexp.MustCompile(`#\s*BT\-start\s*#?[\d\D]*#\s*BT\-end\s*#?`)
	if r.MatchString(v) {
		err = WriteFile(path, r.ReplaceAllString(v, Block(s)))
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	} else {
		err = WriteFile(path, fmt.Sprintf("%s\n%s", v, Block(s)))
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}

func WriteList(path string, s *Set) {
	err := WriteFile(path, strings.Join(s.GetAll(), "\n"))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// get (source | allow | block) list from path name
func GetList(path string) []string {
	d, err := ReadFile(path)
	if err != nil {
		fmt.Printf("%s : Not Found\n", path)
		os.Exit(-1)
	}
	var v []string = []string{}
	for _, i := range strings.Split(d, "\n") {
		if !IsComment(i) {
			v = append(v, i)
		}
	}
	return v
}

func Block(s *Set) string {
	k := s.GetAll()
	if len(k) == 0 {
		return "# BT-start\n# BT-end\n"
	}
	m := ""
	for _, i := range k {
		m = fmt.Sprintf("%s0.0.0.0    %s\n", m, i)
	}
	return fmt.Sprintf("# BT-start\n%s# BT-end\n", m)
}

