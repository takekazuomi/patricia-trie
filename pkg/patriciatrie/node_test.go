package patriciatrie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	label := "test"
	node := NewNode(label)
	
	assert.NotNil(t, node)
	assert.Equal(t, label, node.label)
	assert.NotNil(t, node.children)
	assert.Equal(t, 0, len(node.children))
	assert.False(t, node.isEndOfKey)
	assert.Nil(t, node.value)
}

func TestNode_HasChild(t *testing.T) {
	node := NewNode("test")
	
	// 子ノードが存在しない場合
	assert.False(t, node.HasChild('a'))
	
	// 子ノードを追加
	child := NewNode("child")
	node.AddChild('a', child)
	
	// 子ノードが存在する場合
	assert.True(t, node.HasChild('a'))
	assert.False(t, node.HasChild('b'))
}

func TestNode_GetChild(t *testing.T) {
	node := NewNode("test")
	child := NewNode("child")
	node.AddChild('a', child)
	
	// 存在する子ノードを取得
	result, exists := node.GetChild('a')
	assert.True(t, exists)
	assert.Equal(t, child, result)
	
	// 存在しない子ノードを取得
	result, exists = node.GetChild('b')
	assert.False(t, exists)
	assert.Nil(t, result)
}

func TestNode_AddChild(t *testing.T) {
	node := NewNode("test")
	child := NewNode("child")
	
	// 子ノードを追加
	node.AddChild('a', child)
	
	assert.Equal(t, 1, len(node.children))
	assert.Equal(t, child, node.children['a'])
}

func TestNode_RemoveChild(t *testing.T) {
	node := NewNode("test")
	child := NewNode("child")
	node.AddChild('a', child)
	
	assert.Equal(t, 1, len(node.children))
	
	// 子ノードを削除
	node.RemoveChild('a')
	
	assert.Equal(t, 0, len(node.children))
	assert.False(t, node.HasChild('a'))
}

func TestNode_IsLeaf(t *testing.T) {
	node := NewNode("test")
	
	// 子ノードがない場合は葉ノード
	assert.True(t, node.IsLeaf())
	
	// 子ノードを追加
	child := NewNode("child")
	node.AddChild('a', child)
	
	// 子ノードがある場合は葉ノードではない
	assert.False(t, node.IsLeaf())
}

func TestNode_ChildrenCount(t *testing.T) {
	node := NewNode("test")
	
	// 初期状態では子ノードは0個
	assert.Equal(t, 0, node.ChildrenCount())
	
	// 子ノードを追加
	child1 := NewNode("child1")
	child2 := NewNode("child2")
	node.AddChild('a', child1)
	node.AddChild('b', child2)
	
	assert.Equal(t, 2, node.ChildrenCount())
}