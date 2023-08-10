package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var R1 = regexp.MustCompile(`#\s*BT\-start\s*#*([\d\D])*#\s*BT\-end\s*#*`)
var R2 = regexp.MustCompile(`#\s*BT\-redirect\-start\s*#*[\d\D]*#\s*BT\-redirect\-end\s*#*`)

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

func IsExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
func IsExistOneOfMany(paths ...string) bool {
	for _, i := range paths {
		if IsExist(i) {
			return true
		}
	}
	return false
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
	rr := regexp.MustCompile(`(?i)^\s*localhost\s+([A-Za-z0-9-\.]+)\s*`)
	for _, i := range strings.Split(s, "\n") {
		if !IsComment(i) && r.MatchString(i) {
			v = append(v, r.FindStringSubmatch(i)[1])
		} else if !IsComment(i) && rr.MatchString(i) {
			v = append(v, rr.FindStringSubmatch(i)[1])
		}
	}
	return v
}

func FilterRedirectHosts(s string) [][2]string {
	var v [][2]string = [][2]string{}
	r := regexp.MustCompile(`^\s*([^\.][A-Za-z0-9\.\-]+)\s+([^\.][A-Za-z0-9\.\-]+)\s*$`)
	for _, i := range strings.Split(s, "\n") {
		if !IsComment(i) && r.MatchString(i) {
			j := r.FindStringSubmatch(i)
			v = append(v, [2]string{j[1], j[2]})
		}
	}
	return v
}

func ReadLocalHosts(pth string) []string {
	r, err := ReadFile(pth)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return FilterHosts(R1.FindString(r))
}

func ReadLocalRedirectHosts(pth string) ([][2]string, error) {
	r, err := ReadFile(pth)
	if err != nil {
		return nil, err
	}
	return FilterRedirectHosts(R2.FindString(r)), nil
}

