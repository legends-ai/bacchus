package structures

import "sync"

// StringSet is a set that contains strings
type StringSet struct {
	Map   map[string]bool
	Mutex sync.Mutex
}

// NewStringSet constructs a new StringSet
func NewStringSet() *StringSet {
	return &StringSet{
		Map: map[string]bool{},
	}
}

func (s StringSet) Has(v string) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.Map[v]
}

func (s StringSet) Add(v string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Map[v] = true
}

func (s StringSet) Remove(v string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Map, v)
}

// Values
func (s StringSet) Values() []string {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	r := []string{}
	for v := range s.Map {
		r = append(r, v)
	}
	return r
}
