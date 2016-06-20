package structures

import "sync"

// Set is a set
type Set struct {
	Map   map[interface{}]bool
	Mutex sync.Mutex
}

// NewSet constructs a new Set
func NewSet() *Set {
	return &Set{
		Map: map[interface{}]bool{},
	}
}

func (s Set) Has(v interface{}) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.Map[v]
}

func (s Set) Add(v interface{}) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Map[v] = true
}

func (s Set) Remove(v interface{}) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Map, v)
}

// Values
func (s Set) Values() []interface{} {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	r := []interface{}{}
	for v := range s.Map {
		r = append(r, v)
	}
	return r
}
