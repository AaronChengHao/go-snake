package util

import (
	"image/color"
)

type Node struct {
	X      int64
	Y      int64
	OldX   int64
	OldY   int64
	Child  *Node
	Parent *Node
	Color  color.Color
}

func AddNode(rootNode *Node) bool {
	if rootNode.Child != nil {
		return AddNode(rootNode.Child)
	} else {
		child := &Node{Color: color.Black}
		child.X = rootNode.OldX
		child.Y = rootNode.OldY
		child.Parent = rootNode
		rootNode.Child = child
	}

	return true
}
