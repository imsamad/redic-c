package main

import (
	"fmt"
	"hash/fnv"
)

const (
	MAX_LOADING_FACTOR = 2
	TAB_SIZE           = 2
)

type Node struct {
	next  *Node
	hCode uint64
	Key   string
	Value string
}

func ComputeHash(data string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(data))
	return hash.Sum64()
}

func (n *Node) Hash() uint64 {
	return ComputeHash(n.Key)
}

type HTab struct {
	tab  []*Node
	size uint64
	mask uint64
}

func NewHTab(size uint64) *HTab {
	ht := &HTab{}

	if size < 1 || ((size-1)&size) != 0 {
		panic("size must be a power of 2")
	}

	ht.size = 0
	ht.mask = size - 1
	ht.tab = make([]*Node, size)
	return ht
}

func (ht *HTab) insert(node *Node) {
	pos := ht.mask & node.hCode

	// if ht.tab[pos] != nil {
	ht.size++
	// }

	node.next = ht.tab[pos]
	ht.tab[pos] = node
}

func (ht *HTab) update(key, value string) *Node {
	from := ht.lookup(key)

	if *from == nil {
		return nil
	}
	(*from).Value = value
	return *from
}

// return the pre-decessor node
func (ht *HTab) lookup(key string) **Node {
	hash := ComputeHash(key)

	pos := ht.mask & hash
	from := &ht.tab[pos]

	for curr := *from; curr != nil; curr = curr.next {
		if hash == curr.hCode && curr.Key == key {
			return from
		}
		from = &curr.next
	}
	return from
}

func (ht *HTab) remove(key string) *Node {
	from := ht.lookup(key)

	if *from == nil {
		return nil
	}

	ht.size--
	removed := *from
	*from = removed.next
	return removed
}

type HMap struct {
	new        *HTab
	old        *HTab
	migratePos uint64
}

func NewHMap(size uint64) *HMap {
	return &HMap{
		new: NewHTab(size),
	}
}

func (hm *HMap) Insert(key, value string) {
	node := &Node{Key: key, Value: value, hCode: ComputeHash(key)}
	hm.new.insert(node)
	if hm.old == nil && hm.new.size >= MAX_LOADING_FACTOR*(hm.new.mask+1) {
		hm.triggerRehashing()
	}
}

func (hm *HMap) Find(key string) *Node {
	return *hm.new.lookup(key)
}

func (hm *HMap) Delete(key string) *Node {
	return hm.new.remove(key)
}

func (hm *HMap) triggerRehashing() {
	hm.old = hm.new
	hm.new = NewHTab((hm.new.mask + 1) * 2)
	hm.migratePos = 0
}

func (hm *HMap) Update(key string, value string) *Node {
	return hm.new.update(key, value)
}

func main() {
	hm := NewHMap(TAB_SIZE)
	hm.Insert("name", "samad")
	hm.Insert("age", "25")
	hm.Insert("sex", "male")
	hm.Insert("addr", "localhost")

	fn := hm.Find("sex1")
	if fn == nil {
		fmt.Println("not found")
	} else {
		fmt.Println(fn.Value)
	}
}
