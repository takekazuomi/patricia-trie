package main

import (
	"fmt"
	"log"

	"github.com/takekazu/patricia-trie/pkg/patriciatrie"
)

func main() {
	// パトリシアトライの使用例
	trie := patriciatrie.New()

	// キーを挿入
	keys := []string{"cat", "cats", "dog", "dogs", "elephant"}
	
	fmt.Println("=== パトリシアトライの使用例 ===")
	
	for _, key := range keys {
		if err := trie.Insert(key); err != nil {
			log.Fatalf("キーの挿入に失敗: %v", err)
		}
		fmt.Printf("挿入: %s\n", key)
	}
	
	fmt.Println("\n=== 検索テスト ===")
	
	// 検索テスト
	searchKeys := []string{"cat", "ca", "dog", "elephant", "notfound"}
	
	for _, key := range searchKeys {
		found := trie.Search(key)
		fmt.Printf("検索: %s -> %v\n", key, found)
	}
	
	fmt.Println("\n=== プレフィックス検索 ===")
	
	// プレフィックス検索
	prefixes := []string{"ca", "dog", "el"}
	
	for _, prefix := range prefixes {
		matches := trie.FindByPrefix(prefix)
		fmt.Printf("プレフィックス: %s -> %v\n", prefix, matches)
	}
}