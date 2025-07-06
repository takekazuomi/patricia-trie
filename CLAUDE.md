# CLAUDE.md

このリポジトリでコードを扱う際のClaude Code (claude.ai/code) への指針。

**重要**: このリポジトリでは日本語でコミュニケーション。コメント、ドキュメント、およびClaude Codeとのやり取りは日本語かつ体言止めで記述。markdownはmarkdownlintに従った構文で記述。

## プロジェクト概要

パトリシアトライの実装プロジェクト。パトリシアトライは、単一の子を持つノードを親ノードと結合することで、メモリ使用量を削減し、検索性能を向上させた空間効率化されたトライ木データ構造。

## リポジトリ状況

- 基本的なPatricia Trie実装完了
- ユニットテスト、ベンチマークテスト実装済み
- 大規模ベンチマーク（〜100万キー）と実世界データベンチマーク追加

## 開発セットアップ

### 技術スタック

- **言語**: Go 1.24.2
- **ビルドシステム**: make
- **テストフレームワーク**: stretchr/testify
- **静的解析**: golangci-lint v2.2.1
- **プロジェクト構造**: Go標準（internalは使用せず、pkgの下にコードを配置）

### プロジェクト構造

```text
patricia-trie/
├── Makefile
├── go.mod
├── go.sum
├── README.md
├── CLAUDE.md
├── CONTRIBUTING.md
├── FAQ.md
├── .github/
│   └── workflows/
│       └── ci.yml
├── .claude/
│   └── settings.json
├── cmd/
│   └── example/
│       └── main.go
├── pkg/
│   └── patriciatrie/
│       ├── trie.go
│       ├── trie_test.go
│       ├── node.go
│       ├── node_test.go
│       ├── bench_test.go
│       ├── large_bench_test.go
│       └── realistic_bench_test.go
├── scripts/
│   ├── install-deps.sh
│   └── setup_benchmark_data.sh
├── docs/
│   ├── examples/
│   ├── large-benchmark.md
│   └── sudachi-lex-csv.md
└── testdata/
    ├── .gitignore
    ├── README.md
    ├── japanese/
    └── ipaddresses/
```

### 開発コマンド

```bash
# セットアップ
make setup            # 開発環境セットアップ
make install-deps     # 依存ツールインストール

# ビルド
make build           # バイナリをビルド

# テスト実行
make test            # テストを実行
make test-coverage   # テストカバレッジ取得

# ベンチマーク
make benchmark            # 基本ベンチマーク実行
make setup_benchmark      # ベンチマークデータセットアップ
make benchmark-large      # 大規模データベンチマーク
make benchmark-realistic  # 実世界データベンチマーク

# 静的解析・整形
make lint            # golangci-lint実行
make fmt             # コード整形
make check           # fmt, lint, test を一括実行

# その他
make clean           # 生成ファイル削除
make help            # ヘルプ表示
```

## アーキテクチャの考慮事項

パトリシアトライ実装時の重要ポイント：

- **コアデータ構造**: パス圧縮をサポートするノード構造の設計
- **主要操作**: 挿入、検索、削除、プレフィックスマッチング操作の実装
- **メモリ最適化**: 適切なパス圧縮による効率的なメモリ使用
- **パフォーマンス**: 時間計算量と空間計算量の両方の最適化

## Claude Codeへの制約事項

### 変更禁止ファイル

- `.golangci.yml`: linter設定の変更禁止（`.claude/settings.json`で定義）

### コミットメッセージ

- Claude署名（🤖 Generated with...）を付けない
- 体言止めで簡潔に記述

## 今後の課題

1. GitHub Actionsのローカル実行環境構築（act使用）
2. 並行処理対応の実装
3. より高度な最適化（メモリプール、カスタムアロケータ）
4. プロファイリングとボトルネック分析
5. APIドキュメントの充実

## Git設定

- **デフォルトブランチ**: main
- **CI/CD**: GitHub Actions（Go 1.22, 1.24でのマトリックステスト）

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
| Large-scale benchmark | 大規模ベンチマーク |
| Realistic benchmark | 実世界データベンチマーク |
| Memory efficiency | メモリ効率 |
| Throughput | スループット |
| Latency | レイテンシ |
| Profiling | プロファイリング |
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
| CI/CD | CI/CD |
| GitHub Actions | GitHub Actions |
| act | act |
| Local CI | ローカルCI |