package main

import "fmt"

type Node struct {
	val  int
	next *Node
}

var head *Node

// 10 20 30 40 50
func Display(key int) **Node {
	from := &head

	for curr := *from; curr != nil; curr = curr.next {
		if curr.val == key {
			return from
		}
		from = &curr.next
	}
	return nil
}

func main() {
	size := 8 - 1
	for i := range 100 {
		fmt.Println(i, i&size)
	}
}
