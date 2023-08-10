package main

import (
	"reflect"
	"testing"
)

func TestGetAll(t *testing.T) {
	var set1 = NewSet()
	set1.Add("127.0.0.1")
	set1.Add("rudra.com")
	set1.Add("80.121.44.100")
	set1.Add("example.com")
	set1.Add("rahul.sarkar.com")
	set1.Add("sathi.guchait.in")
	set1.AddRedirect("hello.com", "world.com")
	set1.AddRedirect("foo.in", "zzz.com")
	set1.AddRedirect("10.77.19.100", "fb.com")
	set1.AddRedirect("67.80.13.111", "111.100.40.35")
	got := set1.GetAll()
	want := []string{
		"127.0.0.1",
		"80.121.44.100",
		"example.com",
		"rahul.sarkar.com",
		"rudra.com",
		"sathi.guchait.in",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}

func TestGetAllRedirect(t *testing.T) {
	var set1 = NewSet()
	set1.Add("127.0.0.1")
	set1.Add("rudra.com")
	set1.Add("80.121.44.100")
	set1.Add("example.com")
	set1.Add("rahul.sarkar.com")
	set1.Add("sathi.guchait.in")
	set1.AddRedirect("hello.com", "world.com")
	set1.AddRedirect("foo.in", "zzz.com")
	set1.AddRedirect("10.77.19.100", "fb.com")
	set1.AddRedirect("67.80.13.111", "111.100.40.35")
	got := set1.GetAllRedirect()
	want := [][2]string{
		{"67.80.13.111", "111.100.40.35"},
		{"10.77.19.100", "fb.com"},
		{"hello.com", "world.com"},
		{"foo.in", "zzz.com"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}

func TestContains(t *testing.T) {
	var set1 = NewSet()
	set1.Add("127.0.0.1")
	set1.Add("rudra.com")
	set1.AddRedirect("hello.com", "world.com")
	set1.AddRedirect("foo.in", "zzz.com")

	got1 := set1.Contains("hello.com")
	want1 := false
	if got1 != want1 {
		t.Errorf("got: %v, want: %v\n", got1, want1)
	}

	got2 := set1.Contains("world.com")
	want2 := false
	if got2 != want2 {
		t.Errorf("got: %v, want: %v\n", got2, want2)
	}

	got3 := set1.Contains("rudra.com")
	want3 := true
	if got3 != want3 {
		t.Errorf("got: %v, want: %v\n", got3, want3)
	}

	got4 := set1.ContainsRedirect("tooo.com")
	want4 := false
	if got4 != want4 {
		t.Errorf("got: %v, want: %v\n", got4, want4)
	}
}

func TestContainsRedirect(t *testing.T) {
	var set1 = NewSet()
	set1.Add("127.0.0.1")
	set1.Add("rudra.com")
	set1.AddRedirect("hello.com", "world.com")
	set1.AddRedirect("foo.in", "zzz.com")

	got1 := set1.ContainsRedirect("hello.com")
	want1 := false
	if got1 != want1 {
		t.Errorf("got: %v, want: %v\n", got1, want1)
	}

	got2 := set1.ContainsRedirect("world.com")
	want2 := true
	if got2 != want2 {
		t.Errorf("got: %v, want: %v\n", got2, want2)
	}

	got3 := set1.ContainsRedirect("rudra.com")
	want3 := false
	if got3 != want3 {
		t.Errorf("got: %v, want: %v\n", got3, want3)
	}

	got4 := set1.ContainsRedirect("tooo.com")
	want4 := false
	if got4 != want4 {
		t.Errorf("got: %v, want: %v\n", got4, want4)
	}
}

func TestRemove(t *testing.T) {
	var set1 = NewSet()
	set1.Add("127.0.0.1")
	set1.Add("rudra.com")
	set1.AddRedirect("hello.com", "world.com")
	set1.AddRedirect("foo.in", "zzz.com")

	got := set1.Contains("rudra.com")
	want := true
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}

	set1.Remove("rudra.com")

	got = set1.Contains("rudra.com")
	want = false
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}

	got = set1.Contains("zzz.com")
	want = true
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}

	set1.Remove("zzz.com")

	got = set1.Contains("zzz.com")
	want = false
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}

func TestAdd(t *testing.T) {
	var set1 = NewSet()
	set1.Add("127.0.0.1")

	got := set1.Contains("rudra.com")
	want := false
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}

	set1.Add("rudra.com")
	got = set1.Contains("rudra.com")
	want = true
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}

func TestAddRedirect(t *testing.T) {
	var set1 = NewSet()
	set1.Add("127.0.0.1")
	set1.AddRedirect("hello.com", "world.com")

	got := set1.Contains("rudra.com")
	want := false
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}

	set1.AddRedirect("ghose.com", "rudra.com")
	got = set1.Contains("rudra.com")
	want = true
	if got != want {
		t.Errorf("got: %v, want: %v\n", got, want)
	}
}
