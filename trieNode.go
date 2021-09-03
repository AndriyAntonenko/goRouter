package goRouter

import (
	"net/http"
)

type trieNodeType = int8

const (
	static trieNodeType = iota
	root
	pattern
)

type RouterTrieNode struct {
	path     string
	children *RouterHashTable
	nodeType trieNodeType
	handler  *Handler

	subTries []*RouterTrie
}

func NewRouterRootTrieNode(hashSize int) *RouterTrieNode {
	return &RouterTrieNode{
		children: NewRouterHashTable(hashSize),
		nodeType: root,
	}
}

func NewRouterStaticTrieNode(path string, hashSize int) *RouterTrieNode {
	return &RouterTrieNode{
		path:     path,
		nodeType: static,
		children: NewRouterHashTable(hashSize),
	}
}

func NewRouterParamTrieNode(paramName string, hashSize int) *RouterTrieNode {
	subTries := make([]*RouterTrie, 0)
	subTries = append(subTries, NewRouterTrie("*", hashSize, paramName))

	return &RouterTrieNode{
		nodeType: pattern,
		subTries: subTries,
	}
}

func (tn *RouterTrieNode) addChild(node *RouterTrieNode) {
	tn.children.insert(node)
}

func (tn *RouterTrieNode) handleCall(w http.ResponseWriter, r *http.Request, ps *RouterParams) {
	if tn.handler != nil {
		handler := *tn.handler
		handler(w, r, ps)
	}
}
