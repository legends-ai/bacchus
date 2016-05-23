package structures

type StringSet map[string]bool

func (s StringSet) Has(v string) bool {
	return s[v]
}

func (s StringSet) Add(v string) {
	s[v] = true
}

func (s StringSet) Remove(v string) {
	delete(s, v)
}

// Values
func (s StringSet) Values() []string {
	r := []string{}
	for v := range s {
		r = append(r, v)
	}
	return r
}
