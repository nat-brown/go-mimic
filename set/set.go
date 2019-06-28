// Package set is a simple set implementation with limited functionality.
// It is intended to just do what is necessary for go-mimic.
package set

// Concept taken from https://yourbasic.org/golang/implement-set/
// Empty struct chosen for memory efficiency.
type void struct{}

var present void

// String is a set of strings.
type String map[string]void

// New returns a new, empty set.
func New() String {
	return map[string]void{}
}

// Union is the union of two sets.
func Union(s1, s2 String) String {
	new := New()
	for k := range s1 {
		new.Add(k)
	}
	for k := range s2 {
		new.Add(k)
	}
	return new
}

// Add adds a new entry to the set.
func (s String) Add(addition string) {
	s[addition] = present
}

// Contains returns whether or not the set contains the given string.
func (s String) Contains(value string) bool {
	_, ok := s[value]
	return ok
}
