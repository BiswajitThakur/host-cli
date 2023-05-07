package options

import (
	"fmt"
	"regexp"
)

type Option struct {
	SortName string
	LongName string
	Callback func()
}

type B struct {
	name string
	op   string
	lst  []string
}

type OptionList struct {
	ops *B
	lst []Option
}

func NewOptions(args []string) *OptionList {
	v := &B{}
	v.name = args[0]
	if len(args) >= 2 {
		v.op = args[1]
	}
	if len(args) >= 3 {
		v.lst = args[2:]
	}
	return &OptionList{
		ops: v,
		lst: []Option{},
	}
}

func (f *OptionList) Add(k Option) {
	f.lst = append(f.lst, k)
}

func (f *OptionList) Args() []string {
	return f.ops.lst
}

func (f *OptionList) Start() {
	var isLong bool = !regexp.MustCompile(`^\s*\-[^\-]`).MatchString(f.ops.op)
	var oo string = regexp.MustCompile(`^\-+`).ReplaceAllString(f.ops.op, "")
	for _, i := range f.lst {
		if isLong {
			if regexp.MustCompile(fmt.Sprintf("(?i)^\\s*\\-*%s\\s*$", oo)).MatchString(i.LongName) {
				i.Callback()
				return
			}
		} else {
			if regexp.MustCompile(fmt.Sprintf("(?i)^\\s*\\-+%s\\s*$", oo)).MatchString(i.SortName) {
				i.Callback()
				return
			}
		}
	}
	fmt.Printf("Invalid option : %s\n", f.ops.op)
}

// func main() {
// 	//n := NewArgs(os.Args)
// 	m := NewOptions(os.Args)
// 	m.Add(Option{
// 		"-h",
// 		"--help",
// 		help,
// 	})
// 	m.Add(Option{
// 		"-b",
// 		"--block",
// 		func() { fmt.Println(888) },
// 	})
// 	m.Start()
// 	fmt.Println(m.Args())
// }

// func help() {
// 	fmt.Println(111111)
// }
