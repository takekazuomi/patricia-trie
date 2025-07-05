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
			children: make(map[byte]*Node),
		},
	}
}

// Insert キーをトライに挿入
func (t *Trie) Insert(key string) error {
	// TODO: 実装予定
	return nil
}

// Search キーがトライに存在するかを検索
func (t *Trie) Search(key string) bool {
	// TODO: 実装予定
	return false
}

// Delete キーをトライから削除
func (t *Trie) Delete(key string) error {
	// TODO: 実装予定
	return nil
}

// FindByPrefix 指定されたプレフィックスを持つすべてのキーを検索
func (t *Trie) FindByPrefix(prefix string) []string {
	// TODO: 実装予定
	return nil
}