.PHONY: build test test-coverage benchmark lint fmt clean help

# デフォルトターゲット
.DEFAULT_GOAL := help

# 変数定義
BINARY_NAME := patricia-trie
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html

# ビルド
build: ## バイナリをビルド
	go build -o bin/$(BINARY_NAME) ./cmd/example

# テスト
test: ## テストを実行
	go test -v ./...

# テストカバレッジ
test-coverage: ## テストカバレッジを取得
	go test -v -coverprofile=$(COVERAGE_OUT) ./...
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "カバレッジレポート: $(COVERAGE_HTML)"

# ベンチマーク
benchmark: ## ベンチマークを実行
	go test -bench=. -benchmem ./...

# 静的解析
lint: ## golangci-lintを実行
	golangci-lint run

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

# 開発用のセットアップ
setup: ## 開発環境をセットアップ
	go mod download
	go install golang.org/x/tools/cmd/goimports@latest

# 全体チェック
check: fmt lint test ## フォーマット、リント、テストを実行

# ヘルプ
help: ## このヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'