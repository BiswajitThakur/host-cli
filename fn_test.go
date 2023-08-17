package main

import (
	"os"
	"reflect"
	"testing"
)

var test_file string = "test.file.txt"

func TestIsExist(t *testing.T) {

	// any path that exists
	var got bool = IsExist("fn_test.go")
	var want bool = true
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}

	got = IsExist("any_path_not_exists")
	want = false
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}

func TestIsExistOneOfMany(t *testing.T) {
	var got bool = IsExistOneOfMany("path1", "fn_test.go", "path2")
	var want bool = true
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}

	// dont't create any file named path1, path2, path3
	got = IsExistOneOfMany("path1", "path2", "path3")
	want = false
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}

func TestReadFile(t *testing.T) {

	var _, got = ReadFile("fn_test.go")
	var want error = nil
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}

	_, got = ReadFile("any_path_not_exists")
	if got == want {
		t.Error("no such file or directory, error require\n")
	}
}

type isCmt struct {
	input string
	want  bool
}

func TestIsComment(t *testing.T) {
	var got_want []isCmt = []isCmt{
		{"val", false},
		{"   5feAvZ  ", false},
		{"######", true},
		{"#    ", true},
		{" v ###    ", false},
		{"A5     #", false},
		{"  ##    ", true},
		{" val", false},
		{"# val", true},
		{"#  666  #", true},
		{"#     gg", true},
		{"     gg", false},
		{"    # gg", true},
		{"    #  ##", true},
		{"    ", true},
		{"", true},
	}
	for _, i := range got_want {
		got := IsComment(i.input)
		if got != i.want {
			t.Errorf("got: %v, want: %v\n", got, i.want)
		}
	}
}

func TestFindStr(t *testing.T) {
	var s1 = []string{"abc", "bc", "ca", "ac"}
	got_want := [][]int{
		{FindStr(s1, "ca"), 2},
		{FindStr(s1, "ac"), 3},
		{FindStr(s1, "caa"), -1},
		{FindStr(s1, "abc"), 0},
		{FindStr(s1, "ayff"), -1},
		{FindStr(s1, "bc"), 1},
	}
	for _, i := range got_want {
		got := i[0]
		want := i[1]
		if got != want {
			t.Errorf("got: %v, want: %v\n", got, want)
		}
	}
}

func TestWriteFile(t *testing.T) {
	var got error = WriteFile(test_file, "*Hello World*")
	var want error = nil
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}

	var got1 bool = IsExist(test_file)
	var want1 bool = true
	if got1 != want1 {
		t.Errorf("got: %v, want: %v\n", got1, want1)
	}

	got2, _ := ReadFile(test_file)
	want2 := "*Hello World*"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v\n", got2, want2)
	}
	_ = os.Remove(test_file)
}

func TestFilterHosts(t *testing.T) {
	var h string = `# this is comment
0.0.0.0  fff.com
  # 0.0.0.0 abc.in
127.0.0.1 bt.com

#127.0.0.1  hello.com

   
localhost  kk.in.com
##  ###
123.36.0.44       ee.edu


localhost  foo.com
ddd.com  iii.in
77.80.120.55   googgggll.com
56.34.100.33  102.78.111.111
`

	var got []string = FilterHosts(h)
	var want []string = []string{
		"fff.com",
		"bt.com",
		"kk.in.com",
		"ee.edu",
		"foo.com",
		"googgggll.com",
		"102.78.111.111",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}

func TestFilterRedirectHosts(t *testing.T) {
	var input string = `# Hello
aa.com  eq.com

  yu.in   ff.oo
# ss.com  qwh.com
#####################
ghgh.in   kkllu.gg
# #####
##### #####
    #   eeuu.in   vvvkkk.com
google.com    facebook.com
`
	var got [][2]string = FilterRedirectHosts(input)
	var want [][2]string = [][2]string{
		{"aa.com", "eq.com"},
		{"yu.in", "ff.oo"},
		{"ghgh.in", "kkllu.gg"},
		{"google.com", "facebook.com"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}

func TestReadLocalHosts(t *testing.T) {
	input := `127.0.0.1   localhost
###########
0.0.0.0   kkk.com
# BT-start #
0.0.0.0        yooo.in
127.0.0.1   fooooo.com
#  127.0.0.1   example.com
###   ####

123.66.40.222  google.com

# BT-end  #
127.0.0.1   ioio.io
0.0.0.0   kkk.com
`
	_ = WriteFile(test_file, input)
	var got []string = ReadLocalHosts(test_file)
	var want []string = []string{
		"yooo.in",
		"fooooo.com",
		"google.com",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
	_ = os.Remove(test_file)
}

func TestReadLocalRedirectHosts(t *testing.T) {
	input := `127.0.0.1   localhost
###########
0.0.0.0   kkk.com
# BT-start #
0.0.0.0        yooo.in
127.0.0.1   fooooo.com
#  127.0.0.1   example.com
###   ####

123.66.40.222  google.com

# BT-end  #
127.0.0.1   ioio.io
# BT-redirect-start #
hello.com   world.com
# hh77.in       iu.vn
rudra.ghos.com   73.100.78.99
#############
    fuck.com fucking.com

### ### gyg.com hgg.com
123.45.78.99  rudra.com
77.111.55.89             56.98.12.88

# BT-redirect-end #
0.0.0.0   kkk.com`
	_ = WriteFile(test_file, input)
	got1, got2 := ReadLocalRedirectHosts(test_file)
	want1 := [][2]string{
		{"hello.com", "world.com"},
		{"rudra.ghos.com", "73.100.78.99"},
		{"fuck.com", "fucking.com"},
		{"123.45.78.99", "rudra.com"},
		{"77.111.55.89", "56.98.12.88"},
	}
	var want2 error = nil

	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("got: %v, want: %v\n", got1, want1)
	}

	if got2 != want2 {
		t.Errorf("got: %v, want: %v\n", got2, want2)
	}
	_ = os.Remove(test_file)
}

func TestBuildBlockList(t *testing.T) {
	lst := NewSet()
	lst.Add("aa.com")
	lst.Add("bb.in")
	lst.AddRedirect("zzzz.com", "yyy.in")
	lst.AddRedirect("ppp.com", "uuuuu.com")
	var got string = BuildBlockList(lst)
	var want string = `# BT-start #
0.0.0.0 aa.com
0.0.0.0 bb.in
# BT-end #`
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}

func TestBuildRedirectList(t *testing.T) {
	lst := NewSet()
	lst.Add("aa.com")
	lst.Add("bb.in")
	lst.AddRedirect("zzzz.com", "yyy.in")
	lst.AddRedirect("ppp.com", "uuuuu.com")
	var got string = BuildRedirectList(lst)
	var want string = `# BT-redirect-start #
ppp.com uuuuu.com
zzzz.com yyy.in
# BT-redirect-end #`
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}
