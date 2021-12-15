package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

type MerkleTrieNode struct {
	depth    int
	maxDepth int

	childNum int
	children []*MerkleTrieNode
	childCnt []int

	label []byte
	value string
}

func (n *MerkleTrieNode) UpdateLabel() {
	b := bytes.Buffer{}
	zeroLabel := make([]byte, SIZE)
	for i := 0; i < len(n.children); i++ {
		if n.children[i] == nil {
			b.Write(zeroLabel)
		} else {
			b.Write(n.children[i].label[:])
		}
	}
	t := hashFunc(b.Bytes())
	copy(n.label[:], t[:])
	return
}

func (r *MerkleTrieNode) Insert(s string, start int) bool {
	pos := 0
	for pos < len(r.value) && r.value[pos] == s[start+pos] {
		pos++
	}
	if pos == len(r.value) {
		if start+pos == len(s) {
			// 重复
			return false
		}
		if r.children[m[s[start+pos]]] == nil {
			content, _ := hex.DecodeString(s)
			tmp := hashFunc(content)

			r.children[m[s[start+pos]]] = &MerkleTrieNode{
				depth:    0,
				maxDepth: 0,
				childNum: 0,
				children: make([]*MerkleTrieNode, 16),
				childCnt: make([]int, 16),
				label:    make([]byte, 32),
				value:    s[start+pos:],
			}
			copy(r.children[m[s[start+pos]]].label[:], tmp[:])
			r.childNum++
			r.childCnt[m[s[start+pos]]]++
		} else if r.children[m[s[start+pos]]].Insert(s, start+pos) {
			r.childCnt[m[s[start+pos]]] += 1
		}
		r.maxDepth = max(r.maxDepth, r.children[m[s[start+pos]]].depth)
		r.depth = r.maxDepth + 1
		r.UpdateLabel()
		return true
	}
	// 分裂
	splitNode := &MerkleTrieNode{
		depth:    r.depth,
		maxDepth: r.maxDepth,
		childNum: r.childNum,
		children: r.children,
		childCnt: r.childCnt,
		label:    r.label,
		value:    r.value[pos:],
	}
	r.depth = splitNode.depth + 1
	r.maxDepth = splitNode.depth
	r.childNum = 1
	r.children = make([]*MerkleTrieNode, 16)
	r.children[m[r.value[pos]]] = splitNode
	r.childCnt = make([]int, 16)
	r.childCnt[m[r.value[pos]]] = 1
	r.value = r.value[:pos]

	content, _ := hex.DecodeString(s)
	tmp := hashFunc(content)
	insertNode := &MerkleTrieNode{
		depth:    0,
		maxDepth: 0,
		childNum: 0,
		children: make([]*MerkleTrieNode, 16),
		childCnt: make([]int, 16),
		label:    make([]byte, 32),
		value:    s[start+pos:],
	}
	copy(insertNode.label[:], tmp[:])
	r.children[m[s[start+pos]]] = insertNode
	r.childCnt[m[s[start+pos]]]++
	r.childNum++

	r.UpdateLabel()
	return true
}

func (r *MerkleTrieNode) Delete(s string, start int) bool {
	p := r.children[m[s[start]]]
	if p == nil {
		return false
	}
	pos := 0
	for pos < len(p.value) && p.value[pos] == s[start+pos] {
		pos++
	}
	if pos < len(p.value) {
		return false
	}
	if start+pos == len(s) {
		// 可以删除p
		r.children[m[s[start]]] = nil
		r.childCnt[m[s[start]]]--
		if r.childCnt[m[s[start]]] == 0 {
			r.childNum--
		}
		if r.childNum == 1 {
			c := 0
			for c < len(r.children) && r.children[c] == nil {
				c++
			}
			r.value += r.children[c].value
			r.childNum = r.children[c].childNum
			r.childCnt = r.children[c].childCnt
			r.maxDepth = r.children[c].maxDepth
			r.depth = r.children[c].depth
			r.children = r.children[c].children
		}
		r.UpdateLabel()
		return true
	}
	if p.Delete(s, start+pos) {
		// 删除成功
		r.UpdateLabel()
		return true
	}
	return false
}

