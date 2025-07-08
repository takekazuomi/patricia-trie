# CLAUDE.md

このリポジトリでコードを扱う際のClaude Code (claude.ai/code) への指針。

**重要**: このリポジトリでは日本語でコミュニケーション。コメント、ドキュメント、およびClaude Codeとのやり取りは日本語かつ体言止めで記述。markdownはmarkdownlintに従った構文で記述。

## プロジェクト概要

パトリシアトライの実装プロジェクト。パトリシアトライは、単一の子を持つノードを親ノードと結合することで、メモリ使用量を削減し、検索性能を向上させた空間効率化されたトライ木データ構造。

## リポジトリ状況

- 基本的なPatricia Trie実装完了
- ユニットテスト、ベンチマークテスト実装済み
- 大規模ベンチマーク（〜100万キー）と実世界データベンチマーク追加
- README.mdを更新し、Claude Codeテストプロジェクトであることを明記
- REPL（対話的検索）ツール実装完了

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
│   ├── example/
│   │   └── main.go
│   └── patricia-repl/
│       ├── main.go
│       └── README.md
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
    │   ├── small_lex.zip         # Sudachi辞書（ダウンロード元）
    │   ├── core_lex.zip          # Sudachi辞書（ダウンロード元）
    │   ├── notcore_lex.zip       # Sudachi辞書（ダウンロード元）
    │   ├── 1000.txt              # テスト用漢字語彙（1K語）
    │   ├── small.txt             # Small辞書（約57万語）
    │   ├── core.txt              # Core辞書（約82万語）
    │   ├── notcore.txt           # NotCore辞書（約124万語）
    │   └── full.txt              # 全辞書統合（約259万語）
    └── ipaddresses/
        ├── ipv4_10k.txt          # IPv4アドレス（10K個）
        ├── ipv4_100k.txt         # IPv4アドレス（100K個）
        ├── ipv6_10k.txt          # IPv6アドレス（10K個）
        └── ipv6_100k.txt         # IPv6アドレス（100K個）
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
make build-example   # exampleバイナリをビルド
make build-repl      # REPLバイナリをビルド

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
- **REPLツール**: go-promptを使用した対話的検索インターフェース

## Claude Codeへの制約事項

### 変更禁止ファイル

- `.golangci.yml`: linter設定の変更禁止（`.claude/settings.json`で定義）

### コミットメッセージ

- Claude署名（🤖 Generated with...）を付けない
- 体言止めで簡潔に記述

## 今後の課題

1. GitHub Actionsのローカル実行環境構築（act使用）
2. 並行処理対応の実装
3. REPLツールの機能拡張
   - 実際のノード探索統計の実装（patriciatrieパッケージの拡張）
   - 履歴の逆検索機能（Ctrl+R）
   - 大量結果時のページング機能
4. より高度な最適化（メモリプール、カスタムアロケータ）
5. プロファイリングとボトルネック分析
6. APIドキュメントの充実
7. Goバージョンの更新（Go 1.25リリース後）

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

### コミットメッセージスタイル

体言止めで簡潔に記述。

```text
feat: パトリシアトライの基本構造を実装
fix: 検索処理のバグを修正
test: 挿入操作のテストケースを追加
docs: README.mdの使用例を更新
chore: 依存関係の更新
refactor: コードのリファクタリング
```

#### プレフィックスの使い分け

- **feat**: 新機能の追加
- **fix**: バグ修正
- **docs**: ドキュメントのみの変更
- **chore**: ビルドプロセスや補助ツールの変更
- **test**: テストの追加・修正
- **refactor**: コードのリファクタリング

## 英語・日本語対訳表

このプロジェクトで使用する技術用語の統一表記。新しい用語は随時追加。

| 英語 | 日本語 |
|------|--------|
| act | act |
| Algorithm | アルゴリズム |
| Allocator | アロケータ |
| API | API |
| Assertion | アサーション |
| Benchmark | ベンチマーク |
| Bottleneck | ボトルネック |
| Branch | ブランチ |
| Build system | ビルドシステム |
| CI/CD | CI/CD |
| CLI | CLI |
| Commit message | コミットメッセージ |
| Compression | 圧縮 |
| Configuration | 設定 |
| Conflict | コンフリクト |
| Context | コンテキスト |
| Continuous delivery | 継続的デリバリー |
| Continuous integration | 継続的統合 |
| Corpus | コーパス |
| CSV | CSV |
| Data structure | データ構造 |
| Database | データベース |
| Delete | 削除 |
| Delimiter | 区切り文字 |
| Dependency | 依存関係 |
| Deployment | デプロイ |
| Development branch | 開発ブランチ |
| Dictionary | 辞書 |
| Documentation | ドキュメント |
| Edge | エッジ |
| Encoding | エンコーディング |
| Environment | 環境 |
| Example | 例 |
| Feature branch | 機能ブランチ |
| Feature flag | 機能フラグ |
| Formatting | 整形 |
| Framework | フレームワーク |
| Garbage collection | ガベージコレクション |
| GitHub Actions | GitHub Actions |
| Go module | Go module |
| golangci-lint | golangci-lint |
| Hook | フック |
| Hot fix | ホットフィックス |
| Implementation | 実装 |
| Insert | 挿入 |
| Issue | Issue |
| Key | キー |
| Large-scale benchmark | 大規模ベンチマーク |
| Latency | レイテンシ |
| Leaf | リーフ |
| Library | ライブラリ |
| License | ライセンス |
| Linting | リント |
| Local CI | ローカルCI |
| Logging | ログ |
| Lookup | 検索 |
| Main branch | mainブランチ |
| Makefile | Makefile |
| markdownlint | markdownlint |
| Matrix test | マトリクステスト |
| Memory efficiency | メモリ効率 |
| Memory optimization | メモリ最適化 |
| Memory pool | メモリプール |
| Merge | マージ |
| Metrics | メトリクス |
| Migration | 移行 |
| Mock | モック |
| Monitoring | モニタリング |
| Node | ノード |
| Normalization | 正規化 |
| Package management | パッケージ管理 |
| Parsing | 解析 |
| Path compression | パス圧縮 |
| Patricia Trie | パトリシアトライ |
| Performance | パフォーマンス |
| Pipeline | パイプライン |
| Prefix | プレフィックス |
| Profiling | プロファイリング |
| Pull Request | PR |
| Queue | キュー |
| Radix Tree | 基数木 |
| Realistic benchmark | 実世界データベンチマーク |
| Refactoring | リファクタリング |
| Regression | 回帰 |
| Repository | リポジトリ |
| Resilience | 回復力 |
| Resource | リソース |
| Review | レビュー |
| Rollback | ロールバック |
| Root | ルート |
| Scalability | スケーラビリティ |
| Search | 検索 |
| Short-lived branch | 短命ブランチ |
| Space complexity | 空間計算量 |
| Static analysis | 静的解析 |
| Streaming | ストリーミング |
| Synchronization | 同期 |
| Template | テンプレート |
| Test coverage | テストカバレッジ |
| Test suite | テストスイート |
| Throughput | スループット |
| Time complexity | 時間計算量 |
| Tokenizer | トークナイザー |
| Trie | トライ木 |
| Trunk-based development | Trunk-based開発 |
| Validation | 検証 |
| Value | 値 |
| Vulnerability | 脆弱性 |
| Workflow | ワークフロー |
| Workload | ワークロード |
