package lib

import (
	"container/list"
)

type NodeList struct {
	items list.List
}

type Node struct {
	kind   string
	Intval int
	Strval string
	List   *NodeList
}

// --------------------------------

func Node_newNil() *Node {
	return &(Node{kind: "nil"})
}

func Node_newInt(n int) *Node {
	return &(Node{kind: "int", Intval: n})
}

func Node_newStr(s string) *Node {
	return &(Node{kind: "str", Strval: s})
}

func Node_newList(nodes *NodeList) *Node {
	return &(Node{kind: "list", List: nodes})
}

func (self *Node) KindEq(kind string) bool {
	return self.kind == kind
}

// --------------------------------

func NodeList_new() *NodeList {
	var nodeList NodeList
	return &nodeList
}

func NodeList_empty() *NodeList {
	return NodeList_new()
}

func (self *NodeList) Add(item *Node) {
	self.items.PushBack(item)
}

func (self *NodeList) AddInt(n int) {
	self.Add(Node_newInt(n))
}

func (self *NodeList) NodeList_addStr(str string) {
	self.Add(Node_newStr(str))
}

func (self *NodeList) AddStr(str string) {
	self.Add(Node_newStr(str))
}

func (self *NodeList) AddList(list *NodeList) {
	self.Add(Node_newList(list))
}

func (self *NodeList) AddListAll(list *NodeList) {
	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)
		self.Add(item)
	}
}

func (self *NodeList) Len() int {
	return self.items.Len()
}

func (self *NodeList) Get(index int) *Node {
	i := 0
	for el := self.items.Front(); el != nil; el = el.Next() {
		if i == index {
			node, _ := el.Value.(*Node)
			return node
		}
		i++
	}
	panic("Invalid index")
}

// --------------------------------

type Names struct {
	list list.List
}

func Names_new() *Names {
	return new(Names)
}

func Names_empty() *Names {
	return Names_new()
}

func (self *Names) Add(name string) {
	self.list.PushBack(name)
}

func (self *Names) IndexOf(name string) int {
	i := 0
	for el := self.list.Front(); el != nil; el = el.Next() {
		if el.Value == name {
			return i
		}
		i++
	}
	return -1
}
