#!/usr/bin/env bash

set -euo pipefail

# スクリプトのディレクトリを取得
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "${SCRIPT_DIR}")"

echo "=== 開発依存ツールのインストール ==="

# golangci-lintのインストール
echo "📦 golangci-lintをインストール中..."
mkdir -p "${PROJECT_ROOT}/tmp"
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${PROJECT_ROOT}/tmp" latest
echo "✅ golangci-lint v2がインストールされました"

# goimportsのインストール
echo "📦 goimportsをインストール中..."
go install golang.org/x/tools/cmd/goimports@latest
echo "✅ goimportsがインストールされました"

# direnvのチェックと設定
echo "🔍 direnvの設定をチェック中..."
if command -v direnv >/dev/null 2>&1; then
    echo "✅ direnvが利用可能です"
    cd "${PROJECT_ROOT}"
    direnv allow
    echo "✅ .envrcが許可されました"
else
    echo "⚠️  警告: direnvがインストールされていません"
    echo "   手動でPATHを設定するか、direnvをインストールしてください:"
    echo "   export PATH=\"${PROJECT_ROOT}/tmp:\${PATH}\""
fi

echo ""
echo "🎉 依存ツールのインストールが完了しました！"
echo ""
echo "次のステップ:"
echo "  1. direnvを使用する場合: source .envrc または新しいシェルを開く"
echo "  2. 手動でPATHを設定する場合: export PATH=\"${PROJECT_ROOT}/tmp:\${PATH}\""
echo "  3. リントを実行: make lint"