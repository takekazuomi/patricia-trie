# Patricia Trie REPL

パトリシアトライを使用した対話的な前方一致検索ツール。

## 概要

このREPLツールは、ワードリストファイルを読み込んでパトリシアトライを構築し、対話的に前方一致検索を実行できます。go-promptライブラリを使用したリアルタイム補完機能と、Emacsキーバインドによる快適な操作性が特徴です。

## インストール

```bash
# ビルド
go build -o patricia-repl cmd/patricia-repl/main.go

# または直接実行
go run cmd/patricia-repl/main.go <wordlist>
```

## 使用方法

### 基本的な使い方

```bash
./patricia-repl wordlist.txt
```

### 起動時の動作

1. **履歴の自動読み込み**
   - `~/.config/patricia-repl/history` から過去の検索履歴を自動読み込み
   - 終了時に履歴を自動保存（最大1000件）

2. **リアルタイム補完**
   - Tabキーで補完候補を表示
   - トライ内の単語を動的に補完

### REPLコマンド

- **前方一致検索**: 任意の文字列を入力
- **/verbose**: Verboseモードのオン/オフ切り替え
- **/help**: ヘルプメッセージとキーバインド一覧を表示
- **/exit/quit**: REPLを終了
- **Ctrl+D**: REPLを終了（EOF）

## 実行例

### 基本操作

```bash
$ ./patricia-repl testwords.txt
📚 Loaded 32 words from testwords.txt

Patricia Trie REPL started. Commands: /exit, /quit, /verbose, /help
Use Tab for auto-completion. Emacs key bindings are enabled.

> ca
✓ Found 2 words: cat, cats
> dog
✓ Found 2 words: dog, dogs
> xyz
✗ No matches found for prefix 'xyz'
> /verbose
[info] Verbose mode enabled
> ca
✓ Found 3 words: cat, cats, cattle
  [verbose] Nodes visited: 5, Max depth: 2, Time: 0.125ms
> /help
(ヘルプメッセージとキーバインド一覧が表示)
> /exit

👋 Goodbye!
```

### リアルタイム補完の例

```bash
> /[TAB]
/help     /verbose  /exit     /quit     (コマンドの補完候補)

> /ver[TAB]
> /verbose  (自動補完される)
```

## ワードリストファイルの形式

- 1行に1単語
- 空行は無視される
- UTF-8エンコーディング推奨

### 例: testwords.txt

```text
cat
cats
cattle
dog
dogs
dolphin
elephant
eagle
```

## 機能詳細

### 検索機能
- 入力されたプレフィックスに一致するすべての単語を検索
- 大文字小文字を区別
- 結果は件数に関わらずすべて表示

### リアルタイム補完
- /で始まる入力はコマンドのみを補完
- 通常の文字列入力では補完を行わない（検索のみ）
- Tabキーで候補を表示・選択

### Verboseモード統計情報
- **Nodes visited**: 探索したノード数（現在は仮実装）
- **Max depth**: 探索の最大深度（現在は仮実装）
- **Time**: 検索処理時間（ミリ秒）

### 履歴機能
- 検索履歴を自動保存（最大1000件）
- 保存場所: `~/.config/patricia-repl/history`
- 起動時に自動読み込み

### Emacsキーバインド
- 標準的なEmacsキーバインドをサポート
- カーソル移動、文字削除、単語移動など
- `help`コマンドで一覧表示

## 今後の拡張予定

- [ ] 実際のノード探索統計の実装（patriciatrieパッケージの拡張）
- [x] 検索履歴機能（実装済み）
- [x] コマンド専用補完機能（実装済み）
- [ ] 履歴の逆検索機能（Ctrl+R）
- [ ] 大量結果時のページング
- [ ] 結果のエクスポート機能
- [ ] 設定ファイルサポート

## トラブルシューティング

### Q: 日本語の単語が正しく表示されない
A: ターミナルのエンコーディングがUTF-8に設定されているか確認してください。

### Q: カラー表示が機能しない
A: ターミナルがANSIカラーコードに対応している必要があります。Windowsの場合は、Windows Terminal使用を推奨。

### Q: 補完候補が表示されない
A: Tabキーを押しても候補が表示されない場合は、ターミナルエミュレータの設定を確認してください。一部のターミナルではTabキーが別の機能に割り当てられている場合があります。

### Q: Emacsキーバインドが効かない
A: ターミナルエミュレータによってはAltキーの扱いが異なります。設定でMetaキーの送信を有効にしてください。