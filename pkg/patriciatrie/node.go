package patriciatrie

// Node パトリシアトライのノード構造体
type Node struct {
	// エッジラベル（パス圧縮された文字列）
	label string
	
	// 子ノードのマップ（最初の文字をキーとする）
	children map[byte]*Node
	
	// このノードがキーの終端かどうか
	isEndOfKey bool
	
	// 値（必要に応じて）
	value interface{}
}

// NewNode 新しいノードを作成
func NewNode(label string) *Node {
	return &Node{
		label:    label,
		children: make(map[byte]*Node),
	}
}

// HasChild 指定されたバイトで始まる子ノードが存在するかチェック
func (n *Node) HasChild(b byte) bool {
	_, exists := n.children[b]
	return exists
}

// GetChild 指定されたバイトで始まる子ノードを取得
func (n *Node) GetChild(b byte) (*Node, bool) {
	child, exists := n.children[b]
	return child, exists
}

// AddChild 子ノードを追加
func (n *Node) AddChild(b byte, child *Node) {
	n.children[b] = child
}

// RemoveChild 子ノードを削除
func (n *Node) RemoveChild(b byte) {
	delete(n.children, b)
}

// IsLeaf このノードが葉ノードかどうかチェック
func (n *Node) IsLeaf() bool {
	return len(n.children) == 0
}

// ChildrenCount 子ノードの数を取得
func (n *Node) ChildrenCount() int {
	return len(n.children)
}