package trie

import (
	"github.com/duke-git/lancet/v2/slice"
	"strings"
)

type Node struct {
	Children map[string]*Node
}

// InitTrie 初始化Trie树. 从最顶级开始，“/”分割。d1/d2/d3
func InitTrie(words []string) *Node {
	root := &Node{}
	for _, w := range words {
		AddWord(root, w)
	}
	return root
}
func AddWord(root *Node, word string) {
	ds := strings.Split(word, "/")
	node := root
	for _, d := range ds {
		if node.Children == nil {
			node.Children = make(map[string]*Node)
		}
		if _, ok := node.Children[d]; !ok {
			node.Children[d] = &Node{}
		}
		node = node.Children[d]
	}
}
func Contain(root *Node, word string) bool {
	ds := strings.Split(word, "/")
	for i := 0; i < len(ds); i++ {
		p := root
		j := i
		for j < len(ds) && p.Children != nil {
			d := ds[j]
			if _, ok := p.Children[d]; ok {
				p = p.Children[d]
				j++
			} else {
				break
			}
		}
		if p.Children == nil {
			return true
		}
	}
	return false
}
func Match(root *Node, word string) string {
	m := make([]string, 0)
	ds := strings.Split(word, "/")
	p := root
	for i := 0; i < len(ds); i++ {
		d := ds[i]
		if _, ok := p.Children[d]; ok {
			p = p.Children[d]
			m = append(m, d)
		} else {
			break
		}
	}
	return slice.Join(m, "/")
}
