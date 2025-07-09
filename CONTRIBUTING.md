# コントリビューションガイド

このプロジェクトへの貢献に関するガイドライン。

## Issue・PR管理方針

このプロジェクトでは**最小限PR・包括的Issue方式**を採用している。

### 📋 Issue（技術情報の集約地）

**役割**: 完全な技術ドキュメント・知識ベース

Issueには以下の詳細情報を記載：

- **背景・要件**: 技術的課題、実装理由、目標設定
- **技術設計**: アーキテクチャ方針、データ設計、実装方針  
- **検証・結果**: ベンチマーク結果、パフォーマンス分析、技術的成果
- **使用方法**: セットアップ手順、実行例、設定方法
- **影響範囲**: 変更統計、技術的価値、今後の展開

**メリット**: 技術情報が永続的に保存され、後の参照・学習に活用可能

### 🔄 PR（変更要約のみ）

**役割**: 簡潔な変更チェックリスト

PRには以下の最小限の情報のみ記載：

- **変更概要**: 何を変更したか（1-2行）
- **変更ファイル**: 主要な変更ファイルと変更理由
- **品質保証**: テスト結果、lint結果、レビュー対応状況
- **詳細参照**: `詳細は Issue #XX を参照`で誘導

**メリット**: 開発効率が向上し、レビューがコード品質に集中可能

### 📝 記述例

**Issue例**:
```markdown
## 背景・課題
- 現在の実装では大規模データに対応していない
- 実世界データでの性能評価が必要

## 技術設計
### アーキテクチャ方針
- データ統一: CSV+TXT混在 → 統一TXT形式

## 実装内容
### 1. セットアップスクリプト強化
- ファイル: scripts/setup_benchmark_data.sh
- 内容: 3辞書の個別ダウンロード・処理

## 性能測定結果
| データセット | 語数 | 検索時間 |
|-------------|------|----------|
| Small | 57万語 | 262ns |
```

**PR例**:
```markdown
# 大規模単語リスト対応（Issue #15）

Sudachi辞書の全データ（最大約259万語）を使用した段階的ベンチマークを実装

## 変更ファイル
- scripts/setup_benchmark_data.sh: 3辞書の統一処理
- pkg/patriciatrie/realistic_bench_test.go: 大規模専用ベンチマーク追加

## 品質保証
- ✅ 全テスト成功
- ✅ golangci-lint・markdownlint成功

技術詳細は Issue #15 を参照

Closes #15
```

## 開発スタイル

- **言語**: 日本語（体言止め）
- **コードスタイル**: Go標準
- **テスト**: stretchr/testifyを使用
- **静的解析**: golangci-lintを使用
- **ドキュメント**: markdownlintに従った構文

## Trunk-based開発

このプロジェクトはTrunk-based開発を採用。

### 基本概念

Trunk-based開発は、開発者が1つのブランチ（trunk/main）に対して頻繁に小さなコミットを行う開発手法。

### 特徴

- **単一の真実の源**: mainブランチが常に最新の状態
- **短命なブランチ**: 機能ブランチは数日以内で完了
- **頻繁な統合**: 1日に複数回のマージが発生
- **継続的テスト**: 全てのコミットでテストが実行される

### 従来のGit Flowとの違い

| 項目 | Git Flow | Trunk-based |
|------|----------|-------------|
| ブランチ寿命 | 数週間〜数ヶ月 | 数時間〜数日 |
| 統合頻度 | 週1回程度 | 1日数回 |
| コンフリクト | 発生しやすい | 発生しにくい |
| リリース | 専用ブランチ | mainから直接 |

### メリット

- **コンフリクト削減**: 頻繁な統合によりマージコンフリクトを最小化
- **品質向上**: 小さな変更により問題の早期発見
- **デプロイ効率**: 常にデプロイ可能な状態を維持
- **開発速度**: 長期間のブランチ管理が不要

## 開発ワークフロー

### 1. 機能開発の準備

```bash
# 最新のmainブランチを取得
git checkout main
git pull origin main

# 短命な機能ブランチを作成
git checkout -b feature/node-implementation
```

### 2. 開発サイクル

```bash
# 小さな変更を実装
# テストを追加・更新
make test

# 変更をコミット
git add .
git commit -m "feat: ノード構造体の基本実装"

# 必要に応じて中間プッシュ
git push origin feature/node-implementation
```

### 3. マージの準備

```bash
# mainブランチの最新変更を取得
git checkout main
git pull origin main

# 機能ブランチにマージ
git checkout feature/node-implementation
git merge main

# コンフリクトがある場合は解決
# テストを実行
make test
make lint  # golangci-lintを実行
```

### 4. プルリクエスト

- **小さな単位**: 1つのPRで1つの機能に集中
- **頻繁な提出**: 完了したらすぐにPRを作成
- **迅速なレビュー**: 24時間以内のレビューを目標

## ブランチ命名規則

| 種類 | 命名規則 | 例 |
|------|----------|-----|
| 機能追加 | feature/xxx | feature/search-operation |
| バグ修正 | fix/xxx | fix/memory-leak |
| ドキュメント | docs/xxx | docs/api-reference |
| テスト | test/xxx | test/benchmark |
| リファクタリング | refactor/xxx | refactor/node-structure |

## コミットメッセージ

