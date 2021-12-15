package main

import (
	"crypto/rand"
	"testing"
	"time"
)

func Test_RecursiveMerkleTrie(t *testing.T) {
	root := &MerkleTrieNode{
		depth:    0,
		maxDepth: 0,
		childNum: 0,
		children: make([]*MerkleTrieNode, 16),
		childCnt: make([]int, 16),
		label:    make([]byte, 32),
		value:    "",
	}
	//testDatas := [][]byte{
	//{'a', 'b', 'c', 'b', 'c'},
	//{'a', 'c', 'b', 'c', 'b'},
	//{'a', 'a', 'a', 'a', 'a'},
	//{'a', 'c', 'c', 'c', 'c'},
	//{'b', 'c', 'a', 'c', 'a'},
	//{'b', 'c', 'b', 'c', 'b'},
	//{'c', 'b', 'a', 'b', 'a'},
	//{'c', 'a', 'a', 'a', 'a'},
	//}

	lengths := []int{100, 1000, 10000, 100000, 200000, 500000, 600000, 800000, 1000000}
	repeat := 10
	//lengths := []int{10}
	for _, length := range lengths {
		for i := 0; i < repeat; i++ {
			t.Logf("The %d-th time for number %d", i, length)
			testDatas := make([][]byte, length)
			cur := 0
			for cur < length {
				rng := rand.Reader
				s := make([]byte, 32)
				n, err := rng.Read(s)
				if err != nil || n < 32 {
					continue
				}
				copy(testDatas[cur], s[:])
				cur++
			}
			insertDatas := make([]string, len(testDatas))
			for i := 0; i < len(testDatas); i++ {
				tmp := hashFunc(testDatas[i])
				insertDatas[i] = NewHash(tmp[:]).String()
			}
			for i := 0; i < len(insertDatas); i++ {
				start := time.Now().Unix()
				root.Insert(insertDatas[i], 0)
				end := time.Now().Unix()
				t.Logf("Element number is %d, Tree height is %d, Insert time is %d", i, root.depth, end-start)
			}
		}
	}
}
