# テストデータディレクトリ

このディレクトリには大規模ベンチマーク用のテストデータが格納されます。

## ディレクトリ構成

```text
testdata/
├── japanese/           # 日本語辞書データ
│   ├── small_words.txt       # small_lex.csvから抽出（基本語彙、約50万語）
│   ├── core_words.txt        # core_lex.csvから抽出（追加語彙、約200万語）
│   ├── notcore_words.txt     # notcore_lex.csvから抽出（専門語彙、約300万語）
│   ├── full_words.txt        # 全辞書統合（重複削除、約800万語）
│   ├── 1000.csv              # テスト用（1000語）
│   ├── all.csv               # ベンチマーク用（small_lexのみ）
│   ├── large_bench.csv       # 大規模ベンチマーク用（small+core、約250万語）
│   └── mega_bench.csv        # 超大規模ベンチマーク用（全辞書統合、約800万語）
└── ipaddresses/        # IPアドレスデータ
    ├── ipv4_10k.txt          # IPv4 10K個
    ├── ipv4_100k.txt         # IPv4 100K個
    ├── ipv6_10k.txt          # IPv6 10K個
    └── ipv6_100k.txt         # IPv6 100K個
```

## データの準備

テストデータは以下のコマンドで自動生成されます：

```bash
make setup_benchmark
```

## データソース

- **日本語辞書**: [Sudachi Language Resources](https://registry.opendata.aws/sudachi/) (Apache-2.0ライセンス)
  - **配布URL**: <https://d2ej7fkh96fzlu.cloudfront.net/sudachidict-raw/>
  - **最新版**: 20250515
  - **公式リポジトリ**: [SudachiDict](https://github.com/WorksApplications/SudachiDict)、[Sudachi](https://github.com/WorksApplications/Sudachi)
  - **small_lex.csv**: 基本語彙（約50万語、40MB）- UniDicベース
  - **core_lex.csv**: 追加語彙（約200万語、21MB）- NEologdベース
  - **notcore_lex.csv**: 専門語彙（約300万語、35MB）- 固有名詞・専門用語
- **IPアドレス**: ランダム生成（IPv4/IPv6対応）

### 辞書ファイルの関係性

**重要**: 各CSVファイルは独立した辞書ファイルです。small_lexはcore_lexのサブセットではありません。

- **独立性**: 各辞書ファイルは独立した語彙を持つ
- **累積構造**: Core辞書 = small_lex + core_lex、Full辞書 = small_lex + core_lex + notcore_lex
- **重複削除**: 統合時は`sort -u`による重複削除が必要

根拠：
- [SudachiDict公式](https://github.com/WorksApplications/SudachiDict): 「Core dictionary requires small and core files」
- [Sudachi公式](https://github.com/WorksApplications/Sudachi): 辞書の累積構造を説明

## ベンチマークデータセット

### 日本語辞書データ

| ファイル | 内容 | 語数 | 用途 |
|---------|------|------|------|
| `1000.csv` | 漢字で始まる単語（テスト用） | 1,000語 | 単体テスト |
| `all.csv` | small_lex辞書 | 約50万語 | 基本ベンチマーク |
| `large_bench.csv` | small+core辞書 | 約250万語 | 大規模ベンチマーク |
| `mega_bench.csv` | 全辞書統合 | 約800万語 | 超大規模ベンチマーク |

### IPアドレスデータ

| ファイル | 内容 | 個数 | 用途 |
|---------|------|------|------|
| `ipv4_10k.txt` | IPv4アドレス | 10,000個 | 中規模ベンチマーク |
| `ipv4_100k.txt` | IPv4アドレス | 100,000個 | 大規模ベンチマーク |
| `ipv6_10k.txt` | IPv6アドレス | 10,000個 | 中規模ベンチマーク |
| `ipv6_100k.txt` | IPv6アドレス | 100,000個 | 大規模ベンチマーク |

## 使用方法

### 基本ベンチマーク実行

```bash
# 基本ベンチマーク
make benchmark

# 大規模データベンチマーク
make benchmark-large

# 実世界データベンチマーク
make benchmark-realistic
```

### 大規模データベンチマーク

`BenchmarkTrie_Large_Japanese_Specialized`では、以下の大規模データセットを使用：

- **Core_250W**: small+core辞書（約250万語）
- **Full_800W**: 全辞書統合（約800万語）

各データセットで挿入、検索、プレフィックス検索の性能を測定します。