func WriteHosts(pth string, s *Set) error {
	v, err := ReadFile(pth)
	if err != nil {
		return err
	}
	if R2.MatchString(v) {
		v = R2.ReplaceAllString(v, BuildRedirectList(s))
	} else {
		v = fmt.Sprintf("%s\n%s", v, BuildRedirectList(s))
	}
	if R1.MatchString(v) {
		err = WriteFile(pth, R1.ReplaceAllString(v, BuildBlockList(s)))
		if err != nil {
			return err
		}
	} else {
		err = WriteFile(pth, fmt.Sprintf("%s\n%s", v, BuildBlockList(s)))
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteList(pth string, s *Set) error {
	err := WriteFile(pth, strings.Join(s.GetAll(), "\n"))
	if err != nil {
		return err
	}
	return nil
}

func WriteRedirectList(pth string, s *Set) error {
	s1 := []string{}
	// localRedirectList := NewSet()
	// rdList, err := GetRedirectList(RedirectPath)
	// if err != nil {
	// 	return err
	// }
	// for _, i := range rdList {
	// 	localRedirectList.AddRedirect(i[0], i[1])
	// }
	// for _, i := range s.GetAllRedirect() {
	// 	localRedirectList.AddRedirect(i[0], i[1])
	// }
	for _, i := range s.GetAllRedirect() {
		s1 = append(s1, fmt.Sprintf("%s %s", i[0], i[1]))
	}
	err := WriteFile(pth, strings.Join(s1, "\n"))
	if err != nil {
		return err
	}
	return nil
}

func WriteRedirectListHosts(pth string, s *Set) error {
	v, err := ReadFile(pth)
	if err != nil {
		return err
	}

	rdList, err := GetRedirectList(RedirectPath)
	if err != nil {
		return err
	}
	for _, i := range rdList {
		s.AddRedirect(i[0], i[1])
	}

	if R2.MatchString(v) {
		v = R2.ReplaceAllString(v, BuildRedirectList(s))
	} else {
		v = fmt.Sprintf("%s\n%s", v, BuildRedirectList(s))
	}
	err = WriteFile(pth, v)
	return err
}

// get (source | allow | block) list from path name.
// return type list, comment, error
func GetList(path string) ([]string, []string, error) {
	d, err := ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var lst []string = []string{}
	var cmt []string = []string{}
	for _, i := range strings.Split(d, "\n") {
		if !IsComment(i) {
			lst = append(lst, i)
		} else {
			if !regexp.MustCompile(`^\s*$`).MatchString(i) {
				cmt = append(cmt, i)
			}
		}
	}
	return lst, cmt, nil
}

func GetRedirectList(pth string) ([][2]string, error) {
	v, err := ReadFile(pth)
	if err != nil {
		return nil, err
	}
	return FilterRedirectHosts(v), nil
}

func BuildBlockList(s *Set) string {
	st := [2]string{"# BT-start #", "# BT-end #"}
	k := s.GetAll()
	if len(k) == 0 {
		return fmt.Sprintf("%s\n%s", st[0], st[1])
	}
	l := make([]string, len(k)+2)
	l[0] = st[0]
	var j int = 1
	for _, i := range k {
		l[j] = fmt.Sprintf("0.0.0.0 %s", i)
		j++
	}
	l[j] = st[1]
	return strings.Join(l, "\n")
}

func BuildRedirectList(s *Set) string {
	st := [2]string{"# BT-redirect-start #", "# BT-redirect-end #"}
	k := s.GetAllRedirect()
	if len(k) == 0 {
		return fmt.Sprintf("%s\n%s", st[0], st[1])
	}
	l := make([]string, len(k)+2)
	l[0] = st[0]
	var j int = 1
	for _, i := range k {
		l[j] = fmt.Sprintf("%s %s", i[0], i[1])
		j++
	}
	l[j] = st[1]
	return strings.Join(l, "\n")
}

func DownloadHostsLists(sourcePth string) ([]string, []error) {
	v := []string{}
	e := []error{}
	q, _, err := GetList(sourcePth)
	if err != nil {
		e = append(e, err)
	}
	for _, i := range q {
		c, err := GetContent(i)
		if err != nil {
			e = append(e, err)
		} else {
			for _, j := range FilterHosts(c) {
				v = append(v, j)
			}
		}
	}
	return v, e
}

func OptionBlock(bb *Set, isDownloadSources bool) []error {
	errs := []error{}
	localBlockList := NewSet()
	localAllowList := NewSet()
	localRedirectList := NewSet()
	q, _, err := GetList(BlockPath)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, i := range q {
			localBlockList.Add(i)
		}
	}
	q, _, err = GetList(AllowPath)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, i := range q {
			localAllowList.Add(i)
		}
	}
	rdList, err := GetRedirectList(RedirectPath)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, i := range rdList {
			localBlockList.Remove(i[1])
			localRedirectList.AddRedirect(i[0], i[1])
		}
	}
	for _, i := range bb.GetAll() {
		localAllowList.Remove(i)
		localRedirectList.Remove(i)
		localBlockList.Add(i)
	}
	for _, i := range bb.GetAllRedirect() {
		localRedirectList.AddRedirect(i[0], i[1])
	}
	err = WriteList(BlockPath, localBlockList)
	if err != nil {
		errs = append(errs, err)
	}
	err = WriteList(AllowPath, localAllowList)
	if err != nil {
		errs = append(errs, err)
	}
	err = WriteRedirectList(RedirectPath, localRedirectList)
	if err != nil {
		errs = append(errs, err)
	}
	if isDownloadSources {
		downList, errs1 := DownloadHostsLists(SourcePath)
		if len(errs1) != 0 {
			for _, i := range errs1 {
				errs = append(errs, i)
			}
		}
		for _, i := range downList {
			bb.Add(i)
		}
	} else {
		for _, i := range ReadLocalHosts(HostPath) {
			bb.Add(i)
		}
	}
	for _, i := range localAllowList.GetAll() {
		bb.Remove(i)
	}
	for _, i := range localBlockList.GetAll() {
		bb.Add(i)
	}
	for _, i := range localRedirectList.GetAllRedirect() {
		bb.AddRedirect(i[0], i[1])
	}
	err = WriteHosts(HostPath, bb)
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

