package main

import "sort"

var exists = "1"

type Set struct {
	m map[string]string
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]string)
	return s
}

func (s *Set) Add(value string) {
	s.m[value] = exists
}

func (s *Set) AddRedirect(a string, b string) {
	s.m[b] = a
}

func (s *Set) Remove(value string) {
	delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func (s *Set) ContainsRedirect(val string) bool {
	_, b := s.m[val]
	c := (s.m[val] != exists) && b
	return c
}

func (s *Set) GetAll() []string {
	var v []string = []string{}
	for key, val := range s.m {
		if val == exists {
			v = append(v, key)
		}
	}
	sort.Strings(v)
	return v
}

func (s *Set) GetAllRedirect() [][2]string {
	var v [][2]string = [][2]string{}
	for key, val := range s.m {
		if val != exists {
			v = append(v, [2]string{val, key})
		}
	}
	sort.SliceStable(v, func(i, j int) bool {
		return v[i][1] < v[j][1]
	})
	return v
}
