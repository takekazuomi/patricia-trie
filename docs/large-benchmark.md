# 大規模ベンチマーク仕様書

## 概要

本ドキュメントでは、Patricia Trieの大規模データセットにおける性能特性を測定するためのベンチマーク仕様について説明します。

## 目的

1. **スケーラビリティの評価**: 10K〜1Mキーでの性能特性の把握
2. **実用性の検証**: 日本語・IPアドレスなど実データでの性能測定  
3. **メモリ効率の分析**: データサイズ増加に対するメモリ使用量の評価
4. **最悪ケースの特定**: 性能劣化が発生する条件の特定

## データソース

### 1. 日本語コーパス（Sudachi）

**選定理由:**

- マルチバイト文字によるトライ構造への影響を測定
- 自然言語の語彙分布における実用性評価
- Apache-2.0ライセンスで商用利用可能

**データ特性:**

- コア辞書: ~100K語
- 拡張辞書: ~1M語  
- 文字種: ひらがな、カタカナ、漢字
- 語長分布: 1〜20文字程度
- 共通プレフィックス: 自然言語特有のパターン

**取得方法:**

```bash
curl -L -o small_lex.zip \
  "https://d2ej7fkh96fzlu.cloudfront.net/sudachidict-raw/20250515/small_lex.zip"
```

### 2. IPアドレス

**選定理由:**

- 構造化データでのトライ効率性の評価
- ネットワークプレフィックス検索の実用性測定
- 固定フォーマットでの最適化効果確認

**データ特性:**

*IPv4:*

- 形式: xxx.xxx.xxx.xxx (7〜15文字)
- 文字種: 数字とドット
- 階層構造: ネットワーク/ホスト分離
- プレフィックス: /8, /16, /24 相当

*IPv6:*

- 形式: xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx
- 文字種: 16進数とコロン
- 長さ: 39文字固定
- プレフィックス: /64, /48 相当

## ベンチマーク項目

### 1. 大規模データベンチマーク (`large_bench_test.go`)

#### BenchmarkTrie_Large_Insert

- **目的**: 大量データ挿入時のスループット測定
- **データサイズ**: 10K, 100K, 1M キー
- **測定値**: ops/sec, メモリアロケーション

#### BenchmarkTrie_Large_Search  

- **目的**: 大量データ検索時のレイテンシ測定
- **前提条件**: 事前にすべてのキーを挿入
- **測定値**: search ops/sec

#### BenchmarkTrie_Large_Memory

- **目的**: メモリ効率性の定量評価
- **測定値**: bytes/key, 総メモリ使用量
- **手法**: runtime.MemStatsによる測定

#### BenchmarkTrie_Large_Mixed

- **目的**: 実運用に近い混合ワークロード測定
- **操作比率**: Insert 25%, Search 50%, PrefixSearch 25%
- **測定値**: 混合操作のスループット

#### BenchmarkTrie_Large_PrefixSearch

- **目的**: プレフィックス長による性能影響評価
- **プレフィックス長**: 2, 4, 6, 8文字
- **測定値**: prefix search ops/sec, 結果件数

#### BenchmarkTrie_Large_WorstCase

- **目的**: 最悪ケースシナリオでの性能確認
- **シナリオ**:
  - Sequential: 連番キー（key00000001〜）
  - CommonPrefix: 超長共通プレフィックス
  - SingleChar: 単一文字キー（a, b, c...）

### 2. リアルデータベンチマーク (`realistic_bench_test.go`)

#### 日本語辞書ベンチマーク

```go
BenchmarkTrie_Japanese_Insert        // 辞書構築性能
BenchmarkTrie_Japanese_Search        // 単語検索性能  
BenchmarkTrie_Japanese_PrefixSearch  // 前方一致検索
BenchmarkTrie_Japanese_Memory        // 日本語データでのメモリ効率
```

#### IPアドレスベンチマーク

```go
BenchmarkTrie_IPv4_Insert           // IPv4登録性能
BenchmarkTrie_IPv4_Search           // IPv4検索性能
BenchmarkTrie_IPv4_NetworkPrefix    // ネットワーク検索
BenchmarkTrie_IPv4_Memory           // IPv4データでのメモリ効率
BenchmarkTrie_IPv6_Insert           // IPv6登録性能
BenchmarkTrie_IPv6_Search           // IPv6検索性能
BenchmarkTrie_IPv6_Memory           // IPv6データでのメモリ効率
```

