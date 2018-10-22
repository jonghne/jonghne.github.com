package main

import (
	"fmt"
	"strings"
)


type Node struct {
	count int
	no []byte
	child map[string]*Node
	parent *Node
}

var root *Node = nil
func makeDB() *Node {
	if root == nil {
		root = new(Node)
		root.count = 0
		root.no = []byte{}
		root.child = make(map[string]*Node)
	}
	return root
}

func createNode(key string) *Node {
	node := new(Node)
	node.count = 1
	node.child = make(map[string]*Node)
	node.no = []byte(key)
	return node
}

func (node *Node) setCount(v int) {
	node.count = v
}

func (node *Node) findRoom(adr []byte) (prefix []byte, ok bool) {
	l1 := len(adr)
	l2 := len(node.no)
	l := l2
	if l1 < l2 {
		l = l1
	}

	for i:=0; i<l; i++ {
		if adr[i] == node.no[i] {
			prefix = append(prefix, adr[i])
		} else {
			break
		}
	}

	ok = false
	if len(prefix) > 0 {
		ok = true
	}

	return
}

func (node *Node) removeChild(adr string) {
	delete(node.child, adr)
}

func (node *Node) addChild(adr string, nc *Node) {
	node.child[adr] = nc
	node.count++
	nc.parent = node
}

func (node *Node) searchChild(adr []byte) {
	found := false
	if len(node.child) > 0 {
		for _, n := range node.child {
			if n.no[0] == adr[0] {
				//find next room
				found = true
				n.addNode(adr)
				n.updateCount()
				break
			}
		}
	}
	if !found {
		//no child, new node then
		newNode := createNode(string(adr))
		node.addChild(string(adr), newNode)
	}
}

func (node *Node) addNode(w []byte) {
	if prefix, ok := node.findRoom(w); ok {
		//find address
		if len(prefix)==len(node.no) && len(prefix)==len(w) {
			//same thing
		} else if len(prefix)==len(node.no) && len(prefix) < len(w) {
			node.searchChild(w[len(prefix):])
		} else if len(prefix)<len(node.no) {
			//split happen
			split := createNode(string(prefix))
			split.setCount(node.count)
			//parent-->split-->node
			node.parent.addChild(string(prefix), split)
			node.parent.removeChild(string(node.no))
			//node address cut off prefix
			node.no = node.no[len(prefix):]
			split.addChild(string(node.no), node)
			if len(prefix) < len(w) {
				//add new node
				wordLeft := string(w[len(prefix):])
				newNode := createNode(wordLeft)
				split.addChild(wordLeft, newNode)
			}
		}
	}

}

func (node *Node) updateCount() int {
	if len(node.child) == 0 {
		return node.count
	}
	sum := 0
	for _, n := range node.child {
		cnt := n.updateCount()
		sum = sum + cnt
	}
	node.count = sum
	return sum
}

func insert(resident string) {
	if root == nil {
		makeDB()
	}
	root.searchChild([]byte(resident))
}

func (node *Node) findResident(w []byte) (*Node, bool) {
	if prefix, ok := node.findRoom(w); ok {
		if len(prefix)==len(node.no) && len(prefix) == len(w) {
			//matched
			return node, true
		} else if len(prefix)==len(node.no) && len(prefix) < len(w) {
			left := w[len(prefix):]
			for _, n := range node.child {
				position, flag := n.findResident(left)
				if flag {
					return position, true
				}
			}

		}
	}

	return nil, false

}

func find(word string) (node *Node, ok bool) {
	ok = false
	node = nil
	if root != nil && len(root.child) > 0 {
		for _, n := range root.child {
			node, ok = n.findResident([]byte(word))
			if ok {
				break
			}
		}
	}
	return node, ok
}

func (node *Node) printNode(level int) {
	fmt.Println(strings.Repeat(" ", level), string(node.no), " ", node.count)
	for _, n := range node.child {
		n.printNode(level+1)
	}
}

func printAll() {
	root.printNode(0)
}

func main() {
	insert("apple")
	insert("brother")
	insert("bro")
	insert("interesting")
	printAll()
	insert("interest")
	printAll()
	insert("interested")
	insert("inside")
	insert("in")
	printAll()

	n, err := find("in")
	if err {
		fmt.Println("---------")
		n.printNode(0)
	}

	n, err = find("he")
	if !err {
		fmt.Println("not found")
	}
}