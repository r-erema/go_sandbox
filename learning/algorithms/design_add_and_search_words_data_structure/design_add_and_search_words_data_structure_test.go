package designaddandsearchwordsdatastructure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordDictionary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		scenarioFn func(t *testing.T)
	}{
		{
			name: "scenario 1",
			scenarioFn: func(t *testing.T) {
				t.Helper()

				dict := Constructor()
				dict.AddWord("bad")
				dict.AddWord("dad")
				dict.AddWord("mad")
				assert.False(t, dict.Search("pad"))
				assert.True(t, dict.Search("bad"))
				assert.True(t, dict.Search(".ad"))
				assert.True(t, dict.Search("..d"))
				assert.True(t, dict.Search("b.."))
			},
		},
		{
			name: "scenario 2",
			scenarioFn: func(t *testing.T) {
				t.Helper()

				dict := Constructor()
				dict.AddWord("a")
				dict.AddWord("ab")
				assert.True(t, dict.Search("a"))
				assert.True(t, dict.Search("a."))
			},
		},
		{
			name: "scenario 3",
			scenarioFn: func(t *testing.T) {
				t.Helper()

				dict := Constructor()
				dict.AddWord("at")
				dict.AddWord("and")
				dict.AddWord("an")
				dict.AddWord("add")
				assert.False(t, dict.Search("a"))
				assert.False(t, dict.Search(".at"))

				dict.AddWord("bat")

				assert.True(t, dict.Search(".at"))
			},
		},
		{
			name: "scenario 4",
			scenarioFn: func(t *testing.T) {
				t.Helper()

				dict := Constructor()
				assert.False(t, dict.Search("a"))
			},
		},
		{
			name: "scenario 5",
			scenarioFn: func(t *testing.T) {
				t.Helper()

				dict := Constructor()
				dict.AddWord("a")
				dict.AddWord("ab")
				assert.True(t, dict.Search("."))
				assert.True(t, dict.Search(".."))
			},
		},
		{
			name: "scenario 6",
			scenarioFn: func(t *testing.T) {
				t.Helper()

				dict := Constructor()
				dict.AddWord("at")
				dict.AddWord("and")
				dict.AddWord("an")
				dict.AddWord("add")
				dict.AddWord("bat")
				assert.False(t, dict.Search("b."))
				assert.True(t, dict.Search(".."))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.scenarioFn(t)
		})
	}
}

type WordDictionary struct {
	endOfWord bool
	children  [26]*WordDictionary
}

func Constructor() WordDictionary {
	return WordDictionary{}
}

// Time O(n), since we need iterate entire tree
// Space O(t+n), where n is the length of the string and t is the total number of TrieNodes created in the Trie.
func (trie *WordDictionary) AddWord(word string) {
	curr := trie

	for _, i := range word {
		idx := i - 'a'
		if curr.children[idx] == nil {
			curr.children[idx] = &WordDictionary{}
		}

		curr = curr.children[idx]
	}

	curr.endOfWord = true
}

// Time O(n), since we need iterate entire tree
// Space O(t+n), where n is the length of the string and t is the total number of TrieNodes created in the Trie.
func (trie *WordDictionary) Search(word string) bool {
	var dfs func(i int, node *WordDictionary) bool

	dfs = func(i int, node *WordDictionary) bool {
		for ; i < len(word); i++ {
			char := word[i]

			if char != '.' {
				idx := char - 'a'
				if node == nil || node.children[idx] == nil {
					return false
				}

				node = node.children[idx]

				continue
			}

			for _, child := range node.children {
				if child != nil && dfs(i+1, child) {
					return true
				}
			}

			return false
		}

		return node.endOfWord
	}

	return dfs(0, trie)
}
