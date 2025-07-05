// Package patriciatrie パトリシアトライの実装を提供
package patriciatrie

// Trie パトリシアトライの構造体
type Trie struct {
	root *Node
}

// New 新しいパトリシアトライを作成
func New() *Trie {
	return &Trie{
		root: &Node{
			label:      "",
			isEndOfKey: false,
			children:   make(map[byte]*Node),
			value:      nil,
		},
	}
}

// Insert キーをトライに挿入
func (t *Trie) Insert(key string) error {
	if key == "" {
		t.root.isEndOfKey = true

		return nil
	}

	return t.insertNode(t.root, key)
}

// Search キーがトライに存在するかを検索
func (t *Trie) Search(key string) bool {
	if key == "" {
		return t.root.isEndOfKey
	}

	return t.searchNode(t.root, key)
}

// Delete キーをトライから削除
func (t *Trie) Delete(key string) error {
	if key == "" {
		t.root.isEndOfKey = false

		return nil
	}

	return t.deleteNode(t.root, key)
}

// FindByPrefix 指定されたプレフィックスを持つすべてのキーを検索
func (t *Trie) FindByPrefix(prefix string) []string {
	var result []string
	t.findKeysWithPrefix(t.root, "", prefix, &result)

	return result
}

// insertNode 指定されたノードから始まってキーを挿入
func (t *Trie) insertNode(node *Node, key string) error {
	if len(key) == 0 {
		node.isEndOfKey = true

		return nil
	}

	firstByte := key[0]

	// 子ノードが存在しない場合、新しいノードを作成
	if !node.HasChild(firstByte) {
		newNode := NewNode(key)
		newNode.isEndOfKey = true
		node.AddChild(firstByte, newNode)

		return nil
	}

	// 子ノードが存在する場合
	child, _ := node.GetChild(firstByte)

	// 共通プレフィックスの長さを計算
	commonLen := t.findCommonPrefixLength(child.label, key)

	if commonLen == len(child.label) {
		// 子ノードのラベルが完全にマッチする場合
		// 残りのキーで再帰的に挿入
		remaining := key[commonLen:]

		return t.insertNode(child, remaining)
	}

	if commonLen == len(key) {
		// 挿入するキーが既存ノードのプレフィックスの場合
		// 既存ノードを分割して新しい中間ノードを作成
		return t.splitNode(node, child, firstByte, commonLen)
	}

	// 部分的にマッチする場合、ノードを分割
	return t.splitNodeWithNewBranch(node, child, firstByte, key, commonLen)
}

// findCommonPrefixLength 2つの文字列の共通プレフィックスの長さを計算
func (t *Trie) findCommonPrefixLength(s1, s2 string) int {
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}

	for i := range minLen {
		if s1[i] != s2[i] {
			return i
		}
	}

	return minLen
}

// splitNode 既存ノードを分割（挿入キーが既存ラベルのプレフィックス）
func (t *Trie) splitNode(parent *Node, child *Node, firstByte byte, commonLen int) error {
	// 新しい中間ノードを作成
	intermediateNode := NewNode(child.label[:commonLen])
	intermediateNode.isEndOfKey = true

	// 既存の子ノードのラベルを短縮
	remainingLabel := child.label[commonLen:]
	child.label = remainingLabel

	// 中間ノードに既存の子ノードを接続
	if len(remainingLabel) > 0 {
		intermediateNode.AddChild(remainingLabel[0], child)
	}

	// 親ノードに中間ノードを接続
	parent.children[firstByte] = intermediateNode

	return nil
}

// searchNode 指定されたノードから始まってキーを検索
func (t *Trie) searchNode(node *Node, key string) bool {
	if len(key) == 0 {
		return node.isEndOfKey
	}

	firstByte := key[0]

	// 対応する子ノードが存在しない場合、キーは存在しない
	if !node.HasChild(firstByte) {
		return false
	}

	child, _ := node.GetChild(firstByte)

	// 子ノードのラベルと比較
	if len(key) < len(child.label) {
		// キーが子ノードのラベルより短い場合、完全一致しない
		return false
	}

	// ラベルがキーのプレフィックスとして一致しない場合
	if key[:len(child.label)] != child.label {
		return false
	}

	// ラベルが完全に一致する場合、残りのキーで再帰的に検索
	remaining := key[len(child.label):]

	return t.searchNode(child, remaining)
}

