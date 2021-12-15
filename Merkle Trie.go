package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const HASHSIZE int = 32
const SIZE int = 32
const DEBUGSIZE int = 5

// 0:F  ->  0~16
var m map[byte]int = map[byte]int{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'a': 10,
	'b': 11,
	'c': 12,
	'd': 13,
	'e': 14,
	'f': 15,
}

var hashFunc func(data []byte) [HASHSIZE]byte = sha256.Sum256

type Hash [SIZE]byte

func NewHash(b []byte) Hash {
	if len(b) != SIZE {
		return [SIZE]byte{}
	}
	res := *new(Hash)
	copy(res[:], b[:])
	return res
}
func (hash Hash) String() string {
	for i := 0; i < HASHSIZE/2; i++ {
		hash[i], hash[HASHSIZE-1-i] = hash[HASHSIZE-1-i], hash[i]
	}
	return hex.EncodeToString(hash[:])
}

func newNode(s string) *node {
	return &node{
		childrenNum: 0,
		cnt:         make([]int, 16),
		children:    make([]*node, 16),
		label:       *new(Hash),
		//value:       value,
		valueStr: s,
	}
}

type node struct {
	// 节点所在高度
	childrenNum int
	// 有效节点计数
	cnt []int
	// 孩子节点
	children []*node
	// 计算并存储下一层的hash值
	label Hash
	// 压缩表示
	//value    []byte
	valueStr string
}

