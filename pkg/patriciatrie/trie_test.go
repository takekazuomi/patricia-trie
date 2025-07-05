package patriciatrie

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	trie := New()
	
	assert.NotNil(t, trie)
	assert.NotNil(t, trie.root)
	assert.False(t, trie.root.isEndOfKey)
	assert.Empty(t, trie.root.children)
}

func TestTrie_Insert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		keys []string
	}{
		{
			name: "基本的な挿入",
			keys: []string{"cat", "cats", "dog"},
		},
		{
			name: "空文字列",
			keys: []string{""},
		},
		{
			name: "プレフィックス関係",
			keys: []string{"a", "ab", "abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			trie := New()
			
			for _, key := range tt.keys {
				err := trie.Insert(key)
				require.NoError(t, err)
			}
		})
	}
}

func TestTrie_Search(t *testing.T) {
	t.Parallel()

	trie := New()
	keys := []string{"cat", "cats", "dog", "dogs"}
	
	// キーを挿入
	for _, key := range keys {
		err := trie.Insert(key)
		require.NoError(t, err)
	}
	
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"存在するキー: cat", "cat", true},
		{"存在するキー: cats", "cats", true},
		{"存在するキー: dog", "dog", true},
		{"存在するキー: dogs", "dogs", true},
		{"存在しないキー: ca", "ca", false},
		{"存在しないキー: elephant", "elephant", false},
		{"空文字列", "", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			result := trie.Search(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTrie_Delete(t *testing.T) {
	t.Parallel()

	trie := New()
	keys := []string{"cat", "cats", "dog"}
	
	// キーを挿入
	for _, key := range keys {
		err := trie.Insert(key)
		require.NoError(t, err)
	}
	
	// 削除テスト
	err := trie.Delete("cat")
	require.NoError(t, err)
	
	// 削除されたキーは見つからない
	assert.False(t, trie.Search("cat"))
	
	// 他のキーは残っている
	assert.True(t, trie.Search("cats"))
	assert.True(t, trie.Search("dog"))
}

func TestTrie_FindByPrefix(t *testing.T) {
	t.Parallel()

	trie := New()
	keys := []string{"cat", "cats", "dog", "dogs", "elephant"}
	
	// キーを挿入
	for _, key := range keys {
		err := trie.Insert(key)
		require.NoError(t, err)
	}
	
	tests := []struct {
		name     string
		prefix   string
		expected []string
	}{
		{
			name:     "プレフィックス: ca",
			prefix:   "ca",
			expected: []string{"cat", "cats"},
		},
		{
			name:     "プレフィックス: dog",
			prefix:   "dog",
			expected: []string{"dog", "dogs"},
		},
		{
			name:     "プレフィックス: el",
			prefix:   "el",
			expected: []string{"elephant"},
		},
		{
			name:     "マッチしないプレフィックス",
			prefix:   "xyz",
			expected: []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			result := trie.FindByPrefix(tt.prefix)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}