体言止めで簡潔に記述。

### 形式

```text
type: 変更内容の要約

詳細な説明（必要に応じて）
```

### 種類

- **feat**: 新機能の追加
- **fix**: バグの修正
- **docs**: ドキュメントの更新
- **test**: テストの追加・修正
- **refactor**: コードのリファクタリング
- **perf**: パフォーマンスの改善
- **style**: コードスタイルの修正

### 例

```text
feat: パトリシアトライの挿入操作を実装

- ノードの分割処理を追加
- エッジラベルの管理機能を実装
- 重複キーの処理を改善
```

## 大きな機能の開発

### devブランチの使用

複数の小さな機能を統合する場合：

```bash
# devブランチを作成
git checkout -b dev/patricia-trie-v1

# 小さな機能を順次マージ
git merge feature/node-structure
git merge feature/insert-operation
git merge feature/search-operation

# 統合テストを実行
make test
make benchmark

# mainブランチへPR
```

### 機能フラグ

未完成の機能をmainに統合する場合：

```go
// 機能フラグを使用
if config.EnableNewFeature {
    // 新機能の実装
}
```

## 品質管理

### 必須チェック

- [ ] 全てのテストが通過
- [ ] リント検査が通過
- [ ] ベンチマークが劣化していない
- [ ] ドキュメントが更新されている

### 自動化

```bash
# プリコミットフック
make test lint

# CI/CDパイプライン
make test test-coverage benchmark lint
```

## トラブルシューティング

### コンフリクトの解決

```bash
# mainブランチの最新を取得
git checkout main
git pull origin main

# 機能ブランチにマージ
git checkout feature/xxx
git merge main

# コンフリクトを解決
git add .
git commit -m "fix: コンフリクトを解決"
```

### 長期ブランチの回避

- 3日以上続くブランチは分割を検討
- 機能を小さな単位に分解
- 必要に応じて機能フラグを活用

## Makefileの設計方針

### 安全なMakefileの原則

このプロジェクトでは[Your Makefiles are wrong](https://tech.davis-hansson.com/p/make/)の方針に基づく安全なMakefileを採用：

- **厳格なシェル設定**: `bash -eu -o pipefail`でエラー時の即座停止
- **未定義変数の警告**: 未定義変数使用時に警告表示
- **エラー時クリーンアップ**: `.DELETE_ON_ERROR`で不完全なファイルを自動削除
- **明示的な依存関係**: ビルトインルールを無効化し明示的に定義
- **PHONYターゲット**: ファイル生成しないターゲットを明示

### シンプルさを保つ

Makefileは複雑にしすぎない。以下の原則に従う：

- **単純なコマンドの組み合わせ**: 基本的なタスクの定義に留める
- **複雑なロジックは外部スクリプト**: 条件分岐や複雑な処理は`scripts/`ディレクトリのシェルスクリプトに移譲
- **1ターゲット1責任**: 各ターゲットは単一の明確な責任を持つ
- **依存関係は最小限**: 過度な依存関係は避ける
- **タブインデント**: 伝統的なタブを使用（`.RECIPEPREFIX`は使用しない）
- **センチネルファイル不使用**: シンプルさを保つためセンチネルファイルは使用しない
- **Magic Variables活用**: `$@`, `$<`, `$^`などの自動変数を積極的に使用

### 良い例

```makefile
# シンプルなコマンド実行
test: ## テストを実行
    go test -v ./...

# 外部スクリプトの呼び出し
install-deps: ## 開発依存ツールをインストール
    @./scripts/install-deps.sh
```

### 避けるべき例

```makefile
# 複雑すぎる例（避ける）
install-deps: ## 開発依存ツールをインストール
    @echo "golangci-lintをインストール中..."
    @mkdir -p tmp
    @if [ ! -f ./tmp/golangci-lint ]; then \
        curl -sSfL https://example.com/install.sh | sh -s -- -b ./tmp latest; \
    fi
    @if command -v direnv >/dev/null 2>&1; then \
        direnv allow; \
    else \
        echo "警告: direnvが必要です"; \
    fi
```

### Magic Variablesの活用

Makefileの可読性と保守性を向上させるため、自動変数を積極的に使用：

```makefile
# 良い例: Magic Variablesを使用
build: cmd/example/main.go
    go build -o bin/$(BINARY_NAME) $<

# 避ける例: ハードコーディング
build: cmd/example/main.go
    go build -o bin/patricia-trie cmd/example/main.go
```

#### 主要なMagic Variables

- `$@`: ターゲット名
- `$<`: 最初の依存関係
- `$^`: すべての依存関係（重複除去）
- `$+`: すべての依存関係（重複含む）
- `$?`: ターゲットより新しい依存関係

### 複雑な処理の外部化

- **スクリプト配置**: `scripts/`ディレクトリに配置
- **実行権限**: `chmod +x`で実行可能にする
- **エラーハンドリング**: `set -euo pipefail`を使用
- **プロジェクトルート**: `$(dirname "${BASH_SOURCE[0]}")`で相対パス解決

## 参考文献

- [Trunk-based Development](https://trunkbaseddevelopment.com/)
- [Google's Engineering Practices](https://google.github.io/eng-practices/)
- [Effective Go](https://golang.org/doc/effective_go.html)