// splitNodeWithNewBranch ノードを分割して新しい分岐を作成
func (t *Trie) splitNodeWithNewBranch(parent *Node, child *Node, firstByte byte, key string, commonLen int) error {
	// 共通部分で中間ノードを作成
	intermediateNode := NewNode(key[:commonLen])

	// 既存の子ノードのラベルを更新
	childRemainingLabel := child.label[commonLen:]
	child.label = childRemainingLabel

	// 新しいノードを作成
	newRemainingKey := key[commonLen:]
	newNode := NewNode(newRemainingKey)
	newNode.isEndOfKey = true

	// 中間ノードに両方の子を接続
	if len(childRemainingLabel) > 0 {
		intermediateNode.AddChild(childRemainingLabel[0], child)
	}

	if len(newRemainingKey) > 0 {
		intermediateNode.AddChild(newRemainingKey[0], newNode)
	}

	// 親ノードに中間ノードを接続
	parent.children[firstByte] = intermediateNode

	return nil
}

// deleteNode 指定されたノードから始まってキーを削除
func (t *Trie) deleteNode(node *Node, key string) error {
	if len(key) == 0 {
		// キーが完全に一致した場合、終端フラグを無効化
		node.isEndOfKey = false

		return nil
	}

	firstByte := key[0]

	// 対応する子ノードが存在しない場合、削除対象なし
	if !node.HasChild(firstByte) {
		return nil // キーが存在しないが、エラーではない
	}

	child, _ := node.GetChild(firstByte)

	// ラベルがキーのプレフィックスとして一致しない場合
	if len(key) < len(child.label) || key[:len(child.label)] != child.label {
		return nil // キーが存在しない
	}

	// ラベルが完全に一致する場合、残りのキーで再帰的に削除
	remaining := key[len(child.label):]

	err := t.deleteNode(child, remaining)
	if err != nil {
		return err
	}

	// 削除後、子ノードが不要になった場合の整理
	return t.cleanupAfterDelete(node, child, firstByte)
}

// cleanupAfterDelete 削除後のノード整理
func (t *Trie) cleanupAfterDelete(parent *Node, child *Node, firstByte byte) error {
	// 子ノードが終端でなく、子も持たない場合は削除
	if !child.isEndOfKey && child.ChildrenCount() == 0 {
		parent.RemoveChild(firstByte)

		return nil
	}

	// 子ノードが終端でなく、子を1つだけ持つ場合は圧縮
	if !child.isEndOfKey && child.ChildrenCount() == 1 {
		// 唯一の孫ノードを取得
		var grandchild *Node
		for _, v := range child.children {
			grandchild = v

			break
		}

		// 子ノードのラベルと孫ノードのラベルを結合
		combinedLabel := child.label + grandchild.label
		grandchild.label = combinedLabel

		// 親ノードに孫ノードを直接接続
		parent.children[firstByte] = grandchild
	}

	return nil
}

// findKeysWithPrefix プレフィックスマッチングでキーを検索
func (t *Trie) findKeysWithPrefix(node *Node, currentKey, prefix string, result *[]string) {
	// 現在のキーが指定されたプレフィックスで始まる場合
	if node.isEndOfKey && len(currentKey) >= len(prefix) &&
		(prefix == "" || currentKey[:len(prefix)] == prefix) {
		*result = append(*result, currentKey)
	}

	// 子ノードを探索
	for _, child := range node.children {
		newKey := currentKey + child.label
		// プレフィックスの可能性がある場合のみ再帰
		if len(newKey) >= len(prefix) || len(prefix) >= len(newKey) {
			t.findKeysWithPrefix(child, newKey, prefix, result)
		}
	}
}
