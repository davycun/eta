package trie_test

import (
	"github.com/davycun/eta/pkg/common/utils/trie"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsContain(t *testing.T) {
	trieNode := trie.InitTrie([]string{"a/b/c/d", "a/b", "a/b/x", "1/a/b"})
	assert.Equal(t, "a/b/c/d", trie.Match(trieNode, "a/b/c/d"))
	assert.Equal(t, "a/b/c", trie.Match(trieNode, "a/b/c"))
	assert.Equal(t, "a/b", trie.Match(trieNode, "a/b"))
	assert.Equal(t, "1", trie.Match(trieNode, "1/x/a/b"))
}
