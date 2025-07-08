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
    
    # 各辞書ファイルの単語リスト
    local small_words="${japanese_dir}/small_words.txt"
    local core_words="${japanese_dir}/core_words.txt"
    local notcore_words="${japanese_dir}/notcore_words.txt"
    local full_words="${japanese_dir}/full_words.txt"
    
    # 全ての辞書ファイルをダウンロード・処理
    if [ ! -f "${full_words}" ]; then
        mkdir -p "${japanese_dir}"
        
        # ベースURL
        local base_url="https://d2ej7fkh96fzlu.cloudfront.net/sudachidict-raw/20250515"
        
        # 各辞書ファイルの処理
        echo "  🔽 Sudachi辞書データをダウンロード中..."
        
        # small_lex.zip（基本語彙、約40MB）
        if [ ! -f "${small_words}" ]; then
            echo "    📥 small_lex.zip をダウンロード中..."
            curl -L -o "${japanese_dir}/small_lex.zip" "${base_url}/small_lex.zip"
            unzip -o -q "${japanese_dir}/small_lex.zip" -d "${japanese_dir}"
            cut -d',' -f1 "${japanese_dir}/small_lex.csv" | grep -v '^$' | sort -u > "${small_words}"
            echo "    ✅ small_lex: $(wc -l < "${small_words}")語"
        fi
        
        # core_lex.zip（追加語彙、約21MB）
        if [ ! -f "${core_words}" ]; then
            echo "    📥 core_lex.zip をダウンロード中..."
            curl -L -o "${japanese_dir}/core_lex.zip" "${base_url}/core_lex.zip"
            unzip -o -q "${japanese_dir}/core_lex.zip" -d "${japanese_dir}"
            cut -d',' -f1 "${japanese_dir}/core_lex.csv" | grep -v '^$' | sort -u > "${core_words}"
            echo "    ✅ core_lex: $(wc -l < "${core_words}")語"
        fi
        
        # notcore_lex.zip（専門語彙、約35MB）
        if [ ! -f "${notcore_words}" ]; then
            echo "    📥 notcore_lex.zip をダウンロード中..."
            curl -L -o "${japanese_dir}/notcore_lex.zip" "${base_url}/notcore_lex.zip"
            unzip -o -q "${japanese_dir}/notcore_lex.zip" -d "${japanese_dir}"
            cut -d',' -f1 "${japanese_dir}/notcore_lex.csv" | grep -v '^$' | sort -u > "${notcore_words}"
            echo "    ✅ notcore_lex: $(wc -l < "${notcore_words}")語"
        fi
        
        # 全辞書統合（重複削除）
        echo "  🔄 全辞書データを統合中..."
        sort -u "${small_words}" "${core_words}" "${notcore_words}" > "${full_words}"
        echo "  ✅ 統合辞書: $(wc -l < "${full_words}")語"
        
        # テスト用ファイル作成（既存の処理を維持）
        if [ ! -f "${test_file}" ]; then
            # 小規模テスト用（1000語）- 漢字で始まる単語のみ
            grep -P '^[\x{4e00}-\x{9fa5}]' "${small_words}" | head -1000 > "${test_file}"
        fi
        
        # ベンチマーク用ファイル（small_lexのみ、既存の処理を維持）
        if [ ! -f "${bench_file}" ]; then
            cp "${small_words}" "${bench_file}"
        fi
        
        # 大規模ベンチマーク用データファイル作成
        local large_bench_file="${japanese_dir}/large_bench.csv"
        local mega_bench_file="${japanese_dir}/mega_bench.csv"
        
        # large_bench.csv: core辞書まで含む（約250万語）
        if [ ! -f "${large_bench_file}" ]; then
            echo "  🔄 大規模ベンチマーク用データを作成中..."
            sort -u "${small_words}" "${core_words}" > "${large_bench_file}"
            echo "    ✅ 大規模ベンチマーク用: $(wc -l < "${large_bench_file}")語"
        fi
        
        # mega_bench.csv: 全辞書統合（約800万語）
        if [ ! -f "${mega_bench_file}" ]; then
            echo "  🔄 超大規模ベンチマーク用データを作成中..."
            cp "${full_words}" "${mega_bench_file}"
            echo "    ✅ 超大規模ベンチマーク用: $(wc -l < "${mega_bench_file}")語"
        fi
        
        # full_bench.csv: 全辞書統合（mega_benchのエイリアス）
        local full_bench_file="${japanese_dir}/full_bench.csv"
        if [ ! -f "${full_bench_file}" ]; then
            echo "  🔄 Full辞書ベンチマーク用データを作成中..."
            cp "${full_words}" "${full_bench_file}"
            echo "    ✅ Full辞書ベンチマーク用: $(wc -l < "${full_bench_file}")語"
        fi
        
        echo "  ✅ 日本語辞書データを準備完了"
        echo "    - テスト用: $(wc -l < "${test_file}")語"
        echo "    - ベンチマーク用（small）: $(wc -l < "${bench_file}")語"
        echo "    - 大規模ベンチマーク用: $(wc -l < "${large_bench_file}")語"
        echo "    - 超大規模ベンチマーク用: $(wc -l < "${mega_bench_file}")語"
        echo "    - Full辞書ベンチマーク用: $(wc -l < "${full_bench_file}")語"
        echo "    - フル辞書語彙: $(wc -l < "${full_words}")語"
    else
        echo "  ✅ 日本語辞書データは既に存在します"
        echo "    - テスト用: $(wc -l < "${test_file}")語"
        echo "    - ベンチマーク用（small）: $(wc -l < "${bench_file}")語"
        local large_bench_file="${japanese_dir}/large_bench.csv"
        local mega_bench_file="${japanese_dir}/mega_bench.csv"
        local full_bench_file="${japanese_dir}/full_bench.csv"
        if [ -f "${large_bench_file}" ]; then
            echo "    - 大規模ベンチマーク用: $(wc -l < "${large_bench_file}")語"
        fi
        if [ -f "${mega_bench_file}" ]; then
            echo "    - 超大規模ベンチマーク用: $(wc -l < "${mega_bench_file}")語"
        fi
        if [ -f "${full_bench_file}" ]; then
            echo "    - Full辞書ベンチマーク用: $(wc -l < "${full_bench_file}")語"
        fi
        echo "    - フル辞書語彙: $(wc -l < "${full_words}")語"
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