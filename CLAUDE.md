# CLAUDE.md

このリポジトリでコードを扱う際のClaude Code (claude.ai/code) への指針。

**重要**: このリポジトリでは日本語でコミュニケーション。コメント、ドキュメント、およびClaude Codeとのやり取りは日本語かつ体言止めで記述。markdownはmarkdownlintに従った構文で記述。

## プロジェクト概要

パトリシアトライの実装プロジェクト。パトリシアトライは、単一の子を持つノードを親ノードと結合することで、メモリ使用量を削減し、検索性能を向上させた空間効率化されたトライ木データ構造。

## リポジトリ状況

現在は空のリポジトリで、既存のコード、ビルド設定、ドキュメントは存在しない状態。

## 開発セットアップ

### 技術スタック

- **言語**: Go
- **ビルドシステム**: make
- **テストフレームワーク**: stretchr/testify
- **静的解析**: golangci-lint
- **プロジェクト構造**: Go標準（internalは使用せず、pkgの下にコードを配置）

### プロジェクト構造

```text
patricia-trie/
├── Makefile
├── go.mod
├── go.sum
├── README.md
├── CLAUDE.md
├── cmd/
│   └── example/
│       └── main.go
├── pkg/
│   └── patriciatrie/
│       ├── trie.go
│       ├── trie_test.go
│       ├── node.go
│       └── node_test.go
└── docs/
    └── examples/
```

### 開発コマンド

```bash
# ビルド
make build

# テスト実行
make test

# テストカバレッジ
make test-coverage

# ベンチマーク
make benchmark

# 静的解析
make lint

# 整形
make fmt
```

## アーキテクチャの考慮事項

パトリシアトライ実装時の重要ポイント：

- **コアデータ構造**: パス圧縮をサポートするノード構造の設計
- **主要操作**: 挿入、検索、削除、プレフィックスマッチング操作の実装
- **メモリ最適化**: 適切なパス圧縮による効率的なメモリ使用
- **パフォーマンス**: 時間計算量と空間計算量の両方の最適化

## 次のステップ

1. Go module初期化とMakefile作成
2. 基本的なプロジェクト構造のセットアップ
3. コアパトリシアトライデータ構造の実装
4. stretchr/testifyを使用したテストスイートの追加
5. ベンチマークとパフォーマンステストの実装
6. 使用例とドキュメントの作成

## Git設定

リポジトリは初期化済みだが、まだコミットなし。

### ブランチ戦略

Trunk-based開発を採用。

- **mainブランチ**: 常にデプロイ可能な状態を維持
- **短命な機能ブランチ**: 数日以内で完了する小さな機能開発用
- **devブランチ**: 大きな機能開発時の統合用（必要に応じて）

### 開発ワークフロー

1. mainブランチから短命な機能ブランチを作成
2. 小さな単位で機能開発・テスト実装
3. 頻繁にmainブランチへPR作成・マージ
4. 大きな機能の場合は、devブランチで統合後にmainへPR

### コミットメッセージ

体言止めで簡潔に記述。

```text
feat: パトリシアトライの基本構造を実装
fix: 検索処理のバグを修正
test: 挿入操作のテストケースを追加
docs: README.mdの使用例を更新
```

## 英語・日本語対訳表

このプロジェクトで使用する技術用語の統一表記。新しい用語は随時追加。

| 英語 | 日本語 |
|------|--------|
| Patricia Trie | パトリシアトライ |
| Radix Tree | 基数木 |
| Trie | トライ木 |
| Node | ノード |
| Edge | エッジ |
| Path compression | パス圧縮 |
| Prefix | プレフィックス |
| Insert | 挿入 |
| Search | 検索 |
| Delete | 削除 |
| Lookup | 検索 |
| Key | キー |
| Value | 値 |
| Root | ルート |
| Leaf | 葉 |
| Branch | 分岐 |
| Compression | 圧縮 |
| Algorithm | アルゴリズム |
| Data structure | データ構造 |
| Time complexity | 時間計算量 |
| Space complexity | 空間計算量 |
| Memory optimization | メモリ最適化 |
| Performance | パフォーマンス |
| Implementation | 実装 |
| Test suite | テストスイート |
| Build system | ビルドシステム |
| Package management | パッケージ管理 |
| Go module | Goモジュール |
| Makefile | Makefile |
| Test coverage | テストカバレッジ |
| Benchmark | ベンチマーク |
| Static analysis | 静的解析 |
| Formatting | 整形 |
| Linting | リント |
| golangci-lint | golangci-lint |
| Example | 例 |
| Documentation | ドキュメント |
| markdownlint | markdownlint |
| Pull Request | PR |
| Branch | ブランチ |
| Merge | マージ |
| Review | レビュー |
| Commit message | コミットメッセージ |
| Feature branch | 機能ブランチ |
| Main branch | mainブランチ |
| Development branch | 開発ブランチ |
| Trunk-based development | Trunk-based開発 |
| Short-lived branch | 短命ブランチ |