package transform

import (
	"fmt"
	"strings"
	"testing"
)

// 定义节点结构
type Node struct {
	ID       int
	Name     string
	ParentID int
	SubNodes []*Node
}

// 实现 TreeNode 接口
func (n *Node) GetID() int {
	return n.ID
}

func (n *Node) GetParentID() int {
	return n.ParentID
}

func (n *Node) AddChild(child *Node) {
	n.SubNodes = append(n.SubNodes, child) // 类型断言
}

// 打印树结构
func printTree(node *Node, level int) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", level*2), node.Name)
	for _, child := range node.SubNodes {
		printTree(child, level+1)
	}
}

func TestBuildTree(t *testing.T) {
	nodes := []*Node{
		{ID: 6, Name: "部门B1", ParentID: 5},
		{ID: 3, Name: "部门A1", ParentID: 2},
		{ID: 4, Name: "部门A2", ParentID: 2},
		{ID: 7, Name: "部门B2", ParentID: 6},
		{ID: 1, Name: "总部", ParentID: 0},
		{ID: 5, Name: "分公司B", ParentID: 1},
		{ID: 2, Name: "分公司A", ParentID: 1},
	}
	roots, err := BuildTree(nodes)
	fmt.Println("-------", roots, err)
	// 打印结果
	for i := range roots {
		printTree(roots[i], 0)
	}
}
