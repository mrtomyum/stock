package model

type Location struct {
	//ID       uint64  `json:"id"`
	//Base
	//ParentID uint64  `json:"parentID"`
	Node
	Name     string  `json:"name"`
	Type     LocType `json:"type"`
}



// Todo: Make Node{} to be "Base struct" for Tree-able struct.
type Node struct {
	ID       uint64  `json:"-"`
	ParentID uint64  `json:"-"`
	Child    []*Node `json:"nodes,omitempty"`
}

func (this *Node) Size() int {
	var size int = len(this.Child)
	for _, c := range this.Child {
		size += c.Size()
	}
	return size
}

func (this *Node) Add(nodes ...*Node) bool {
	var size = this.Size()
	for _, node := range nodes {
		if node.ParentID == this.ID {
			this.Child = append(this.Child, node)
		} else {
			for _, c := range this.Child {
				if c.Add(node) {
					break
				}
			}
		}
	}
	return this.Size() == size + len(nodes)
}

