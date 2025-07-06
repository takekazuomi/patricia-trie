# 安全なMakefileのためのプリアンブル
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif

.PHONY: build test test-coverage benchmark benchmark-large benchmark-realistic setup_benchmark lint fmt clean clean-all clean-testdata install-deps setup mod-tidy check ci-local ci-full help

# フルCIワークフロー
ci-full: ## テスト、静的解析、ビルドを実行
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) build
# デフォルトターゲット
.DEFAULT_GOAL := help

# 変数定義
BINARY_NAME := patricia-trie
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html

# ビルド
build: cmd/example/main.go ## バイナリをビルド
	go build -o bin/$(BINARY_NAME) $<

# テスト
test: ## テストを実行
	go test -v ./...

# テストカバレッジ
test-coverage: ## テストカバレッジを取得
	go test -v -coverprofile=$(COVERAGE_OUT) ./...
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "カバレッジレポート: $(COVERAGE_HTML)"

# ベンチマーク
benchmark: ## 基本ベンチマークを実行
	go test -bench=. -benchmem ./...

# 大規模ベンチマーク
benchmark-large: setup_benchmark ## 大規模データセットでベンチマークを実行
	go test -bench=BenchmarkTrie_Large -benchmem -timeout=30m ./pkg/patriciatrie

# リアルデータベンチマーク
benchmark-realistic: setup_benchmark ## 日本語・IPアドレスデータでベンチマークを実行
	go test -bench=BenchmarkTrie_Japanese -benchmem -timeout=30m ./pkg/patriciatrie
	go test -bench=BenchmarkTrie_IPv -benchmem -timeout=30m ./pkg/patriciatrie

# ベンチマーク用データセットアップ
setup_benchmark: ## ベンチマーク用テストデータをセットアップ
	@./scripts/setup_benchmark_data.sh

# 静的解析
# リンターが見つからない場合は make install-deps を実行してください
lint: ## golangci-lintとmarkdownlintを実行
	golangci-lint run
	markdownlint-cli2 README.md CLAUDE.md FAQ.md CONTRIBUTING.md testdata/README.md docs/*.md

# フォーマット
fmt: ## コードをフォーマット
	go fmt ./...
	goimports -w .

# 依存関係の整理
mod-tidy: ## go mod tidyを実行
	go mod tidy

# クリーンアップ
clean: ## 生成されたファイルを削除
	rm -rf bin/
	rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)

# 完全クリーンアップ
clean-all: clean ## 依存ツールも含めて完全削除
	rm -rf tmp/

# テストデータクリーンアップ
clean-testdata: ## ベンチマーク用テストデータを削除
	rm -rf testdata/*.txt testdata/*/*.txt testdata/*.csv testdata/*/*.csv

# 依存ツールのインストール
install-deps: ## 開発に必要なツールをインストール
	@./scripts/install-deps.sh

# 開発用のセットアップ
setup: install-deps ## 開発環境をセットアップ
	go mod download

# 全体チェック
check: fmt lint test ## フォーマット、リント、テストを実行

# GitHub Actions CI をローカルで再現
ci-local: ## GitHub Actions CI パイプラインをローカルで実行
	@echo "🚀 GitHub Actions CI パイプラインをローカルで実行中..."
	@echo ""
	@echo "📋 Step 1: 依存関係をダウンロード"
	@go mod download
	@echo ""
	@echo "🧪 Step 2: テストを実行"
	@$(MAKE) test
	@echo ""
	@echo "📊 Step 3: テストカバレッジを生成"
	@$(MAKE) test-coverage
	@echo ""
	@echo "🔍 Step 4: 静的解析を実行"
	@$(MAKE) lint
	@echo ""
	@echo "🔨 Step 5: ビルドを実行"
	@$(MAKE) build
	@echo ""
	@echo "⚡ Step 6: ベンチマークを実行（基本）"
	@$(MAKE) benchmark
	@echo ""
	@echo "✅ GitHub Actions CI パイプライン完了！"
	@echo "   すべてのステップが正常に完了しました。"
	@echo "   GitHub Actionsでも同様に成功するはずです。"

# ヘルプ
help: ## このヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'