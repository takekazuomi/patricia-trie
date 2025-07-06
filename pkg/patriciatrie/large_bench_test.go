package patriciatrie

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
)

// BenchmarkTrie_Large_Insert 大規模データでの挿入性能測定
func BenchmarkTrie_Large_Insert(b *testing.B) {
	dataSizes := []int{10000, 100000, 1000000}

	for _, size := range dataSizes {
		b.Run(fmt.Sprintf("Keys_%d", size), func(b *testing.B) {
			keys := generateLargeRandomKeys(size, 20)

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

// BenchmarkTrie_Large_Search 大規模データでの検索性能測定
func BenchmarkTrie_Large_Search(b *testing.B) {
	dataSizes := []int{10000, 100000, 1000000}

	for _, size := range dataSizes {
		b.Run(fmt.Sprintf("Keys_%d", size), func(b *testing.B) {
			keys := generateLargeRandomKeys(size, 20)
			trie := New()

			// 事前データ挿入
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

// BenchmarkTrie_Large_Memory 大規模データでのメモリ使用量測定
func BenchmarkTrie_Large_Memory(b *testing.B) {
	dataSizes := []int{10000, 100000, 500000}

	for _, size := range dataSizes {
		b.Run(fmt.Sprintf("Keys_%d", size), func(b *testing.B) {
			keys := generateLargeRandomKeys(size, 15)

			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				// メモリ測定開始
				var m1, m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m1)

				trie := New()
				for _, key := range keys {
					_ = trie.Insert(key)
				}

				runtime.GC()
				runtime.ReadMemStats(&m2)

				// メモリ使用量を記録（非公式だが参考値として）
				memUsed := m2.Alloc - m1.Alloc
				b.ReportMetric(float64(memUsed)/float64(size), "bytes/key")
			}
		})
	}
}

// BenchmarkTrie_Large_Mixed 大規模データでの混合操作性能測定
func BenchmarkTrie_Large_Mixed(b *testing.B) {
	dataSizes := []int{50000, 100000}

	for _, size := range dataSizes {
		b.Run(fmt.Sprintf("Keys_%d", size), func(b *testing.B) {
			keys := generateLargeRandomKeys(size, 15)
			trie := New()

			// 初期データの半分を挿入
			for i := range len(keys) / 2 {
				_ = trie.Insert(keys[i])
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := range b.N {
				key := keys[i%len(keys)]

				switch i % 4 {
				case 0: // 25% 挿入
					_ = trie.Insert(key)
				case 1, 2: // 50% 検索
					_ = trie.Search(key)
				case 3: // 25% プレフィックス検索
					prefix := key[:minValue(len(key), 3)]
					_ = trie.FindByPrefix(prefix)
				}
			}
		})
	}
}

// BenchmarkTrie_Large_PrefixSearch 大規模データでのプレフィックス検索性能
func BenchmarkTrie_Large_PrefixSearch(b *testing.B) {
	dataSize := 100000
	prefixLengths := []int{2, 4, 6, 8}

	for _, prefixLen := range prefixLengths {
		b.Run(fmt.Sprintf("PrefixLen_%d", prefixLen), func(b *testing.B) {
			keys := generateLargeRandomKeys(dataSize, 15)
			trie := New()

			// データ挿入
			for _, key := range keys {
				_ = trie.Insert(key)
			}

			// プレフィックス候補を準備
			prefixes := make([]string, 1000)
			for i := range prefixes {
				key := keys[rand.Intn(len(keys))]
				if len(key) >= prefixLen {
					prefixes[i] = key[:prefixLen]
				} else {
					prefixes[i] = key
				}
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := range b.N {
				prefix := prefixes[i%len(prefixes)]
				results := trie.FindByPrefix(prefix)
				_ = results
			}
		})
	}
}

// BenchmarkTrie_Large_WorstCase 最悪ケースシナリオの性能測定
func BenchmarkTrie_Large_WorstCase(b *testing.B) {
	scenarios := []struct {
		name string
		keys []string
	}{
		{
			name: "Sequential",
			keys: generateSequentialKeys(10000),
		},
		{
			name: "CommonPrefix",
			keys: generateKeysWithVeryLongCommonPrefix("verylongcommonprefix", 50, 10000),
		},
		{
			name: "SingleChar",
			keys: generateSingleCharKeys(10000),
		},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				trie := New()
				for _, key := range scenario.keys {
					_ = trie.Insert(key)
				}

				// 検索も含める
				for i := range minValue(1000, len(scenario.keys)) {
					key := scenario.keys[i]
					_ = trie.Search(key)
				}
			}
		})
	}
}

// generateLargeRandomKeys 大規模ランダムキー生成
func generateLargeRandomKeys(count, maxLength int) []string {
	// Go 1.20以降では自動的にシードされるためSeedは不要
	keys := make([]string, count)

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for i := range count {
		length := rand.Intn(maxLength-1) + 2 // 2〜maxLength

		key := make([]byte, length)
		for j := range length {
			key[j] = charset[rand.Intn(len(charset))]
		}

		keys[i] = string(key)
	}

	return keys
}

// generateSequentialKeys 連番キー生成（最悪ケース用）
func generateSequentialKeys(count int) []string {
	keys := make([]string, count)
	for i := range count {
		keys[i] = fmt.Sprintf("key%08d", i)
	}

	return keys
}

// generateKeysWithVeryLongCommonPrefix 非常に長い共通プレフィックスを持つキー生成
func generateKeysWithVeryLongCommonPrefix(prefix string, suffixLength, count int) []string {
	// Go 1.20以降では自動的にシードされるためSeedは不要
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

// generateSingleCharKeys 単一文字キー生成（最悪ケース用）
func generateSingleCharKeys(count int) []string {
	keys := make([]string, count)
	for i := range count {
		keys[i] = string(rune('a' + (i % 26)))
	}

	return keys
}

// minValue ヘルパー関数
func minValue(a, b int) int {
	if a < b {
		return a
	}

	return b
}
