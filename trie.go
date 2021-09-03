package router

import (
	"strings"
)

type RouterTrie struct {
	paramName string
	root      *RouterTrieNode
	hashSize  int
}

func NewRouterTrie(headValue string, hashSize int, paramName string) *RouterTrie {
	return &RouterTrie{
		root:      NewRouterRootTrieNode(hashSize),
		hashSize:  hashSize,
		paramName: paramName,
	}
}

func (t *RouterTrie) AddNode(path string, handler Handler) {
	if path == "/" {
		t.root.handler = &handler
		return
	}

	parts := strings.Split(path, "/")
	currentNode := t.root

	if len(parts) > 1 && parts[0] == "" {
		parts = parts[1:]
	}

	for index, part := range parts {
		var child *RouterTrieNode

		paramName, isParam := getParam(part)

		if isParam {
			child = currentNode.children.lookupPattern()
		} else {
			child = currentNode.children.lookupStatic(part)
		}

		if child == nil {
			if isParam {
				child = NewRouterParamTrieNode(paramName, t.hashSize)
				currentNode.addChild(child)
				currentNode = child.subTries[0].root
			} else {
				child = NewRouterStaticTrieNode(part, t.hashSize)
				currentNode.addChild(child)
				currentNode = child
			}

			continue
		}

		if child.nodeType == pattern && child.subTries != nil {
			for _, trie := range child.subTries {
				if trie.paramName == paramName {
					trie.AddNode(strings.Join(parts[index+1:], "/"), handler)
					return
				}
			}

			newTrie := NewRouterTrie("", t.hashSize, paramName)
			child.subTries = append(child.subTries, newTrie)
			child = newTrie.root
			currentNode = child
		}

		currentNode = child
	}

	currentNode.handler = &handler
}

// TODO: FIX THIS METHOD!!!
func (t *RouterTrie) Lookup(path string, ps *RouterParams) *RouterTrieNode {
	if path == "/" || path == "" {
		return t.root
	}

	parts := strings.Split(path, "/")
	currentNode := t.root

	if len(parts) > 1 && parts[0] == "" {
		parts = parts[1:]
	}

	for index, part := range parts {

		child := currentNode.children.lookupStatic(part)
		if child == nil {
			child = currentNode.children.lookupPattern()

			if child == nil {
				return nil
			}

			if child.subTries != nil {
				for _, trie := range child.subTries {
					subTrieChild := trie.Lookup(strings.Join(parts[index+1:], "/"), ps)
					if subTrieChild != nil {
						ps.addParam(trie.paramName, part)
						return subTrieChild
					}
				}

				return nil
			}
		}

		currentNode = child
	}

	return currentNode
}

func getParam(part string) (paramName string, isParam bool) {
	firstChar := part[0]
	paramName = ""
	isParam = false

	if string(firstChar) == ":" {
		paramName = string(part[1:])
		isParam = true
	}

	return
}