func OptionAllowFn(bb *Set) []error {
	errs := []error{}
	bl := NewSet()
	rdList, err := GetRedirectList(RedirectPath)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, i := range rdList {
			bl.AddRedirect(i[0], i[1])
		}
	}
	if len(bb.GetAll()) == 0 {
		err := WriteHosts(HostPath, bl)
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		return errs
	}
	localBlockList := NewSet()
	localAllowList := NewSet()

	q, _, err := GetList(BlockPath)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, i := range q {
			localBlockList.Add(i)
		}
	}
	for _, i := range bb.GetAll() {
		localAllowList.Add(i)
		localBlockList.Remove(i)
	}
	q, _, err = GetList(AllowPath)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, i := range q {
			localAllowList.Add(i)
			localBlockList.Remove(i)
		}
	}
	for _, i := range ReadLocalHosts(HostPath) {
		bl.Add(i)
	}
	for _, i := range localBlockList.GetAll() {
		bl.Add(i)
	}
	for _, i := range localAllowList.GetAll() {
		bl.Remove(i)
	}
	err = WriteList(BlockPath, localBlockList)
	if err != nil {
		errs = append(errs, err)
	}
	err = WriteList(AllowPath, localAllowList)
	if err != nil {
		errs = append(errs, err)
	}
	err = WriteHosts(HostPath, bl)
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

func AddRedirectListFn(bb *Set) []error {
	errs := []error{}
	localBlockList := NewSet()
	localAllowList := NewSet()
	localRedirectList := NewSet()
	q, _, err := GetList(BlockPath)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, i := range q {
			localBlockList.Add(i)
		}
	}
	q, _, err = GetList(AllowPath)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, i := range q {
			localAllowList.Add(i)
		}
	}
	rdList, err := GetRedirectList(RedirectPath)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, i := range rdList {
			localBlockList.Remove(i[1])
			localAllowList.Remove(i[1])
			localRedirectList.AddRedirect(i[0], i[1])
		}
	}
	for _, i := range bb.GetAllRedirect() {
		localAllowList.Remove(i[1])
		localBlockList.Remove(i[1])
		localRedirectList.AddRedirect(i[0], i[1])
	}
	// for _, i := range bb.GetAllRedirect() {
	// 	localRedirectList.AddRedirect(i[0], i[1])
	// }
	err = WriteList(BlockPath, localBlockList)
	if err != nil {
		errs = append(errs, err)
	}
	err = WriteList(AllowPath, localAllowList)
	if err != nil {
		errs = append(errs, err)
	}
	err = WriteRedirectList(RedirectPath, localRedirectList)
	if err != nil {
		errs = append(errs, err)
	}
	for _, i := range ReadLocalHosts(HostPath) {
		bb.Add(i)
	}

	for _, i := range localAllowList.GetAll() {
		bb.Remove(i)
	}
	for _, i := range localBlockList.GetAll() {
		bb.Add(i)
	}
	for _, i := range localRedirectList.GetAllRedirect() {
		bb.AddRedirect(i[0], i[1])
	}
	err = WriteHosts(HostPath, bb)
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

func RemoveFromRedirect(bb *Set) error {
	localRedirectList := NewSet()
	rdList, err := GetRedirectList(RedirectPath)
	if err != nil {
		return err
	}
	for _, i := range rdList {
		localRedirectList.AddRedirect(i[0], i[1])
	}
	for _, i := range bb.GetAllRedirect() {
		localRedirectList.Remove(i[1])
	}

	err = WriteRedirectList(RedirectPath, localRedirectList)
	if err != nil {
		return err
	}
	v, err := ReadFile(HostPath)
	if err != nil {
		return err
	}
	if R2.MatchString(v) {
		v = R2.ReplaceAllString(v, BuildRedirectList(localRedirectList))
	} else {
		v = fmt.Sprintf("%s\n%s", v, BuildRedirectList(localRedirectList))
	}
	err = WriteFile(HostPath, v)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFromAllowList(s *Set) error {
	localAllowList := NewSet()
	q, _, err := GetList(AllowPath)
	if err != nil {
		return err
	}
	for _, i := range q {
		localAllowList.Add(i)
	}
	for _, i := range s.GetAll() {
		localAllowList.Remove(i)
	}
	return WriteList(AllowPath, localAllowList)
}

func RemoveFromBlockList(s *Set) error {
	localBlockList := NewSet()
	q, _, err := GetList(BlockPath)
	if err != nil {
		return err
	}
	for _, i := range q {
		localBlockList.Add(i)
	}
	for _, i := range s.GetAll() {
		localBlockList.Remove(i)
	}
	return WriteList(BlockPath, localBlockList)
}
