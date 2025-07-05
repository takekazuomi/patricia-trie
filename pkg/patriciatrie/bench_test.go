package patriciatrie

import (
	"fmt"
	"math/rand"
	"testing"
)

// BenchmarkTrie_Insert 挿入操作のベンチマーク
func BenchmarkTrie_Insert(b *testing.B) {
	trie := New()
	keys := generateRandomKeys(1000)

	b.ResetTimer()
	b.ReportAllocs()

	for i := range b.N {
		key := keys[i%len(keys)]
		_ = trie.Insert(key)
	}
}

// BenchmarkTrie_Search 検索操作のベンチマーク
func BenchmarkTrie_Search(b *testing.B) {
	trie := New()
	keys := generateRandomKeys(1000)

	// 事前にキーを挿入
	for _, key := range keys {
		_ = trie.Insert(key)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := range b.N {
		key := keys[i%len(keys)]
		_ = trie.Search(key)
	}
}

// BenchmarkTrie_InsertAndSearch 挿入と検索の混合ベンチマーク
func BenchmarkTrie_InsertAndSearch(b *testing.B) {
	trie := New()
	keys := generateRandomKeys(1000)

	b.ResetTimer()
	b.ReportAllocs()

	for i := range b.N {
		key := keys[i%len(keys)]
		if i%2 == 0 {
			_ = trie.Insert(key)
		} else {
			_ = trie.Search(key)
		}
	}
}

// BenchmarkTrie_MemoryUsage メモリ使用量のベンチマーク
func BenchmarkTrie_MemoryUsage(b *testing.B) {
	keySizes := []int{10, 100, 1000, 10000}

	for _, size := range keySizes {
		b.Run(fmt.Sprintf("Keys_%d", size), func(b *testing.B) {
			keys := generateRandomKeys(size)

			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				trie := New()
				for _, key := range keys {
					_ = trie.Insert(key)
				}
			}
		})
	}
}

// BenchmarkTrie_PrefixLength プレフィックス長による性能測定
func BenchmarkTrie_PrefixLength(b *testing.B) {
	prefixLengths := []int{5, 10, 20, 50}

	for _, length := range prefixLengths {
		b.Run(fmt.Sprintf("PrefixLen_%d", length), func(b *testing.B) {
			keys := generateKeysWithCommonPrefix("prefix", length, 1000)
			trie := New()

			// 事前にキーを挿入
			for _, key := range keys {
				_ = trie.Insert(key)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := range b.N {
				key := keys[i%len(keys)]
				_ = trie.Search(key)
			}
		})
	}
}

// generateRandomKeys ランダムなキーを生成
func generateRandomKeys(count int) []string {
	keys := make([]string, count)

	const charset = "abcdefghijklmnopqrstuvwxyz"

	const maxLength = 10

	for i := range count {
		length := rand.Intn(maxLength) + 1

		key := make([]byte, length)
		for j := range length {
			key[j] = charset[rand.Intn(len(charset))]
		}

		keys[i] = string(key)
	}

	return keys
}

// generateKeysWithCommonPrefix 共通プレフィックスを持つキーを生成
func generateKeysWithCommonPrefix(prefix string, suffixLength, count int) []string {
	keys := make([]string, count)

	const charset = "abcdefghijklmnopqrstuvwxyz"

	for i := range count {
		suffix := make([]byte, suffixLength)
		for j := range suffixLength {
			suffix[j] = charset[rand.Intn(len(charset))]
		}

		keys[i] = prefix + string(suffix)
	}

	return keys
}