func (n *node) UpdateLabel() {
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

type demoTree struct {
	// 根节点
	root *node
	// 用于记录树中全部有效节点有多少个
	cnt int
	//
}

func newDemoTree() *demoTree {
	return &demoTree{
		root: newNode(""),
		cnt:  0,
	}
}

// Query 查询指定内容是否存在，如果存在则返回true，否则返回false
func (t *demoTree) Query(content []byte) bool {
	if len(content) != SIZE {
		return false
	}
	str := NewHash(content[:]).String()
	p := t.root.children[m[str[0]]]
	if p == nil {
		return false
	}
	length := 0
	for length < len(str) {
		loc := 0
		for loc < len(p.valueStr) && length+loc < len(str) && p.valueStr[loc] == str[length+loc] {
			loc++
		}
		// 不存在的content
		if loc != len(p.valueStr) {
			return false
		}
		length += len(p.valueStr)
		if length == len(str) {
			break
		}
		p = p.children[m[str[length]]]
	}
	return true
}

// Delete 删除一个有效节点，如果删除成功，返回true，否则返回false
// 如果成功则需要更新计数
func (t *demoTree) Delete(content []byte) bool {
	if len(content) != SIZE {
		return false
	}
	str := NewHash(content).String()
	p := t.root.children[m[str[0]]]
	if p == nil {
		return false
	}
	// 构建一个栈，在删除后进行合并
	s := make([]*node, 0, SIZE*2+1)
	s = append(s, t.root)
	s = append(s, p)
	cur := 1
	length := 0
	for length < len(content) {
		loc := 0
		for loc < len(s[cur].valueStr) && length+loc < len(str) && s[cur].valueStr[loc] == str[length+loc] {
			loc++
		}
		// 不存在的content
		if loc != len(s[cur].valueStr) {
			return false
		}
		length += len(s[cur].valueStr)
		if length == len(str) {
			break
		}
		s = append(s, s[cur].children[m[str[length]]])
		cur++
	}
	// 完成删除
	s[len(s)-2].children[m[s[len(s)-1].valueStr[0]]] = nil
	s[len(s)-2].cnt[m[s[len(s)-1].valueStr[0]]]--
	if s[len(s)-2].cnt[m[s[len(s)-1].valueStr[0]]] == 0 {
		s[len(s)-2].childrenNum--
		if s[len(s)-2].childrenNum == 1 && s[len(s)-2] != t.root {
			// 删除掉之后导致节点只有一个孩子
			// 找到唯一的孩子，然后合并value
			j := 0
			for ; j < len(s[len(s)-2].children) && s[len(s)-2].children[j] == nil; j++ {
				// pass
			}
			s[len(s)-2].children[j].valueStr = s[len(s)-2].valueStr + s[len(s)-2].children[j].valueStr
			s[len(s)-3].children[m[s[len(s)-2].valueStr[0]]] = s[len(s)-2].children[j]
		}
	}
	// 去除叶子节点
	s = s[:len(s)-1]
	// 更新非叶节点
	for len(s) > 0 {
		s[len(s)-1].UpdateLabel()
		s = s[:len(s)-1]
	}
	t.cnt--
	return true
}

// Insert 向一个子树中添加子节点
func (t *demoTree) Insert(content []byte) {
	if len(content) != SIZE {
		return
	}
	str := NewHash(content).String()
	// 记录路径，最多32的长度
	s := make([]*node, 0, 32)
	// 从根节点开始
	p := t.root
	// 最后需要更新路径上的标签
	s = append(s, p)
	// 若初始时树为空
	if p.children[m[str[0]]] == nil {
		// 根节点的digest等于content的digest
		digest := hashFunc(content)
		p.children[m[str[0]]] = newNode(str)
		p.cnt[m[str[0]]] = 1
		copy(p.children[m[str[0]]].label[:], digest[:])
		p.childrenNum++
		// 计算根节点的digest
		p.UpdateLabel()
		return
	}
	p = p.children[m[str[0]]]
	s = append(s, p)
	length, loc := 0, 0
	for length < len(str) {
		// 比较当前节点中value和content的内容
		for ; length+loc < len(str) && loc < len(p.valueStr) && str[length+loc] == p.valueStr[loc]; loc++ {
			// pass
		}
		// 完全相同，可以进入下一层
		if loc == len(p.valueStr) {
			length += loc
			s = append(s, p)
			if length == len(str) {
				// 说明已经存在了
				break
			}
			// 第一个
			if p.children[m[str[length]]] == nil {
				p.children[m[str[length]]] = newNode(str[length:])
				p.cnt[m[str[length]]] = 1
				p.childrenNum++
				tmp := hashFunc(content)
				copy(p.children[m[str[length]]].label[:], tmp[:])
				for len(s) > 0 {
					s[len(s)-1].UpdateLabel()
					s = s[:len(s)-1]
				}
				return
			}
			p = p.children[m[str[length]]]
			// 重置
			loc = 0
			continue
		} else {
			// 不完全匹配，需要对p完成分裂
			// 原来的的节点下面如果只有一个孩子，甚至已经是叶节点
			if p.childrenNum == 0 {
				// 叶节点分裂成非叶节点
				tmp := newNode(p.valueStr[loc:])
				tmp.label = p.label
				p.children[m[p.valueStr[loc]]] = tmp
				p.cnt[m[p.valueStr[loc]]] = 1
				p.childrenNum = 1
				s = append(s, p)
			} else {
				// 创建一个新节点接管孩子
				p.children[m[p.valueStr[loc]]] = newNode(p.valueStr[loc:])
				p.children[m[p.valueStr[loc]]].children = p.children
				p.children[m[p.valueStr[loc]]].childrenNum = p.childrenNum
				p.children[m[p.valueStr[loc]]].cnt = p.cnt
				p.children[m[p.valueStr[loc]]].label = p.label
				// 置空
				for i := 0; i != m[p.valueStr[loc]]; i++ {
					p.children[i] = nil
					p.cnt[i] = 0
				}
				for i := m[p.valueStr[loc]]; i < len(p.children); i++ {
					p.children[i] = nil
					p.cnt[i] = 0
				}
				p.childrenNum = 1
				s = append(s, p)
			}
			// 修改当前节点的value
			p.valueStr = p.valueStr[:loc]
			length += loc
			p.children[m[str[length]]] = newNode(str[length:])
			p.cnt[m[str[length]]] = 1
			p.childrenNum++
			for len(s) > 0 {
				s[len(s)-1].UpdateLabel()
				s = s[:len(s)-1]
			}
			t.cnt++
			return
			//for i := len(s) - 1; i >= 0; i-- {
			//	b := bytes.Buffer{}
			//	for j := 0; j < len(s[i].children); j++ {
			//		if s[i].children[j] != nil {
			//			b.Write(s[i].children[j].label[:])
			//		}
			//	}
			//	digest := hashFunc(b.Bytes())
			//	copy(s[i].label[:], digest[:])
			//}
		}
	}
}
func (t *demoTree) Show() {
	var dfs func(*node)
	dfs = func(cur *node) {
		if cur == nil {
			return
		}
		fmt.Printf("(%s", cur.valueStr)
		for i := 0; i < len(cur.children); i++ {
			if cur.children[i] != nil {
				dfs(cur.children[i])
			}
		}
		fmt.Printf(")")
	}
	dfs(t.root)
	fmt.Println()
}

//func init() {
//	m = make(map[byte]int)
//	for i := 0; i < 10; i++ {
//		m[byte('0'+i)] = i
//	}
//	for i := 0; i < 26; i++ {
//		m[byte('a'+i)] = 10 + i
//	}
//}
