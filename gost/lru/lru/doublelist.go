package lru

import "fmt"

type Node struct {
	pre, next *Node
	key       string
	value     interface{}
}

type DoubleList struct {
	head, tail *Node
}

func NewDL() *DoubleList {
	return &DoubleList{}
}

func (d *DoubleList) RemoveLast() {
	if d.tail != nil {
		lastNode := d.tail
		d.tail = lastNode.pre
		d.tail.next = nil
		lastNode.pre = nil
	}
}

func (d *DoubleList) InsertHead(newNode *Node) *Node {
	if d.head == nil {
		d.head = newNode
		d.tail = newNode
	} else {
		newNode.next = d.head
		d.head.pre = newNode
		d.head = newNode
	}
	return newNode
}

func (d *DoubleList) MoveToHead(node *Node) {
	if node == d.head {
		return
	}

	if node != nil {
		node.pre.next = node.next
		if node.next != nil {
			node.next.pre = node.pre
		}
		node.next = d.head
		d.head.pre = node
		d.head = node
	}
}

func (d *DoubleList) Show() {
	move := d.head
	res := []interface{}{}
	for move != nil {
		temp := []interface{}{move.key, move.value}
		res = append(res, temp)
		move = move.next
	}
	fmt.Println(res...)
}

func TestDL() {
	// dl := NewDL()
	// a := dl.InsertHead("a", 1)
	// dl.InsertHead("b", "abc")
	// dl.InsertHead("c", "111")

	// dl.Show()

	// dl.MoveToHead(a)

	// dl.RemoveLast()
	// dl.Show()
	// dl.RemoveLast()
	// dl.Show()
}
