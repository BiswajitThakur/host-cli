package main

import "sort"

var exists = struct{}{}

type Set struct {
	m map[string]struct{}
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]struct{})
	return s
}

func (s *Set) Add(value string) {
	s.m[value] = exists
}

func (s *Set) Remove(value string) {
	delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func (s *Set) GetAll() []string {
	var v []string = []string{}
	for key, _ := range s.m {
		v = append(v, key)
	}
	sort.Strings(v)
	return v
}
