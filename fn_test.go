package main

import (
    "testing"
)

func TestIsExist(t *testing.T){

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

func TestIsExistOneOfMany(t *testing.T){
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

func TestReadFile(t *testing.T){

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
func TestIsComment(t *testing.T){
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

func TestFindStr(t *testing.T){
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
