#!/usr/bin/env bash

set -euo pipefail

# スクリプトのディレクトリを取得
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "${SCRIPT_DIR}")"
TESTDATA_DIR="${PROJECT_ROOT}/testdata"

echo "=== ベンチマーク用テストデータのセットアップ ==="

# 日本語辞書データのセットアップ
setup_japanese_data() {
    echo "📚 日本語辞書データをセットアップ中..."
    
    local japanese_dir="${TESTDATA_DIR}/japanese"
    local test_file="${japanese_dir}/1000.csv"
    local bench_file="${japanese_dir}/all.csv"
    
    # 小規模テスト用（1000語）と全データ用に分割
    if [ ! -f "${test_file}" ] || [ ! -f "${bench_file}" ]; then
        echo "  🔽 Sudachi辞書データ（small_lex.zip）をダウンロード中..."
        mkdir -p "${japanese_dir}"
        
        # Sudachi raw lexicon data をダウンロード
        local temp_file="${japanese_dir}/small_lex.zip"
        curl -L -o "${temp_file}" \
            "https://d2ej7fkh96fzlu.cloudfront.net/sudachidict-raw/20250515/small_lex.zip"
        
        # 解凍（上書き確認なし）
        unzip -o -q "${temp_file}" -d "${japanese_dir}"
        
        # CSVファイルから単語を抽出（1列目が表層形、2列目が左文脈ID、3列目が右文脈ID...）
        echo "  🔄 辞書データを処理中..."
        
        # small_lex.csv の1列目（表層形）を抽出してテスト用・ベンチマーク用作成
        if [ ! -f "${japanese_dir}/small_lex.csv" ]; then
            echo "  ❌ エラー: small_lex.csv が見つかりません"
            echo "     ダウンロードまたは解凍に失敗した可能性があります"
            exit 1
        fi
        
        # CSVの1列目（表層形）を抽出、重複削除、空行除去
        cut -d',' -f1 "${japanese_dir}/small_lex.csv" | \
            grep -v '^$' | \
            sort -u > "${japanese_dir}/temp_words.txt"
        
        # 小規模テスト用（1000語）- 漢字で始まる単語のみ
        grep -P '^[\x{4e00}-\x{9fa5}]' "${japanese_dir}/temp_words.txt" | \
            head -1000 > "${test_file}"
        
        # 全データ用（ベンチマーク用）
        cp "${japanese_dir}/temp_words.txt" "${bench_file}"
        
        # 一時ファイル削除
        #rm -f "${temp_file}" "${japanese_dir}/small_lex.csv" "${japanese_dir}/temp_words.txt"
        
        echo "  ✅ 日本語辞書データを準備完了"
        echo "    - テスト用: $(wc -l < "${test_file}")語"
        echo "    - ベンチマーク用: $(wc -l < "${bench_file}")語"
    else
        echo "  ✅ 日本語辞書データは既に存在します"
        echo "    - テスト用: $(wc -l < "${test_file}")語"
        echo "    - ベンチマーク用: $(wc -l < "${bench_file}")語"
    fi
}

# IPアドレスデータの生成
setup_ipaddress_data() {
    echo "🌐 IPアドレスデータを生成中..."
    
    local ip_dir="${TESTDATA_DIR}/ipaddresses"
    mkdir -p "${ip_dir}"
    
    # IPv4アドレス生成
    local ipv4_10k="${ip_dir}/ipv4_10k.txt"
    local ipv4_100k="${ip_dir}/ipv4_100k.txt"
    
    if [ ! -f "${ipv4_10k}" ] || [ ! -f "${ipv4_100k}" ]; then
        echo "  🔄 IPv4アドレスを生成中..."
        
        # Go script for generating IP addresses
        cat > "${ip_dir}/generate_ips.go" << 'EOF'
package main

import (
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "time"
)

func generateIPv4(count int) []string {
    rand.Seed(time.Now().UnixNano())
    ips := make([]string, count)
    
    for i := 0; i < count; i++ {
        ips[i] = fmt.Sprintf("%d.%d.%d.%d",
            rand.Intn(224)+1,   // 1-224 (避けるべき範囲を除外)
            rand.Intn(256),     // 0-255
            rand.Intn(256),     // 0-255
            rand.Intn(254)+1)   // 1-254
    }
    return ips
}

func generateIPv6(count int) []string {
    rand.Seed(time.Now().UnixNano())
    ips := make([]string, count)
    
    for i := 0; i < count; i++ {
        ips[i] = fmt.Sprintf("2001:db8:%04x:%04x:%04x:%04x:%04x:%04x",
            rand.Intn(65536), rand.Intn(65536), rand.Intn(65536), rand.Intn(65536),
            rand.Intn(65536), rand.Intn(65536))
    }
    return ips
}

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Usage: go run generate_ips.go <type> <count> <output>")
        os.Exit(1)
    }
    
    ipType := os.Args[1]
    count, _ := strconv.Atoi(os.Args[2])
    output := os.Args[3]
    
    var ips []string
    if ipType == "ipv4" {
        ips = generateIPv4(count)
    } else if ipType == "ipv6" {
        ips = generateIPv6(count)
    }
    
    file, err := os.Create(output)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    for _, ip := range ips {
        fmt.Fprintln(file, ip)
    }
    
    fmt.Printf("Generated %d %s addresses\n", count, ipType)
}
EOF
        
        # IPv4アドレス生成
        cd "${ip_dir}"
        go run generate_ips.go ipv4 10000 "ipv4_10k.txt"
        go run generate_ips.go ipv4 100000 "ipv4_100k.txt"
        
        # IPv6アドレス生成
        go run generate_ips.go ipv6 10000 "ipv6_10k.txt"
        go run generate_ips.go ipv6 100000 "ipv6_100k.txt"
        
        # 一時ファイル削除
        rm -f generate_ips.go
        cd "${PROJECT_ROOT}"
        
        echo "  ✅ IPアドレスデータを生成完了"
    else
        echo "  ✅ IPアドレスデータは既に存在します"
    fi
}

# データサイズ報告
report_data_sizes() {
    echo ""
    echo "📊 生成されたテストデータサイズ:"
    
    if [ -d "${TESTDATA_DIR}" ]; then
        find "${TESTDATA_DIR}" -name "*.txt" -o -name "*.csv" | sort | while read -r file; do
            echo "  $(basename "${file}"): $(wc -l < "${file}") lines ($(du -h "${file}" | cut -f1))"
        done
    fi
}

# メイン処理
main() {
    setup_japanese_data
    setup_ipaddress_data
    report_data_sizes
    
    echo ""
    echo "🎉 ベンチマーク用テストデータのセットアップが完了しました！"
    echo ""
    echo "次のステップ:"
    echo "  1. 大規模ベンチマーク実行: make benchmark-large"
    echo "  2. リアルデータベンチマーク: make benchmark-realistic"
}

main "$@"