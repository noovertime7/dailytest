package main

import (
	"fmt"
	"strings"
)

// TrieNode 定义前缀树的节点结构
type TrieNode struct {
	children    map[string]*TrieNode
	isEnd       bool
	permissions []string // 存储权限信息
}

// Trie 定义前缀树的结构
type Trie struct {
	root *TrieNode
}

// NewTrie 创建一个新的前缀树
func NewTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[string]*TrieNode)}}
}

// Insert 向前缀树中插入路径及其对应的权限信息
func (t *Trie) Insert(path string, permissions []string) {
	node := t.root
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == "" {
			continue
		}
		if _, found := node.children[part]; !found {
			node.children[part] = &TrieNode{children: make(map[string]*TrieNode)}
		}
		node = node.children[part]
	}
	node.isEnd = true
	node.permissions = permissions
}

// CheckPermissions 检查是否具有所需的权限
func (t *Trie) CheckPermissions(path string, requiredPermissions []string) bool {
	node := t.root
	var maxMatchNode *TrieNode

	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == "" {
			continue
		}

		// 精确匹配
		if nextNode, found := node.children[part]; found {
			node = nextNode
			if node.isEnd {
				maxMatchNode = node
			}
			continue
		}

		// 动态参数匹配
		dynamicMatchFound := false
		for key, nextNode := range node.children {
			if strings.HasPrefix(key, ":") {
				node = nextNode
				if node.isEnd {
					maxMatchNode = node
				}
				dynamicMatchFound = true
				break
			}
		}
		if !dynamicMatchFound {
			return false
		}
	}

	if maxMatchNode == nil {
		return false
	}

	return hasRequiredPermissions(maxMatchNode.permissions, requiredPermissions)
}

// hasRequiredPermissions 检查是否具有所有必需的权限
func hasRequiredPermissions(nodePermissions, requiredPermissions []string) bool {
	permSet := make(map[string]bool)
	for _, perm := range nodePermissions {
		permSet[perm] = true
	}
	for _, reqPerm := range requiredPermissions {
		if !permSet[reqPerm] {
			return false
		}
	}
	return true
}

// main 主函数展示如何使用前缀树进行权限校验
func main() {
	trie := NewTrie()
	trie.Insert("/api/v1/agent/:service/info", []string{"admin"})
	trie.Insert("/api/user/delete/:id", []string{"admin", "superuser"})

	testCases := []struct {
		path                string
		requiredPermissions []string
		expectedResult      bool
	}{
		{"/api/v1/agent/abc/aaaa", []string{"admin"}, true}, // should be false
		{"/api/user/delete/123", []string{"admin"}, true},
		{"/api/user/delete/123", []string{"superuser"}, true},
		{"/api/user/delete/123", []string{"guest"}, false},
	}

	for _, tc := range testCases {
		result := trie.CheckPermissions(tc.path, tc.requiredPermissions)
		if result != tc.expectedResult {
			fmt.Printf("Test case failed: %s, required permissions: %v, expected result: %v, actual result: %v\n", tc.path, tc.requiredPermissions, tc.expectedResult, result)
		}
		fmt.Printf("Test case passed: %s, required permissions: %v, expected result: %v, actual result: %v\n", tc.path, tc.requiredPermissions, tc.expectedResult, result)
	}
}
