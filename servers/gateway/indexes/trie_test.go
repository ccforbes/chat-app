package indexes

import (
	"reflect"
	"testing"
)

//TODO: implement automated tests for your trie data structure
func TestTrieAddAndFind(t *testing.T) {
	cases := []struct {
		name          string
		hint          string
		root          *TrieNode
		entries       map[string]int64
		prefix        string
		correctValues []int64
		max           int
		isNil         bool
	}{
		{
			"Get Single Entry",
			"Check to see if you're adding each letter to a child node and recursing down each node",
			NewTrieNode(),
			map[string]int64{"test": 1},
			"t",
			[]int64{1},
			1,
			false,
		},
		{
			"Get Double Entry",
			"Check to see if you're adding each letter to a child node and recursing down each node",
			NewTrieNode(),
			map[string]int64{"test": 1, "testtwo": 2},
			"t",
			[]int64{1, 2},
			2,
			false,
		},
		{
			"Get Triple Entry",
			"Check to see if you're adding each letter to a child node and recursing down each node",
			NewTrieNode(),
			map[string]int64{"test": 1, "two": 2, "three": 3},
			"t",
			[]int64{1, 3, 2},
			3,
			false,
		},
		{
			"Get One From Triple Entry",
			"Make sure you stop after retrieving values after you've hit the max",
			NewTrieNode(),
			map[string]int64{"test": 1, "two": 2, "three": 3},
			"t",
			[]int64{1},
			1,
			false,
		},
		{
			"Get Two From Triple Entry",
			"Make sure you stop after retrieving values after you've hit the max",
			NewTrieNode(),
			map[string]int64{"test": 1, "two": 2, "three": 3},
			"t",
			[]int64{1, 3},
			2,
			false,
		},
		{
			"Get One Value From Three Entries",
			"Make sure you stop after retrieving values after you've hit the max",
			NewTrieNode(),
			map[string]int64{"test": 1, "two": 2, "three": 3},
			"te",
			[]int64{1},
			2,
			false,
		},
		{
			"No Values",
			"No results should be shown for a search prefix that doesn't point to any nodes",
			NewTrieNode(),
			map[string]int64{"test": 1, "two": 2, "three": 3},
			"b",
			[]int64{},
			3,
			true,
		},
		{
			"No Max",
			"Find should returned nil if there max is 0",
			NewTrieNode(),
			map[string]int64{"test": 1, "two": 2, "three": 3},
			"b",
			nil,
			0,
			true,
		},
		{
			"Empty Tree",
			"Find should returned nil if the tree is empty",
			NewTrieNode(),
			map[string]int64{},
			"test",
			nil,
			5,
			true,
		},
		{
			"Empty Prefix",
			"Find should returned nil if the prefix is empty",
			NewTrieNode(),
			map[string]int64{"test": 1, "two": 2, "three": 3},
			"",
			nil,
			3,
			true,
		},
	}

	for _, c := range cases {
		for key, value := range c.entries {
			c.root.Add(key, value)
		}
		foundValues := c.root.Find(c.prefix, c.max)
		if reflect.DeepEqual(foundValues, c.correctValues) == false && c.isNil == false {
			t.Errorf("case %s: incorrect values returned\nexpected values: %v; returned: %v\nHINT: %s", c.name, c.correctValues, foundValues, c.hint)
		}
		if c.isNil && len(foundValues) != 0 {
			t.Error("Seaching for nonexistent prefix should return no values")
		}
	}
}

func TestLen(t *testing.T) {
	cases := []struct {
		name         string
		hint         string
		root         *TrieNode
		entries      map[string]map[int64]struct{}
		numOfEntries int
	}{
		{
			"No Entries",
			"An empty tree should return 0",
			NewTrieNode(),
			map[string]map[int64]struct{}{},
			0,
		},
		{
			"One Entry",
			"Make sure to recursively add the number of entries only when there are values at the node",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {1: {}}},
			1,
		},
		{
			"Two Entries",
			"Make sure to recursively add the number of entries only when there are values at the node",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {1: {}}, "two": {2: {}}},
			2,
		},
		{
			"Three Entries",
			"Make sure to recursively add the number of entries only when there are values at the node",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {1: {}}, "two": {2: {}}, "bop": {3: {}}},
			3,
		},
		{
			"One Word, Two Values",
			"Make sure to recursively add the number of entries only when there are values at the node",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {1: {}, 2: {}}},
			2,
		},
		{
			"One Word, Three Values",
			"Make sure to recursively add the number of entries only when there are values at the node",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {1: {}, 2: {}, 3: {}}},
			3,
		},
	}
	for _, c := range cases {
		for key, valueSet := range c.entries {
			for value := range valueSet {
				c.root.Add(key, value)
			}
		}
		returnedLength := c.root.Len()
		if returnedLength != c.numOfEntries {
			t.Errorf("case %s: incorrect num of entries returned\nexpected: %d; returned %d\nHINT: %s", c.name, c.numOfEntries, returnedLength, c.hint)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name         string
		hint         string
		root         *TrieNode
		entries      map[string]map[int64]struct{}
		keys         []string
		values       []int64
		numOfEntries int
	}{
		{
			"Deleting One Item",
			"There should still be items available",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {1: {}}, "two": {2: {}}},
			[]string{"test"},
			[]int64{1},
			1,
		},
		{
			"Deleting All Items",
			"There should be no items available",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {1: {}}, "two": {2: {}}},
			[]string{"test", "two"},
			[]int64{1, 2},
			0,
		},
		{
			"Deleting One Value from a Node's Value Set",
			"There should still be a value left in one set",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {2: {}, 1: {}}},
			[]string{"test"},
			[]int64{2},
			1,
		},
		{
			"Deleting Key/Value but Node Has Children",
			"There should still be access to these nodes because the node with deleted values has children",
			NewTrieNode(),
			map[string]map[int64]struct{}{"testing": {2: {}}, "test": {1: {}}},
			[]string{"test"},
			[]int64{1},
			1,
		},
		{
			"Invalid Key",
			"Nothing should be deleted, key is invalid",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {1: {}}},
			[]string{"other"},
			[]int64{1},
			1,
		},
		{
			"Invalid Value",
			"Nothing should be deleted, value is invalid",
			NewTrieNode(),
			map[string]map[int64]struct{}{"test": {1: {}}},
			[]string{"test"},
			[]int64{2},
			1,
		},
	}

	for _, c := range cases {
		for key, valueSet := range c.entries {
			for value := range valueSet {
				c.root.Add(key, value)
			}
		}
		for i, key := range c.keys {
			c.root.Remove(key, c.values[i])
		}
		if c.numOfEntries != c.root.Len() {
			t.Errorf("case %s: delete not performed correctly\nexpected num of entries: %d; returned %d\nHINT: %s", c.name, c.numOfEntries, c.root.Len(), c.hint)
		}
	}
}
