package patriciatrie

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// BenchmarkTrie_Japanese_Insert 日本語辞書データでの挿入性能
func BenchmarkTrie_Japanese_Insert(b *testing.B) {
	datasets := []struct {
		name string
		file string
	}{
		{"Test_1K", "testdata/japanese/1000.csv"},
		{"Bench_All", "testdata/japanese/all.csv"},
	}

	for _, dataset := range datasets {
		b.Run(dataset.name, func(b *testing.B) {
			words, err := loadWordsFromFile(dataset.file)
			if err != nil {
				b.Skipf("テストデータが見つかりません: %s (make setup_benchmarkを実行してください)", dataset.file)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				trie := New()
				for _, word := range words {
					_ = trie.Insert(word)
				}
			}

			b.ReportMetric(float64(len(words)), "words")
		})
	}
}

// BenchmarkTrie_Japanese_Search 日本語辞書データでの検索性能
func BenchmarkTrie_Japanese_Search(b *testing.B) {
	datasets := []struct {
		name string
		file string
	}{
		{"Test_1K", "testdata/japanese/1000.csv"},
		{"Bench_All", "testdata/japanese/all.csv"},
	}

	for _, dataset := range datasets {
		b.Run(dataset.name, func(b *testing.B) {
			words, err := loadWordsFromFile(dataset.file)
			if err != nil {
				b.Skipf("テストデータが見つかりません: %s", dataset.file)
			}

			trie := New()
			for _, word := range words {
				_ = trie.Insert(word)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := range b.N {
				word := words[i%len(words)]
				_ = trie.Search(word)
			}
		})
	}
}

// BenchmarkTrie_Japanese_PrefixSearch 日本語データでのプレフィックス検索
func BenchmarkTrie_Japanese_PrefixSearch(b *testing.B) {
	words, err := loadWordsFromFile("testdata/japanese/all.csv")
	if err != nil {
		b.Skipf("テストデータが見つかりません (make setup_benchmarkを実行してください)")
	}

	trie := New()
	for _, word := range words {
		_ = trie.Insert(word)
	}

	// 様々な長さのプレフィックスを準備
	prefixLengths := []int{1, 2, 3, 4}

	for _, prefixLen := range prefixLengths {
		b.Run(fmt.Sprintf("PrefixLen_%d", prefixLen), func(b *testing.B) {
			prefixes := generateJapanesePrefixes(words, prefixLen, 1000)

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

// BenchmarkTrie_IPv4_Insert IPv4アドレスでの挿入性能
func BenchmarkTrie_IPv4_Insert(b *testing.B) {
	datasets := []struct {
		name string
		file string
	}{
		{"10K", "testdata/ipaddresses/ipv4_10k.txt"},
		{"100K", "testdata/ipaddresses/ipv4_100k.txt"},
	}

	for _, dataset := range datasets {
		b.Run(dataset.name, func(b *testing.B) {
			ips, err := loadWordsFromFile(dataset.file)
			if err != nil {
				b.Skipf("テストデータが見つかりません: %s (make setup_benchmarkを実行してください)", dataset.file)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				trie := New()
				for _, ip := range ips {
					_ = trie.Insert(ip)
				}
			}

			b.ReportMetric(float64(len(ips)), "addresses")
		})
	}
}

// BenchmarkTrie_IPv4_Search IPv4アドレスでの検索性能
func BenchmarkTrie_IPv4_Search(b *testing.B) {
	datasets := []struct {
		name string
		file string
	}{
		{"10K", "testdata/ipaddresses/ipv4_10k.txt"},
		{"100K", "testdata/ipaddresses/ipv4_100k.txt"},
	}

	for _, dataset := range datasets {
		b.Run(dataset.name, func(b *testing.B) {
			ips, err := loadWordsFromFile(dataset.file)
			if err != nil {
				b.Skipf("テストデータが見つかりません: %s", dataset.file)
			}

			trie := New()
			for _, ip := range ips {
				_ = trie.Insert(ip)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := range b.N {
				ip := ips[i%len(ips)]
				_ = trie.Search(ip)
			}
		})
	}
}

// BenchmarkTrie_IPv4_NetworkPrefix IPv4ネットワークプレフィックス検索
func BenchmarkTrie_IPv4_NetworkPrefix(b *testing.B) {
	ips, err := loadWordsFromFile("testdata/ipaddresses/ipv4_100k.txt")
	if err != nil {
		b.Skipf("テストデータが見つかりません (make setup_benchmarkを実行してください)")
	}

	trie := New()
	for _, ip := range ips {
		_ = trie.Insert(ip)
	}

	// ネットワークプレフィックス（/8, /16, /24相当）を準備
	prefixTypes := []struct {
		name   string
		length int
	}{
		{"Class_A", 3},  // "192" (xxx.*)
		{"Class_B", 7},  // "192.168" (xxx.xxx.*)
		{"Class_C", 11}, // "192.168.1" (xxx.xxx.xxx.*)
	}

	for _, prefixType := range prefixTypes {
		b.Run(prefixType.name, func(b *testing.B) {
			prefixes := generateIPv4Prefixes(ips, prefixType.length, 100)

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

// BenchmarkTrie_IPv6_Insert IPv6アドレスでの挿入性能
func BenchmarkTrie_IPv6_Insert(b *testing.B) {
	datasets := []struct {
		name string
		file string
	}{
		{"10K", "testdata/ipaddresses/ipv6_10k.txt"},
		{"100K", "testdata/ipaddresses/ipv6_100k.txt"},
	}

	for _, dataset := range datasets {
		b.Run(dataset.name, func(b *testing.B) {
			ips, err := loadWordsFromFile(dataset.file)
			if err != nil {
				b.Skipf("テストデータが見つかりません: %s (make setup_benchmarkを実行してください)", dataset.file)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				trie := New()
				for _, ip := range ips {
					_ = trie.Insert(ip)
				}
			}

			b.ReportMetric(float64(len(ips)), "addresses")
		})
	}
}

// BenchmarkTrie_IPv6_Search IPv6アドレスでの検索性能
func BenchmarkTrie_IPv6_Search(b *testing.B) {
	datasets := []struct {
		name string
		file string
	}{
		{"10K", "testdata/ipaddresses/ipv6_10k.txt"},
		{"100K", "testdata/ipaddresses/ipv6_100k.txt"},
	}

	for _, dataset := range datasets {
		b.Run(dataset.name, func(b *testing.B) {
			ips, err := loadWordsFromFile(dataset.file)
			if err != nil {
				b.Skipf("テストデータが見つかりません: %s", dataset.file)
			}

			trie := New()
			for _, ip := range ips {
				_ = trie.Insert(ip)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := range b.N {
				ip := ips[i%len(ips)]
				_ = trie.Search(ip)
			}
		})
	}
}

// BenchmarkTrie_Japanese_Memory 日本語データでのメモリ使用量測定
func BenchmarkTrie_Japanese_Memory(b *testing.B) {
	datasets := []struct {
		name string
		file string
	}{
		{"Test_1K", "testdata/japanese/1000.csv"},
		{"Bench_All", "testdata/japanese/all.csv"},
	}

	for _, dataset := range datasets {
		b.Run(dataset.name, func(b *testing.B) {
			words, err := loadWordsFromFile(dataset.file)
			if err != nil {
				b.Skipf("テストデータが見つかりません: %s (make setup_benchmarkを実行してください)", dataset.file)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				// メモリ測定開始
				var m1, m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m1)

				trie := New()
				for _, word := range words {
					_ = trie.Insert(word)
				}

				runtime.GC()
				runtime.ReadMemStats(&m2)

				// メモリ使用量を記録
				memUsed := m2.Alloc - m1.Alloc
				b.ReportMetric(float64(memUsed)/float64(len(words)), "bytes/word")
				b.ReportMetric(float64(len(words)), "words")
			}
		})
	}
}

// BenchmarkTrie_IPv4_Memory IPv4アドレスでのメモリ使用量測定
func BenchmarkTrie_IPv4_Memory(b *testing.B) {
	datasets := []struct {
		name string
		file string
	}{
		{"10K", "testdata/ipaddresses/ipv4_10k.txt"},
		{"100K", "testdata/ipaddresses/ipv4_100k.txt"},
	}

	for _, dataset := range datasets {
		b.Run(dataset.name, func(b *testing.B) {
			ips, err := loadWordsFromFile(dataset.file)
			if err != nil {
				b.Skipf("テストデータが見つかりません: %s (make setup_benchmarkを実行してください)", dataset.file)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				// メモリ測定開始
				var m1, m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m1)

				trie := New()
				for _, ip := range ips {
					_ = trie.Insert(ip)
				}

				runtime.GC()
				runtime.ReadMemStats(&m2)

				// メモリ使用量を記録
				memUsed := m2.Alloc - m1.Alloc
				b.ReportMetric(float64(memUsed)/float64(len(ips)), "bytes/address")
				b.ReportMetric(float64(len(ips)), "addresses")
			}
		})
	}
}

// BenchmarkTrie_IPv6_Memory IPv6アドレスでのメモリ使用量測定
func BenchmarkTrie_IPv6_Memory(b *testing.B) {
	datasets := []struct {
		name string
		file string
	}{
		{"10K", "testdata/ipaddresses/ipv6_10k.txt"},
		{"100K", "testdata/ipaddresses/ipv6_100k.txt"},
	}

	for _, dataset := range datasets {
		b.Run(dataset.name, func(b *testing.B) {
			ips, err := loadWordsFromFile(dataset.file)
			if err != nil {
				b.Skipf("テストデータが見つかりません: %s (make setup_benchmarkを実行してください)", dataset.file)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				// メモリ測定開始
				var m1, m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m1)

				trie := New()
				for _, ip := range ips {
					_ = trie.Insert(ip)
				}

				runtime.GC()
				runtime.ReadMemStats(&m2)

				// メモリ使用量を記録
				memUsed := m2.Alloc - m1.Alloc
				b.ReportMetric(float64(memUsed)/float64(len(ips)), "bytes/address")
				b.ReportMetric(float64(len(ips)), "addresses")
			}
		})
	}
}

// BenchmarkTrie_Mixed_Workload 混合ワークロード（日本語+IPv4）
func BenchmarkTrie_Mixed_Workload(b *testing.B) {
	japanese, err1 := loadWordsFromFile("testdata/japanese/1000.csv")

	ipv4, err2 := loadWordsFromFile("testdata/ipaddresses/ipv4_10k.txt")
	if err1 != nil || err2 != nil {
		b.Skipf("テストデータが見つかりません (make setup_benchmarkを実行してください)")
	}
	// データを混合
	mixedData := make([]string, 0, len(japanese)+len(ipv4))
	mixedData = append(mixedData, japanese...)
	mixedData = append(mixedData, ipv4...)

	trie := New()
	for _, data := range mixedData {
		_ = trie.Insert(data)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := range b.N {
		data := mixedData[i%len(mixedData)]

		switch i % 3 {
		case 0: // 33% 検索
			_ = trie.Search(data)
		case 1: // 33% プレフィックス検索（短い）
			if len(data) >= 2 {
				_ = trie.FindByPrefix(data[:2])
			}
		case 2: // 33% プレフィックス検索（長い）
			if len(data) >= 4 {
				_ = trie.FindByPrefix(data[:4])
			}
		}
	}
}

// loadWordsFromFile ファイルから単語リストを読み込み
func loadWordsFromFile(filename string) ([]string, error) {
	// プロジェクトルートからの相対パスを解決
	if !filepath.IsAbs(filename) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		// pkg/patriciatrieから2つ上がプロジェクトルート
		projectRoot := filepath.Join(wd, "..", "..")
		filename = filepath.Join(projectRoot, filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() { _ = file.Close() }()

	var words []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			words = append(words, line)
		}
	}

	return words, scanner.Err()
}

// generateJapanesePrefixes 日本語単語からプレフィックスを生成
func generateJapanesePrefixes(words []string, prefixLen, count int) []string {
	prefixes := make([]string, 0, count)

	for _, word := range words {
		if len([]rune(word)) >= prefixLen {
			runes := []rune(word)
			prefix := string(runes[:prefixLen])
			prefixes = append(prefixes, prefix)

			if len(prefixes) >= count {
				break
			}
		}
	}

	// 不足分は先頭から循環（無限ループ回避のため一度だけ実行）
	if len(prefixes) < count {
		for _, word := range words {
			if len([]rune(word)) >= prefixLen {
				runes := []rune(word)
				prefix := string(runes[:prefixLen])
				prefixes = append(prefixes, prefix)

				if len(prefixes) >= count {
					break
				}
			}
		}
	}

	return prefixes
}

// generateIPv4Prefixes IPv4アドレスからネットワークプレフィックスを生成
func generateIPv4Prefixes(ips []string, prefixLen, count int) []string {
	prefixes := make([]string, 0, count)

	for _, ip := range ips {
		if len(ip) >= prefixLen {
			prefix := ip[:prefixLen]
			prefixes = append(prefixes, prefix)

			if len(prefixes) >= count {
				break
			}
		}
	}

	// 不足分は先頭から循環（無限ループ回避のため一度だけ実行）
	if len(prefixes) < count {
		for _, ip := range ips {
			if len(ip) >= prefixLen {
				prefix := ip[:prefixLen]
				prefixes = append(prefixes, prefix)

				if len(prefixes) >= count {
					break
				}
			}
		}
	}

	return prefixes
}