type newnode struct {
	// 孩子的最大深度
	maxDepth int
	// 自己的深度
	depth int

	// 孩子个数
	childrenNum int
	// 有效节点计数
	cnt []int
	// 孩子节点
	children []*newnode
	// 计算并存储下一层的hash值
	label Hash
	// 压缩表示
	//value    []byte
	valueStr string
}

// Insert 递归更新
func (n *newnode) Insert(content string) bool {
	pos := 0
	for pos < len(n.valueStr) && content[pos] == n.valueStr[pos] {
		pos++
	}
	if pos == len(n.valueStr) {
		if pos == len(content) {
			// 重复元素，无需更新
			fmt.Println("Duplicated element")
			return false
		}
		// 说明可以进入下一层的节点
		child := m[content[len(n.valueStr)]]
		// 如果没有相应的下一层
		if n.children[child] == nil {
			n.children[child] = newnewNode(content[len(n.valueStr):])
			n.cnt[child] = 1
			n.maxDepth = max(n.maxDepth, 1)
			n.depth = n.maxDepth + 1
			n.childrenNum++
		} else {
			if n.children[child].Insert(content[len(n.valueStr):]) {
				n.cnt[child] += 1
				n.maxDepth = max(n.maxDepth, n.children[child].depth+1)
				n.depth = n.maxDepth + 1
			}
		}
		n.cnt[child]++
	} else {
		// 分裂
		splitNode := newnewNode(n.valueStr[pos:])
		splitNode.depth = n.depth
		splitNode.maxDepth = n.maxDepth
		splitNode.cnt = n.cnt
		splitNode.children = n.children
		splitNode.childrenNum = n.childrenNum
		splitNode.label = n.label

		n.valueStr = n.valueStr[:pos]
		n.children = make([]*newnode, 16)
		n.children[m[splitNode.valueStr[0]]] = splitNode
		n.cnt = make([]int, 16)
		for i := 0; i < len(splitNode.cnt); i++ {
			n.cnt[splitNode.valueStr[0]] += splitNode.cnt[i]
		}
		n.childrenNum = 1
		n.maxDepth = splitNode.maxDepth + 1
		n.depth = n.maxDepth + 1

		n.children[m[content[pos]]] = newnewNode(content[pos:])
		n.cnt[m[content[pos]]] += 1
		n.childrenNum += 1
	}
	n.UpdateLabel()
	return true
}
func (n *newnode) Delete(content string) bool {
	pos := 0
	for pos < len(n.valueStr) && content[pos] == n.valueStr[pos] {
		pos++
	}
	if pos != len(n.valueStr) {
		// 不存在
		return false
	}
	// 到了叶子节点
	if pos == len(content) {

	}
	if n.children[m[content[len(n.valueStr)]]].Delete(content[len(n.valueStr):]) {
		n.cnt[m[content[len(n.valueStr)]]]--
		// 如果相应位置上本就只剩下一个孩子，删除后，需要合并
		if n.cnt[m[content[len(n.valueStr)]]] == 0 {
			n.childrenNum--
			// 合并
			if n.childrenNum == 1 {
				// 找到这个唯一的孩子
				unique := 0
				for unique < len(n.children) && n.children[unique] == nil {
					unique++
				}
				n.valueStr += n.children[unique].valueStr
				n.depth = n.children[unique].depth
				n.maxDepth = n.children[unique].maxDepth
				n.cnt = n.children[unique].cnt
				n.children = n.children[unique].children
				n.childrenNum = n.children[unique].childrenNum
			}
		}
		n.UpdateLabel()
		return true
	}
	return false
}

func (n *newnode) UpdateLabel() {
	b := bytes.Buffer{}
	zeroLabel := make([]byte, SIZE)
	for i := 0; i < len(n.children); i++ {
		if n.children[i] == nil {
			b.Write(zeroLabel)
		} else {
			b.Write(n.children[i].label[:])
		}
	}
	t := hashFunc(b.Bytes())
	copy(n.label[:], t[:])
	return
}

func newnewNode(s string) *newnode {
	return &newnode{
		maxDepth:    0,
		depth:       0,
		childrenNum: 0,
		cnt:         make([]int, 16),
		children:    make([]*newnode, 16),
		label:       *new(Hash),
		valueStr:    s,
	}
}
