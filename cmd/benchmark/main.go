package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/mackee/go-readability"
)

func main() {
	// コマンドライン引数の解析
	var (
		cpuprofile  = flag.String("cpuprofile", "", "CPUプロファイルの出力先ファイル")
		memprofile  = flag.String("memprofile", "", "メモリプロファイルの出力先ファイル")
		iterations  = flag.Int("iterations", 100, "繰り返し回数")
		htmlFile    = flag.String("html", "", "HTMLファイルのパス（指定しない場合はテストケースを使用）")
		testCaseDir = flag.String("testcase", "../../testdata/fixtures/001", "テストケースのディレクトリ（htmlが指定されていない場合に使用）")
	)
	flag.Parse()

	// CPUプロファイリングの設定
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "CPUプロファイルの作成に失敗しました: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "CPUプロファイリングの開始に失敗しました: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	// HTMLの読み込み
	var html string
	var err error
	if *htmlFile != "" {
		// 指定されたHTMLファイルを読み込む
		htmlBytes, err := ioutil.ReadFile(*htmlFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "HTMLファイルの読み込みに失敗しました: %v\n", err)
			os.Exit(1)
		}
		html = string(htmlBytes)
	} else {
		// テストケースのHTMLを読み込む
		sourcePath := filepath.Join(*testCaseDir, "source.html")
		htmlBytes, err := ioutil.ReadFile(sourcePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "テストケースのHTMLファイルの読み込みに失敗しました: %v\n", err)
			os.Exit(1)
		}
		html = string(htmlBytes)
	}

	// オプションの設定
	options := readability.ReadabilityOptions{}

	// 処理時間の測定
	startTime := time.Now()
	var result readability.ReadabilityArticle

	// 指定された回数だけ抽出を実行
	for i := 0; i < *iterations; i++ {
		result, err = readability.Extract(html, options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "抽出に失敗しました: %v\n", err)
			os.Exit(1)
		}
	}

	// 処理時間の表示
	elapsedTime := time.Since(startTime)
	fmt.Printf("処理時間: %v (%d回の繰り返し)\n", elapsedTime, *iterations)
	fmt.Printf("1回あたりの平均処理時間: %v\n", elapsedTime/time.Duration(*iterations))

	// 結果の表示
	fmt.Printf("タイトル: %s\n", result.Title)
	fmt.Printf("著者: %s\n", result.Byline)
	fmt.Printf("ページタイプ: %s\n", result.PageType)
	fmt.Printf("ノード数: %d\n", result.NodeCount)

	// メモリ使用量の表示
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)

	// メモリプロファイリングの設定
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "メモリプロファイルの作成に失敗しました: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		runtime.GC() // メモリプロファイリング前にGCを実行
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "メモリプロファイルの書き込みに失敗しました: %v\n", err)
			os.Exit(1)
		}
	}
}

// バイトをメガバイトに変換する関数
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