#### 混合ワークロード

```go
BenchmarkTrie_Mixed_Workload         // 日本語+IPv4混合
```

## セットアップ手順

### 1. テストデータ準備

```bash
# ベンチマーク用データを自動セットアップ
make setup_benchmark
```

**実行内容:**

1. Sudachi辞書データのダウンロード・変換
2. IPv4/IPv6アドレスのランダム生成
3. データファイルの配置

**生成ファイル:**

```text
testdata/
├── japanese/
│   ├── 1000.csv              # 小規模テスト用 1K語（漢字開始のみ）
│   └── all.csv               # ベンチマーク用 全語彙（約57万語）
└── ipaddresses/
    ├── ipv4_10k.txt          # IPv4 10K個
    ├── ipv4_100k.txt         # IPv4 100K個  
    ├── ipv6_10k.txt          # IPv6 10K個
    └── ipv6_100k.txt         # IPv6 100K個
```

### 2. ベンチマーク実行

```bash
# 基本ベンチマーク
make benchmark

# 大規模データベンチマーク  
make benchmark-large

# リアルデータベンチマーク
make benchmark-realistic

# テストデータクリーンアップ
make clean-testdata
```

## 実装詳細

### データ生成アルゴリズム

#### 日本語辞書処理

```bash
# CSVファイルから表層形（1列目）を抽出
cut -d',' -f1 small_lex.csv | grep -v '^$' | sort -u

# 小規模テスト用: 漢字で始まる単語のみを抽出
grep -P '^[\x{4e00}-\x{9fa5}]' temp_words.txt | head -1000
```

#### IPアドレス生成  

```go
// IPv4: 1.0.0.1 〜 224.255.255.254 の範囲でランダム生成
func generateIPv4(count int) []string {
    // プライベート/予約アドレスを除外した実用的な範囲
}

// IPv6: 2001:db8::/32 プレフィックスでテスト用アドレス生成  
func generateIPv6(count int) []string {
    // RFC 3849 ドキュメント用プレフィックス使用
}
```

### パフォーマンス測定

#### メモリ効率測定

```go
var m1, m2 runtime.MemStats
runtime.GC()
runtime.ReadMemStats(&m1)

// トライ構築

runtime.GC()  
runtime.ReadMemStats(&m2)
memUsed := m2.Alloc - m1.Alloc
```

#### カスタムメトリクス

```go
b.ReportMetric(float64(memUsed)/float64(keyCount), "bytes/key")
b.ReportMetric(float64(len(results)), "matches")
```

## 期待される成果

### 性能特性の定量化

- **スループット**: 各操作の ops/sec
- **レイテンシ**: p50, p95, p99レスポンス時間  
- **メモリ効率**: bytes/key, 圧縮率
- **スケーラビリティ**: データサイズ増加に対する性能変化

### 実用性の評価

- **日本語処理**: マルチバイト文字でのトライ効率
- **ネットワーク用途**: IPアドレス管理での実用性
- **混合ワークロード**: 実運用シナリオでの性能

### 最適化指針  

- **ボトルネック特定**: CPU/メモリ使用パターン分析
- **改善案策定**: アルゴリズム・データ構造の最適化方向
- **ベースライン確立**: 将来の性能改善測定基準

## 他実装との比較

### 比較対象

- Go標準 `map[string]interface{}`
- サードパーティトライ実装
- B-tree, Hash table等の代替データ構造

### 比較観点

- **メモリ使用量**: 同一データセットでの比較
- **操作性能**: Insert/Search/Delete性能
- **プレフィックス検索**: 前方一致性能（トライ特有の優位性）

## 継続的改善

### 自動化

- CI/CDでの定期ベンチマーク実行
- 性能劣化の早期検出
- 改善効果の定量測定

### 拡張性

- 新しいデータセット追加
- 追加ベンチマークシナリオ
- より大規模データでの測定（10M+キー）
