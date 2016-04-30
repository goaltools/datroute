package datroute

import (
	"net/http"
)

// Dict is a dictionary structure that is used by routing package instead of map
// for small sets of data.
// On average efficency of getting an element from map is O(c + 1).
// At the same time efficency of iterating over a slice is O(n).
// And when n is small, O(n) < O(c + 1). That's why we are using slice and simple loop
// rather than a map.
type Dict struct {
	Keys   []string
	Values []*http.HandlerFunc
}

// NewDict allocates and returns a dict structure.
func NewDict() *Dict {
	return &Dict{
		Keys:   []string{},
		Values: []*http.HandlerFunc{},
	}
}

// Set expects key and value as input parameters that are
// saved to the dict.
func (t *Dict) Set(k string, v *http.HandlerFunc) {
	// Check whether we have already had such key.
	if _, i := t.Get(k); i >= 0 {
		// If so, update it.
		t.Values[i] = v
		return
	}
	// Otherwise, add a new key-value pair.
	t.Keys = append(t.Keys, k)
	t.Values = append(t.Values, v)
}

// Get receives a key as input and returns associated value
// and its index. If the value is not found nil, -1 are returned.
func (t *Dict) Get(k string) (*http.HandlerFunc, int) {
	for i := range t.Keys {
		if t.Keys[i] == k {
			return t.Values[i], i
		}
	}
	return nil, -1
}

// Join receives a new dict and joins with the old one
// calling Set for every key - value pair.
func (t *Dict) Join(d *Dict) {
	// Iterate through all keys of a new dict.
	for i := range d.Keys {
		// Add them to the main dict.
		t.Set(d.Keys[i], d.Values[i])
	}
}
