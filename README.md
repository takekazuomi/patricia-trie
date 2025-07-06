# パトリシアトライ（Patricia Trie）

> **Note**
>
> このプロジェクトは[Claude Code](https://claude.ai/code)のテスト用リポジトリ。
> 内容をプロンプトで指示して、コードのほとんどはClaude Codeで書いている。
>
> テストも生成AIが書いており現時点では、正確性の確認は限定的。
> コードの流用はお勧めしない。

パトリシアトライまたはRadix Treeとも呼ばれる、空間効率化されたトライ木データ構造のGo実装。

## 概要

パトリシアトライは通常のトライ木を改良したデータ構造。単一の子を持つノードを親ノードと結合することで、メモリ使用量を大幅に削減し、検索性能を向上。

### 特徴

- **空間効率**: 通常のトライ木と比較して大幅なメモリ使用量削減
- **高速検索**: O(k)の時間計算量（kはキーの長さ）
- **プレフィックス検索**: 指定されたプレフィックスを持つすべてのキーを効率的に検索
- **動的挿入・削除**: 実行時にキーの追加・削除が可能

## 動作原理

### 通常のトライ木との違い

通常のトライ木では、文字列の各文字に対して個別のノードを作成。

```text
        root
       /    \
      c      s
      |      |
      a      u
      |      |
      t      n
     /       |
    s        \
             sun
```

パトリシアトライでは、単一の子を持つノードを圧縮：

```text
        root
       /    \
     "cat"  "sun"
     /
   "s"
```

### パス圧縮の仕組み

1. **単一パス検出**: 一直線に続く単一の子ノードを特定
2. **パス結合**: 複数のノードを1つのノードに結合
3. **エッジラベル**: 結合されたパスを1つのエッジラベルとして格納

## 主要操作

### 挿入（Insert）

新しいキーを木に挿入。既存のノードとの分岐点を見つけて適切な位置に配置。

### 検索（Search）

指定されたキーの存在を確認。エッジラベルを辿って目的のキーに到達できるかを判定。

### 削除（Delete）

キーを削除し、必要に応じてノードを再結合してパス圧縮を維持。

### プレフィックス検索（Prefix Search）

指定されたプレフィックスを持つすべてのキーを効率的に検索。

## アルゴリズム計算量

| 操作 | 時間計算量 | 空間計算量 |
|------|------------|------------|
| 挿入 | O(k) | O(1) |
| 検索 | O(k) | O(1) |
| 削除 | O(k) | O(1) |
| プレフィックス検索 | O(k + m) | O(m) |

- k: キーの長さ
- m: マッチするキーの数

## インストール

```bash
go get github.com/takekazu/patricia-trie/pkg/patriciatrie
```

## クイックスタート

```bash
# リポジトリのクローン
git clone https://github.com/takekazu/patricia-trie.git
cd patricia-trie

# 開発環境のセットアップ
make setup

# テストの実行
make test

# ベンチマークの実行
make benchmark
```

## 実装済み機能

- ✅ 基本的なパトリシアトライ構造
- ✅ 挿入操作（Insert）
- ✅ 検索操作（Search）
- ✅ 削除操作（Delete）
- ✅ プレフィックス検索（FindByPrefix）
- ✅ 全キー取得（AllKeys）
- ✅ サイズ取得（Size）
- ✅ 空判定（IsEmpty）

## 使用例

```go
package main

import (
    "fmt"
    "github.com/takekazu/patricia-trie/pkg/patriciatrie"
)

func main() {
    // トライの作成
    trie := patriciatrie.New()
    
    // キーと値の挿入
    trie.Insert("cat", "猫")
    trie.Insert("cats", "猫たち")
    trie.Insert("dog", "犬")
    trie.Insert("dogs", "犬たち")
    
    // 検索
    if value, found := trie.Search("cat"); found {
        fmt.Printf("Found: %s\n", value) // Found: 猫
    }
    
    // プレフィックス検索
    keys := trie.FindByPrefix("cat")
    fmt.Println(keys) // [cat cats]
    
    // 削除
    trie.Delete("cat")
    
    // サイズ取得
    fmt.Printf("Size: %d\n", trie.Size()) // Size: 3
}
```

## ベンチマーク

### 基本ベンチマーク

```bash
make benchmark
```

### 大規模データベンチマーク（100万キー）

```bash
make setup_benchmark      # ベンチマークデータのセットアップ
make benchmark-large      # 大規模ベンチマークの実行
```

### 実世界データベンチマーク（日本語辞書・IPアドレス）

```bash
make benchmark-realistic  # 実世界データベンチマークの実行
```

## 開発

### 必要な環境

- Go 1.24.2以上
- make
- golangci-lint v2.2.1（静的解析用）

### 開発コマンド

```bash
make help          # 利用可能なコマンド一覧
make test          # テスト実行
make lint          # 静的解析
make fmt           # コード整形
make check         # fmt, lint, testを一括実行
```

### プロジェクト構造

```text
patricia-trie/
├── pkg/patriciatrie/     # パトリシアトライ実装
│   ├── trie.go          # メインのトライ構造
│   ├── node.go          # ノード構造
│   └── *_test.go        # テストファイル
├── cmd/example/         # 使用例
├── testdata/            # テスト用データ
└── docs/                # ドキュメント
```

## 貢献

[CONTRIBUTING.md](CONTRIBUTING.md)を参照。

## FAQ

[FAQ.md](FAQ.md)を参照。

## 参考文献

- [Radix tree - Wikipedia](https://en.wikipedia.org/wiki/Radix_tree)
- [パトリシアトライ - 論文](https://dl.acm.org/doi/10.1145/321479.321481)
- [Data Structures and Algorithms - Robert Sedgewick](https://algs4.cs.princeton.edu/)
