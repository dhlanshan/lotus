package transform

import (
	"errors"
	"reflect"
)

// TreeNode 是通用节点接口，任意节点需实现以下方法
type TreeNode[T any] interface {
	GetID() int       // 返回当前节点的 ID
	GetParentID() int // 返回当前节点的父节点 ID
	AddChild(child T) // 添加子节点
}

// BuildTree 是通用的树构建方法
func BuildTree[T TreeNode[T]](nodes []T) ([]T, error) {
	if !(reflect.TypeOf(nodes[0]).Kind() == reflect.Ptr) {
		return nil, errors.New("nodes must be a pointer")
	}
	// ID 映射表
	nodeMap := make(map[int]T, len(nodes))
	var roots []T

	// 将所有节点添加到映射表中
	for i := range nodes {
		nodeMap[nodes[i].GetID()] = nodes[i]
	}

	// 构建父子关系
	for _, node := range nodeMap {
		if node.GetParentID() == 0 {
			roots = append(roots, node)
		} else {
			if parent, exists := nodeMap[node.GetParentID()]; exists {
				parent.AddChild(node)
			}
		}
	}
	return roots, nil
}
