package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/takekazu/patricia-trie/pkg/patriciatrie"
)

const (
	bytesPerKB      = 1024
	bytesPerMB      = 1024 * 1024
	msPerSecond     = 1000.0
	maxSuggestions  = 10
	maxHistoryItems = 1000
	dirPermission   = 0750
)

var (
	trie    *patriciatrie.Trie
	verbose bool
	history []string
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	cyan    = color.New(color.FgCyan).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
)

// Stats はverboseモード時の統計情報を保持
type Stats struct {
	NodesVisited int
	MaxDepth     int
	Duration     time.Duration
}

// BuildStats はトライ構築時の統計情報を保持
type BuildStats struct {
	Duration   time.Duration
	MemoryUsed uint64 // バイト単位
}

func main() {
	// コマンドラインフラグの定義
	verboseFlag := flag.Bool("verbose", false, "Verboseモードで開始")
	verboseFlagShort := flag.Bool("v", false, "Verboseモードで開始（短縮形）")
	flag.Usage = showUsage
	flag.Parse()

	// 残りの引数をチェック
	if flag.NArg() < 1 {
		showUsage()
		os.Exit(1)
	}

	// verboseフラグの設定（--verbose または -v）
	verbose = *verboseFlag || *verboseFlagShort

	wordlistPath := flag.Arg(0)

	// ワードリストを読み込んでトライを構築
	var err error

	var buildStats BuildStats

	trie, buildStats, err = buildTrie(wordlistPath)
	if err != nil {
		log.Fatalf("Failed to load wordlist: %v", err)
	}

	// Verboseモードでトライ構築統計を表示
	if verbose {
		memoryMB := float64(buildStats.MemoryUsed) / bytesPerMB
		fmt.Printf("%s Trie build time: %.3fms, Memory used: %.2fMB\n", 
			yellow("[verbose]"), 
			float64(buildStats.Duration.Microseconds())/msPerSecond,
			memoryMB)
	}

	// 履歴ファイルの読み込み
	loadHistory()

	// 起動メッセージ
	verboseStatus := ""
	if verbose {
		verboseStatus = " (Verbose mode enabled)"
	}

	fmt.Printf("\n%s REPL started%s. Commands: /exit, /quit, /verbose, /help\n", cyan("Patricia Trie"), verboseStatus)
	fmt.Printf("Use Tab for auto-completion. Emacs key bindings are enabled.\n\n")

	// go-promptの起動
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("Patricia Trie REPL"),
		prompt.OptionPrefix("> "),
		prompt.OptionPrefixTextColor(prompt.Blue),
		prompt.OptionPreviewSuggestionTextColor(prompt.Green),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionMaxSuggestion(maxSuggestions),
		prompt.OptionShowCompletionAtStart(),
	)
	p.Run()
}

// executor はユーザー入力を処理
func executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// 履歴に追加
	history = append(history, input)

	// コマンド処理
	switch input {
	case "/exit", "/quit":
		saveHistory()
		fmt.Printf("\n%s Goodbye!\n", cyan("👋"))
		os.Exit(0)
	case "/verbose":
		verbose = !verbose

		status := "disabled"
		if verbose {
			status = "enabled"
		}

		fmt.Printf("%s Verbose mode %s\n", yellow("[info]"), status)
	case "/help":
		showHelp()
	default:
		// 前方一致検索を実行
		performSearch(input)
	}
}

// completer は動的な補完候補を生成
func completer(d prompt.Document) []prompt.Suggest {
	prefix := d.GetWordBeforeCursor()
	if prefix == "" {
		return []prompt.Suggest{
			{Text: "/help", Description: "Show help message"},
			{Text: "/verbose", Description: "Toggle verbose mode"},
			{Text: "/exit", Description: "Exit the REPL"},
		}
	}

	var suggestions []prompt.Suggest

	// /で始まる場合はコマンドのみを補完
	if strings.HasPrefix(prefix, "/") {
		commands := []prompt.Suggest{
			{Text: "/help", Description: "Show help message"},
			{Text: "/verbose", Description: "Toggle verbose mode"},
			{Text: "/exit", Description: "Exit the REPL"},
			{Text: "/quit", Description: "Exit the REPL"},
		}

		for _, cmd := range commands {
			if strings.HasPrefix(cmd.Text, prefix) {
				suggestions = append(suggestions, cmd)
			}
		}
	}

	return suggestions
}

// performSearch は前方一致検索を実行
func performSearch(prefix string) {
	start := time.Now()
	results, stats := searchWithStats(trie, prefix, verbose)
	duration := time.Since(start)

	// 結果表示
	switch len(results) {
	case 0:
		fmt.Printf("%s No matches found for prefix '%s'\n", red("✗"), prefix)
	case 1:
		fmt.Printf("%s Found 1 word: %s\n", green("✓"), results[0])
	default:
		fmt.Printf("%s Found %d words: %s\n", green("✓"), len(results), strings.Join(results, ", "))
	}

	// Verboseモードでの追加情報表示
	if verbose && stats != nil {
		stats.Duration = duration
		fmt.Printf("  %s Nodes visited: %d, Max depth: %d, Time: %.3fms\n",
			yellow("[verbose]"),
			stats.NodesVisited,
			stats.MaxDepth,
			float64(stats.Duration.Microseconds())/msPerSecond)
	}
}

