# テストデータディレクトリ

このディレクトリには大規模ベンチマーク用のテストデータが格納されます。

## ディレクトリ構成

```text
testdata/
├── japanese/           # 日本語辞書データ
│   ├── small_lex.zip         # small_lex辞書（zip）
│   ├── core_lex.zip          # core_lex辞書（zip）
│   ├── notcore_lex.zip       # notcore_lex辞書（zip）
│   ├── small_lex.csv         # small_lex辞書（CSV、展開済み）
│   ├── core_lex.csv          # core_lex辞書（CSV、展開済み）
│   ├── notcore_lex.csv       # notcore_lex辞書（CSV、展開済み）
│   ├── small.txt             # small_lex見出し語のみ（約57万語）
│   ├── core.txt              # core_lex見出し語のみ（約82万語）
│   ├── notcore.txt           # notcore_lex見出し語のみ（約124万語）
│   ├── 1000.txt              # テスト用漢字語彙（1000語）
│   └── full.txt              # 全辞書統合（重複削除、約259万語）
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
  - **small_lex.csv**: 基本語彙（約50万語、40MB）
  - **core_lex.csv**: 追加語彙（約200万語、21MB）
  - **notcore_lex.csv**: 専門語彙・固有名詞（約300万語、35MB）
- **IPアドレス**: ランダム生成（IPv4/IPv6対応）

### 辞書ファイルの関係性

**重要**: 各CSVファイルは独立した辞書ファイルです。small_lexはcore_lexのサブセットではありません。

- **独立性**: 各辞書ファイルは独立した語彙を持つ
- **累積構造**: Core辞書 = small_lex + core_lex、Full辞書 = small_lex + core_lex + notcore_lex
- **重複削除**: 統合時は`sort -u`による重複削除が必要

根拠：

- [SudachiDict公式](https://github.com/WorksApplications/SudachiDict): 「Core dictionary requires small and core files」
- [Sudachi公式](https://github.com/WorksApplications/Sudachi): 辞書の累積構造を説明

### 実データ分析結果（20250515版）

実際のsudachi辞書データを使用した重複分析：

| 項目 | small_lex.csv | core_lex.csv | 重複状況 |
|------|---------------|--------------|----------|
| 総行数 | 765,611行 | 861,167行 | 完全重複: 3行 |
| 語彙数 | 573,818語 | 819,383語 | 語彙重複: 17,607語 |

**重複例**:

- **完全重複**: サマセット、大庭、富士吉田（地名・人名）
- **語彙重複**: 数字、記号、Unicode文字等

**検証方法**:

```bash
# 重複確認
comm -12 <(sort small_lex.csv) <(sort core_lex.csv)

# 重複数
comm -12 <(sort small_lex.csv) <(sort core_lex.csv) | wc -l
```

**結論**: 各辞書の独立性が実証され、統合時の重複削除処理の必要性が確認されました。

## ベンチマークデータセット

### 日本語辞書データ

| ファイル | 内容 | 語数 | 用途 |
|---------|------|------|------|
| `1000.txt` | 漢字で始まる単語（テスト用） | 1,000語 | 単体テスト |
| `small.txt` | small_lex見出し語のみ | 約57万語 | 基本ベンチマーク |
| `core.txt` | core_lex見出し語のみ | 約82万語 | 中規模ベンチマーク |
| `notcore.txt` | notcore_lex見出し語のみ | 約124万語 | 大規模ベンチマーク |
| `full.txt` | 全辞書統合（重複削除） | 約259万語 | 超大規模ベンチマーク |

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

- **Core**: core_lex見出し語のみ（約82万語）
- **NotCore**: notcore_lex見出し語のみ（約124万語）
- **Full**: 全辞書統合（約259万語）

各データセットで挿入、検索、プレフィックス検索の性能を測定します。

### データ形式の統一

- **TXTファイル**: 見出し語のみ（重複削除済み）
- **CSVファイル**: 元の辞書データ（品詞、読み等を含む18フィールド）

ベンチマークでは統一的に見出し語のみのTXTファイルを使用し、CSV解析のオーバーヘッドを排除した純粋な語彙性能を測定します。
