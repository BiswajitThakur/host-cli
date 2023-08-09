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
    want bool
}
func TestIsComment(t *testing.T) {
    var got_want []isCmt = []isCmt{
        isCmt{"val", false},
        isCmt{"   5feAvZ  ", false},
        isCmt{"######", true},
        isCmt{"#    ", true},
        isCmt{" v ###    ", false},
        isCmt{"A5     #", false},
        isCmt{"  ##    ", true},
        isCmt{" val", false},
        isCmt{"# val", true},
        isCmt{"#  666  #", true},
        isCmt{"#     gg", true},
        isCmt{"     gg", false},
        isCmt{"    # gg", true},
        isCmt{"    #  ##", true},
        isCmt{"    ", true},
        isCmt{"", true},
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
        { FindStr(s1, "ca"), 2},
        { FindStr(s1, "ac"), 3},
        { FindStr(s1, "caa"), -1},
        { FindStr(s1, "abc"), 0},
        { FindStr(s1, "ayff"), -1},
        { FindStr(s1, "bc"), 1},
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


func TestFilterHosts(t *testing.T){
    var h string = `# this is comment
0.0.0.0  fff.com
  # 0.0.0.0 abc.in
127.0.0.1 bt.com

#127.0.0.1  hello.com

   
localhost  kk.in.com
##  ###
123.36.0.44       ee.edu


localhost  foo.com
ddd.com  iii.in`

    var got []string = FilterHosts(h)
    var want []string = []string{
        "fff.com",
        "bt.com",
        "kk.in.com",
        "ee.edu",
        "foo.com",
   //   "iii.in"
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
    	[2]string{"aa.com", "eq.com"},
    	[2]string{"yu.in", "ff.oo"},
    	[2]string{"ghgh.in", "kkllu.gg"},
    	[2]string{"google.com", "facebook.com"},
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