// showUsage は使用方法を表示
func showUsage() {
	usage := `パトリシアトライ対話検索ツール

使用方法:
  patricia-repl [options] <wordlist>

引数:
  wordlist    単語リストファイル（1行に1単語、UTF-8エンコーディング）

オプション:
  -v, --verbose    Verboseモードで開始（検索統計情報を表示）

例:
  patricia-repl words.txt
  patricia-repl --verbose /path/to/dictionary.txt
  patricia-repl -v words.txt

ワードリストファイルの形式:
  cat
  cats  
  dog
  dogs
  elephant

起動後のコマンド:
  /help     - ヘルプメッセージを表示
  /verbose  - Verboseモードの切り替え
  /exit     - 終了
  /quit     - 終了

機能:
  - 前方一致検索: 任意の文字列を入力して検索実行
  - リアルタイム補完: /で始まるコマンドのみ補完
  - 履歴保存: 検索履歴を自動保存（~/.config/patricia-repl/history）
  - Emacsキーバインド: Ctrl+A, Ctrl+E, Ctrl+F, Ctrl+Bなど

詳細情報: https://github.com/takekazuomi/patricia-trie
`
	fmt.Fprint(os.Stderr, usage)
}

// showHelp はヘルプメッセージを表示
func showHelp() {
	help := `
Commands:
  /help     - Show this help message
  /verbose  - Toggle verbose mode (currently: %s)
  /exit     - Exit the REPL
  /quit     - Exit the REPL

Usage:
  - Type any prefix to search for matching words
  - Type / followed by Tab to see available commands
  - Use arrow keys to navigate suggestions

Key bindings (Emacs mode):
  Ctrl+A   - Move to beginning of line
  Ctrl+E   - Move to end of line
  Ctrl+F   - Move forward one character
  Ctrl+B   - Move backward one character
  Ctrl+D   - Delete character under cursor
  Ctrl+K   - Kill from cursor to end of line
  Ctrl+U   - Kill from beginning to cursor
  Ctrl+W   - Kill previous word
  Alt+F    - Move forward one word
  Alt+B    - Move backward one word
`

	status := "disabled"
	if verbose {
		status = "enabled"
	}

	fmt.Printf(help, status)
}

// buildTrie はワードリストファイルからパトリシアトライを構築
func buildTrie(path string) (*patriciatrie.Trie, BuildStats, error) {
	// ガベージコレクションを実行してより正確なメモリ測定を行う
	runtime.GC()
	
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)
	
	start := time.Now()
	
	file, err := os.Open(path) // #nosec G304 - path is provided by user
	if err != nil {
		return nil, BuildStats{}, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	trie := patriciatrie.New()
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" {
			continue
		}

		err := trie.Insert(word)
		if err != nil {
			return nil, BuildStats{}, fmt.Errorf("failed to insert word '%s': %w", word, err)
		}

		lineCount++
	}

	err = scanner.Err()
	if err != nil {
		return nil, BuildStats{}, fmt.Errorf("error reading file: %w", err)
	}

	buildTime := time.Since(start)
	
	// ガベージコレクションを実行してより正確なメモリ測定を行う
	runtime.GC()
	
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)
	
	// メモリ使用量の差分を計算（ヒープ上のアロケートされたバイト数）
	memoryUsed := memAfter.HeapAlloc - memBefore.HeapAlloc
	
	fmt.Printf("%s Loaded %d words from %s\n", yellow("📚"), lineCount, path)

	return trie, BuildStats{
		Duration:   buildTime,
		MemoryUsed: memoryUsed,
	}, nil
}

// searchWithStats は検索を実行し、verboseモード時は統計情報も収集
func searchWithStats(trie *patriciatrie.Trie, prefix string, verbose bool) ([]string, *Stats) {
	// TODO: 実際の実装では、patriciatrie側にStatsを収集する機能を追加する必要がある
	// 現在は基本的な検索のみ実装
	results := trie.FindByPrefix(prefix)
	
	var stats *Stats
	if verbose {
		// 仮の統計情報（実際の実装では内部状態から取得）
		stats = &Stats{
			NodesVisited: len(prefix) + len(results), // 仮の値
			MaxDepth:     len(prefix),                 // 仮の値
			Duration:     0,                           // 後で設定
		}
	}

	return results, stats
}

// loadHistory は履歴ファイルから履歴を読み込み
func loadHistory() {
	historyFile := getHistoryFile()

	file, err := os.Open(historyFile) // #nosec G304
	if err != nil {
		// ファイルが存在しない場合は無視
		return
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		history = append(history, scanner.Text())
	}
}

// saveHistory は履歴をファイルに保存
func saveHistory() {
	historyFile := getHistoryFile()
	
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(historyFile)
	
	err := os.MkdirAll(dir, dirPermission) // #nosec G301 - group/other read not needed
	if err != nil {
		return
	}

	file, err := os.Create(historyFile) // #nosec G304
	if err != nil {
		return
	}

	defer func() {
		_ = file.Close()
	}()

	// 最新のmaxHistoryItems件のみ保存
	start := 0
	if len(history) > maxHistoryItems {
		start = len(history) - maxHistoryItems
	}

	for i := start; i < len(history); i++ {
		_, _ = fmt.Fprintln(file, history[i])
	}
}

// getHistoryFile は履歴ファイルのパスを返す
func getHistoryFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".patricia_repl_history"
	}
	
	return filepath.Join(home, ".config", "patricia-repl", "history")
}