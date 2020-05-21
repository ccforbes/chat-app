package indexes

import (
	"sort"
	"sync"
)

//int64set is a set of int64 values
type int64set map[int64]struct{}

//add adds a value to the set and returns
//true if the value didn't already exist in the set
func (s int64set) add(value int64) bool {
	_, found := s[value]
	if found {
		return false
	}
	s[value] = struct{}{}
	return true
}

//remove removes a value from the set and returns
//true if that value was in the set, false otherwise
func (s int64set) remove(value int64) bool {
	_, found := s[value]
	if !found {
		return false
	}
	delete(s, value)
	return true
}

//has returns true if value is in the set,
//or false if it is not in the set
func (s int64set) has(value int64) bool {
	_, found := s[value]
	return found
}

func (s int64set) all() []int64 {
	var values []int64
	for value := range s {
		values = append(values, value)
	}
	return values
}

//TrieNode implements a trie data structure mapping strings to int64s
//that is safe for concurrent use
type TrieNode struct {
	Children map[rune]*TrieNode
	Values   int64set
	mx       sync.RWMutex
}

//NewTrieNode constructs a new TrieNode
func NewTrieNode() *TrieNode {
	return &TrieNode{}
}

//Len returns the number of entries in the trie
func (t *TrieNode) Len() int {
	t.mx.RLock()
	defer t.mx.RUnlock()
	entryCount := t.numOfEntries()
	return entryCount
}

//len is a private helper of Len that recursively sums
//the number of values in the trie
func (t *TrieNode) numOfEntries() int {
	entryCount := len(t.Values)
	for child := range t.Children {
		entryCount += t.Children[child].numOfEntries()
	}
	return entryCount
}

//Add adds a key and value to the trie
func (t *TrieNode) Add(key string, value int64) {
	t.mx.Lock()
	defer t.mx.Unlock()
	runes := []rune(key)
	t.add(runes, value)
}

// add is a private helper method that adds a key and value to the trie
func (t *TrieNode) add(key []rune, value int64) {
	// if children do not exist, make sure it is an empty map
	if len(t.Children) == 0 {
		t.Children = make(map[rune]*TrieNode)
	}
	// if the child does not exist, create a new trie node and store it there
	if t.Children[key[0]] == nil {
		t.Children[key[0]] = NewTrieNode()
	}
	if len(key) == 1 {
		// if the child values dont exist, make sure it is an empty int64set
		// you could probably do this by doing int64set{}
		if len(t.Children[key[0]].Values) == 0 {
			t.Children[key[0]].Values = make(map[int64]struct{})
		}
		// add the value and then return
		t.Children[key[0]].Values.add(value)
		return
	}
	// otherwise, call the add method again on the child node (recursively)
	t.Children[key[0]].add(key[1:len(key)], value)
}

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *TrieNode) Find(prefix string, max int) []int64 {
	t.mx.RLock()
	defer t.mx.RUnlock()

	if len(t.Children) == 0 || prefix == "" || max <= 0 {
		return nil
	}

	// iterate through trie until at end of prefix | O(1)
	currNode := t
	for _, char := range prefix {
		if currNode.Children[char] == nil {
			return nil
		}
		currNode = currNode.Children[char]
	}
	// create int64 slice
	var returnSlice []int64
	currNode.find(&returnSlice, max)
	return returnSlice
}

func (t *TrieNode) find(list *[]int64, max int) {
	// add all current values in node to list (or until hit max)
	values := t.Values.all()
	canGet := max - len(*list)
	if len(values) > canGet {
		*list = append(*list, values[0:canGet]...)
		return
	}
	*list = append(*list, values...)
	// if max reached or no children, just return.
	if len(*list) == max || len(t.Children) == 0 {
		return
	}
	// sort children
	children := make([]rune, 0, len(t.Children))
	for k := range t.Children {
		children = append(children, k)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i] < children[j]
	})
	// for every child, recurse and add to list and check for max
	for _, child := range children {
		t.Children[child].find(list, max)
		if len(*list) == max {
			return
		}
	}
	return
}

//Remove removes a key/value pair from the trie
//and trims branches with no values.
func (t *TrieNode) Remove(key string, value int64) {
	// split key into runes
	t.mx.Lock()
	defer t.mx.Unlock()
	runes := []rune(key)
	t.remove(runes, value)
}

func (t *TrieNode) remove(key []rune, value int64) {
	if len(key) == 0 {
		t.Values.remove(value)
		return
	}
	focusChild, ok := t.Children[key[0]]
	if !ok {
		return
	}
	focusChild.remove(key[1:], value)
	if len(focusChild.Children) == 0 && len(focusChild.Values) == 0 {
		delete(t.Children, key[0])
	}
}
