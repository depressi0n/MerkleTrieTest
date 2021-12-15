package main

import "fmt"

func main() {
	trie := newDemoTree()
	testDatas := [][]byte{
		{'a', 'b', 'c', 'b', 'c'},
		{'a', 'c', 'b', 'c', 'b'},
		{'a', 'a', 'a', 'a', 'a'},
		{'a', 'c', 'c', 'c', 'c'},
		{'b', 'c', 'a', 'c', 'a'},
		{'b', 'c', 'b', 'c', 'b'},
		{'c', 'b', 'a', 'b', 'a'},
		{'c', 'a', 'a', 'a', 'a'},
	}
	for _, data := range testDatas {
		//tt := hashFunc(data)
		//s := bytesToString(tt[:])
		//fmt.Println(s)
		tmp := hashFunc(data)
		trie.Insert(tmp[:])
		//trie.Insert(data)
		trie.Show()
		//fmt.Println()
	}
	tmp := hashFunc(testDatas[3])
	fmt.Println(trie.Query(tmp[:]))
	fmt.Println(trie.Query([]byte{'d', 'a', 'c'}))
	for _, data := range testDatas {
		//tt := hashFunc(data)
		//s := bytesToString(tt[:])
		//fmt.Println(s)
		tmp = hashFunc(data)
		trie.Delete(tmp[:])
		//trie.Insert(data)
		trie.Show()
		//fmt.Println()
	}
	//trie.Delete(testDatas[5])
	//trie.Show()
	//cnt := 0
	//data := testDatas[0]
	//for cnt < 100 {
	//	tt := hashFunc(data)
	//	s := bytesToString(tt[:])
	//	//fmt.Println(s)
	//	trie.Insert([]byte(s[:32]))
	//	data=[]byte(s[:32])
	//	//trie.Insert(data)
	//	trie.Show()
	//	fmt.Println()
	//	cnt++
	//}
}
