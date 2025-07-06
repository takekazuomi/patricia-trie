# テストデータディレクトリ

このディレクトリには大規模ベンチマーク用のテストデータが格納されます。

## ディレクトリ構成

```
testdata/
├── japanese/           # 日本語辞書データ
│   ├── sudachi_core.txt      # 辞書エントリ（100K語）
│   └── sudachi_extended.txt  # 拡張辞書（1M語）
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
- **IPアドレス**: ランダム生成（IPv4/IPv6対応）