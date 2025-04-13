# Go-Readability移植プロジェクトメモ

このドキュメントは、TypeScriptライブラリ [@mizchi/readability](https://github.com/mizchi/readability) をGoに移植する過程で必要な情報をまとめたものです。

## 環境設定

### Go バージョン

**Go 1.24.2** を使用します。Go 1.24の主な新機能：

- 実行ファイル依存関係の管理改善 (`tool` ディレクティブを使用)
- 新しいマップ実装（Swiss Tables ベース）- 2〜3%のCPUオーバーヘッド削減
- 標準ライブラリの改善（`os.Root`によるディレクトリ制限付きファイルシステムアクセスなど）
- 多くのパッケージで `encoding.BinaryAppender` と `encoding.TextAppender` インターフェースのサポート追加

詳細なリリースノートは [external-docs/go-docs/go1.24.md](external-docs/go-docs/go1.24.md) を参照してください。

### 開発ツール

移植作業で使用する主なツールは以下の通りです：

#### staticcheck

[staticcheck](https://staticcheck.io/) は、高度なGoコード静的解析ツールです。golangci-lintの代わりに使用します。
Go 1.24の新機能である `go get -tool` コマンドを使用してインストールします。

**特徴**:
- 100以上の検査項目
- パフォーマンス、セキュリティ、コーディングスタイル、バグなどの様々な問題を検出
- コンパイルエラーのない、微妙な問題も検出可能

#### gotests

[gotests](https://github.com/cweill/gotests) は、Goのソースコードを基にしてテーブル駆動テストを自動生成するツールです。

**特徴**:
- 関数やメソッドのインターフェースを解析して適切なテストを生成
- テーブル駆動テストのテンプレートをサポート
- エディタとの統合が容易（VS Code拡張などがある）

#### dlv (Delve)

[dlv (Delve)](https://github.com/go-delve/delve) は、Goプログラム用のデバッガです。

**特徴**:
- Goの言語仕様に合わせた設計
- ゴルーチンのデバッグサポート
- ブレークポイント、変数検査、スタックトレースなどの機能
- エディタとの統合（VS Code Go拡張で対応）

### 開発ツールのインストール

aquaはGo本体のインストールに使用し、その他の開発ツールはGo 1.24の新機能である `go get -tool` コマンドでインストールします。

```bash
# Go本体のインストール
aqua g golang/go@go1.24.2

# 開発ツールのインストール
go get -v -tool github.com/dominikh/go-tools/cmd/staticcheck
go get -v -tool github.com/cweill/gotests/...
go get -v -tool github.com/go-delve/delve/cmd/dlv
```

### 開発ツールの使用方法

Go 1.24から導入された `go get -tool` コマンドでインストールしたツールは、`go tool <ツール名>` の形式で使用します：

```bash
# staticcheckの使用例
go tool staticcheck ./...

# gotestsの使用例
go tool gotests -all ./path/to/file.go

# dlv (デバッガー)の使用例
go tool dlv debug ./path/to/main.go

# golangci-lintの使用例
go tool golangci-lint run
```

**重要**: golangci-lintを実行する際は、スタンドアロンのバイナリではなく、`go tool` コマンドを使用してください。これにより、Go 1.24のツールチェーンと統合されたlinterが使用されます。

これらのツールは従来の方法（$GOPATH/bin にインストールしてPATHに追加）と異なり、Go 1.24の新しいツール管理システムを使用しています。これによりバージョン管理とツールのインストールがより簡単になりました。

### プロジェクト構造

Go言語の慣習とTypeScriptの元コードの構造を考慮して、以下のディレクトリ構造を採用しています：

```
go-readability/                 # リポジトリルート - メインパッケージ
  ├── *.go                      # コアとなる機能を実装するGoファイル
  │                             # (aria.go, classify.go, core.go, dom.go, ...)
  ├── internal/                 # 非公開パッケージ（内部実装用）
  │   ├── dom/                  # DOM関連の内部実装
  │   ├── parser/               # HTML解析の内部実装
  │   └── util/                 # ユーティリティ関数
  ├── cmd/                      # コマンドラインアプリケーション用
  │   └── readability/          # CLIエントリーポイント
  └── testdata/                 # テスト関連ファイル（Goの慣例に従う）
      └── fixtures/             # テスト用フィクスチャ
```

この構造により、ライブラリの利用者は `import "github.com/mackee/go-readability"` とするだけで、メインの機能（例：`readability.Extract()`）にアクセスできます。

## 参考リソース

### 移植元プロジェクト

- @mizchi/readability: https://github.com/mizchi/readability
  - 現在のコードベースは `external-docs/mizchi-readability` にクローン済み

### 関連ライブラリとリファレンス

- [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html) - GoのHTMLパーサー
  - TypeScriptの`htmlparser2`に相当する役割を担う

### HTMLパース関連の比較情報

TypeScriptの`htmlparser2`とGolangの`golang.org/x/net/html`の主な違い：

| 機能 | htmlparser2 (TS) | golang.org/x/net/html (Go) |
|------|-----------------|----------------------------|
| パース方式 | ストリーミングパーサー | 完全なDOM構築 |
| イベント | イベントベース | Nodeツリー生成 |
| メモリ使用量 | 比較的少ない | DOM全体を保持 |
| 柔軟性 | 高い（カスタマイズ可能） | 標準的 |

## 実装方針メモ

### 基本戦略

1. まず基本的なデータ構造を定義
   - 仮想DOM構造体（VNode, VElement, VText）
   - ReadabilityOptions, ReadabilityArticle 
   
2. HTMLパーサーとノード操作関数を実装
   - golang.org/x/net/htmlを使用

3. 抽出ロジックを段階的に実装
   - 前処理
   - スコアリング
   - メインコンテンツ特定
   - 後処理

4. テスト作成
   - オリジナルテストケースの移植
   - 新しいテストケースの追加

## その他のメモ

（今後追加予定）
