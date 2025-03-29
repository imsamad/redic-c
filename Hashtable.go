package main

import (
	"fmt"
	"hash/fnv"
)

// Hashable interface allows defining custom keys.
type Hashable interface {
	Hash() uint64
	Equals(other Hashable) bool
}

// MyNode represents a key-value pair in the hash map.
type MyNode struct {
	Key   string
	Value int
}

// Hashable Implementation for MyNode
func (mn *MyNode) Hash() uint64 {
	hasher := &FNVHasher{}
	return hasher.ComputeHash(mn.Key)
}

func (mn *MyNode) Equals(other Hashable) bool {
	o, ok := other.(*MyNode)
	return ok && mn.Key == o.Key
}

// Hasher interface for different hashing strategies.
type Hasher interface {
	ComputeHash(data string) uint64
}

// FNVHasher implements the Hasher interface using FNV-1a.
type FNVHasher struct{}

func (h *FNVHasher) ComputeHash(data string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(data))
	return hash.Sum64()
}

// HNode represents a node in the hash table.
type HNode struct {
	next  *HNode
	hCode uint64
	value Hashable
}

// NewHNode constructs a new node.
func NewHNode(value Hashable) *HNode {
	return &HNode{
		hCode: value.Hash(),
		value: value,
	}
}

// Attach a node to a given head.
func (n *HNode) Attach(head **HNode) {
	n.next = *head
	*head = n
}

// Detach a node from a given head.
func (n *HNode) Detach(head **HNode) *HNode {
	*head = n.next
	n.next = nil
	return n
}

// HTab represents a hash table with chaining for collision resolution.
type HTab struct {
	tab  []*HNode
	mask uint64
	size uint64
}

// NewHTab initializes a hash table with a given size (must be power of 2).
func NewHTab(size uint64) *HTab {
	if size < 0 || (size&(size-1)) != 0 {
		panic("size must be a power of 2")
	}
	return &HTab{
		tab:  make([]*HNode, size),
		mask: size - 1,
	}
}

// insert adds a node into the hash table.
func (ht *HTab) insert(node *HNode) {
	pos := node.hCode & ht.mask
	node.Attach(&ht.tab[pos])
	ht.size++
}

// lookup finds a node.
func (ht *HTab) lookup(key Hashable) **HNode {
	pos := key.Hash() & ht.mask
	from := &ht.tab[pos]
	for cur := *from; cur != nil; cur = cur.next {
		if cur.hCode == key.Hash() && cur.value.Equals(key) {
			return from
		}
		from = &cur.next
	}
	return nil
}

// detach removes a node.
func (ht *HTab) detach(from **HNode) *HNode {
	node := *from
	if node != nil {
		ht.size--
		return node.Detach(from)
	}
	return nil
}

// HMap represents a hash map with progressive rehashing.
type HMap struct {
	newer      *HTab
	older      *HTab
	migratePos uint64
}

// NewHMap creates a new hash map.
func NewHMap() *HMap {
	return &HMap{
		newer: NewHTab(4),
	}
}

// HelpRehashing migrates nodes progressively.
func (hm *HMap) helpRehashing() {
	const rehashingWork = 128
	var nwork uint64

	for nwork < rehashingWork && hm.older != nil && hm.older.size > 0 {
		if hm.migratePos >= uint64(len(hm.older.tab)) {
			break
		}
		from := &hm.older.tab[hm.migratePos]
		if *from == nil {
			hm.migratePos++
			continue
		}
		hm.newer.insert(hm.older.detach(from))
		nwork++
	}

	if hm.older != nil && hm.older.size == 0 {
		hm.older = nil
	}
}

// TriggerRehashing starts rehashing when needed.
func (hm *HMap) triggerRehashing() {
	hm.older = hm.newer
	hm.newer = NewHTab((hm.newer.mask + 1) * 2)
	hm.migratePos = 0
}

// Insert adds an item to the hash map.
func (hm *HMap) Insert(item Hashable) {
	hm.helpRehashing()
	if hm.newer == nil {
		hm.newer = NewHTab(4)
	}

	node := NewHNode(item)
	hm.newer.insert(node)

	const maxLoadFactor = 8
	if hm.older == nil && hm.newer.size >= (hm.newer.mask+1)*maxLoadFactor {
		hm.triggerRehashing()
	}
}

// Find retrieves an item by key.
func (hm *HMap) Find(key Hashable) Hashable {
	hm.helpRehashing()
	if from := hm.newer.lookup(key); from != nil {
		return (*from).value
	}
	if hm.older != nil {
		if from := hm.older.lookup(key); from != nil {
			return (*from).value
		}
	}
	return nil
}

// Remove deletes an item by key.
func (hm *HMap) Remove(key Hashable) Hashable {
	hm.helpRehashing()
	if from := hm.newer.lookup(key); from != nil {
		return hm.newer.detach(from).value
	}
	if hm.older != nil {
		if from := hm.older.lookup(key); from != nil {
			return hm.older.detach(from).value
		}
	}
	return nil
}

// Size returns the total number of elements.
func (hm *HMap) Size() uint64 {
	size := hm.newer.size
	if hm.older != nil {
		size += hm.older.size
	}
	return size
}

// Clear resets the hash map.
func (hm *HMap) Clear() {
	hm.newer = nil
	hm.older = nil
}
// Example Usage
func _main() {
	hmap := NewHMap()

	// Insert nodes
	hmap.Insert(&MyNode{Key: "Alice", Value: 25})
	hmap.Insert(&MyNode{Key: "Bob", Value: 30})
	hmap.Insert(&MyNode{Key: "Charlie", Value: 35})

	// Lookup
	search := &MyNode{Key: "Bob"}
	found := hmap.Find(search)
	if found != nil {
		fmt.Printf("Found: %s -> %d\n", search.Key, found.(*MyNode).Value)
	} else {
		fmt.Println("Not found")
	}

	// Delete
	hmap.Remove(search)
	found = hmap.Find(search)
	if found != nil {
		fmt.Printf("Found: %s -> %d\n", search.Key, found.(*MyNode).Value)
	} else {
		fmt.Println("Deleted successfully")
	}

	// Print size
	fmt.Printf("Size of hash map: %d\n", hmap.Size())
}
