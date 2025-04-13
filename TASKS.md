# @mizchi/readability → Go移植タスクリスト

このファイルは、TypeScriptで書かれた@mizchi/readabilityライブラリをGoに移植するための具体的なタスクリストです。各タスクは可能な限り小さく、テスト可能な粒度で設計されています。

## 1. 環境セットアップ

- [x] リポジトリのクローン
- [x] aquaを使用して必要なツールをインストール
  - [x] golangci-lint (lintツール)
  - [x] go get -tool honnef.co/go/tools/cmd/staticcheck@latest (静的解析ツール)
  - [x] go get -tool github.com/cweill/gotests/...@latest (テスト生成ツール)
  - [x] go get -tool github.com/go-delve/delve/cmd/dlv@latest (デバッガー)
- [x] Go 1.24.2の確認
- [x] 基本的なディレクトリ構造の作成

## 2. 基本データ構造の実装

- [x] **2.1. 仮想DOM構造体の定義**
  - [x] VNodeType列挙型の実装
  - [x] VNode基本インターフェースの実装
  - [x] VText構造体の実装
  - [x] VElement構造体の実装
  - [x] VDocument構造体の実装
  - [x] テスト: 構造体の基本操作

- [x] **2.2. ReadabilityOptions構造体の実装**
  - [x] 基本オプション構造体の定義
  - [x] デフォルト値の設定
  - [x] テスト: オプション動作確認

- [x] **2.3. ReadabilityArticle構造体の実装**
  - [x] 記事データ構造体の定義
  - [x] PageType列挙型の実装
  - [x] テスト: 基本構造確認

## 3. パーサーの実装

- [x] **3.1. HTML解析基盤**
  - [x] golang.org/x/net/htmlの導入
  - [x] HTMLノードとVNodeの変換関数の実装
  - [x] テスト: 基本的なHTML解析

- [x] **3.2. ParseHTML関数の実装**
  - [x] HTML文字列をVDocumentに変換する機能
  - [x] テスト: 様々なHTMLの解析

- [x] **3.3. SerializeToHTML関数の実装**
  - [x] VDocumentをHTML文字列に変換する機能
  - [x] テスト: シリアライズの検証

## 4. ユーティリティ関数の実装

- [x] **4.1. DOM操作関数**
  - [x] getElementsByTagName実装
  - [x] isProbablyVisible実装
  - [x] getNodeAncestors実装
  - [x] createElement実装
  - [x] テスト: 各関数の動作確認

- [x] **4.2. テキスト関連関数**
  - [x] getInnerText実装
  - [x] getLinkDensity実装
  - [x] getTextDensity実装
  - [x] テスト: テキスト処理の確認

- [x] **4.3. 正規表現定数**
  - [x] REGEXPS定数群の実装
  - [x] DEFAULT_TAGS_TO_SCORE定数の実装
  - [x] テスト: 正規表現マッチング

## 5. コア機能の実装

- [x] **5.1. 基本スコアリング**
  - [x] initializeNode関数の実装
  - [x] getClassWeight関数の実装
  - [x] テスト: スコアリングロジック

- [x] **5.2. 候補抽出**
  - [x] findMainCandidates関数の実装
  - [x] isProbablyContent関数の実装
  - [x] テスト: 候補抽出の検証

- [x] **5.3. メタデータ抽出**
  - [x] getArticleTitle関数の実装
  - [x] getArticleByline関数の実装
  - [x] テスト: メタデータ抽出

- [x] **5.4. ページ分類**
  - [x] classifyPageType関数の実装
  - [x] isSignificantNode関数の実装
  - [x] isSemanticTag関数の実装
  - [x] テスト: ページ分類の検証

- [x] **5.5. 構造要素検出**
  - [x] findStructuralElements関数の実装
  - [x] addSignificantElementsByClassOrId関数の実装
  - [x] テスト: 構造要素検出

## 6. 前処理・後処理の実装

- [x] **6.1. 前処理**
  - [x] preprocessDocument関数の実装
  - [x] 不要な要素の除去
  - [x] テスト: 前処理の効果

- [x] **6.2. フォーマット**
  - [x] toHTML関数の実装
  - [x] stringify関数の実装
  - [x] formatDocument関数の実装
  - [x] extractTextContent関数の実装
  - [x] countNodes関数の実装
  - [x] テスト: フォーマット機能

## 7. メイン抽出機能の実装

- [x] **7.1. extractContent関数**
  - [x] コンテンツ抽出メイン処理
  - [x] テスト: 抽出ロジック

- [x] **7.2. extract関数**
  - [x] HTMLからの抽出エントリーポイント
  - [x] 各種オプション処理
  - [x] テスト: 実際のHTML文字列からの抽出

- [x] **7.3. createExtractor関数**
  - [x] カスタム抽出関数ファクトリ
  - [x] テスト: カスタム抽出の動作

## 8. アクセシビリティ機能の実装

- [x] **8.1. ARIA構造体定義**
  - [x] AriaNodeType列挙型
  - [x] AriaNode構造体
  - [x] AriaTree構造体
  - [x] テスト: ARIA構造体

- [x] **8.2. ARIA解析基本関数**
  - [x] buildAriaNode関数の実装
  - [x] countAriaNodes関数の実装
  - [x] テスト: ARIA解析基本機能

- [x] **8.3. ARIA主要関数**
  - [x] buildAriaTree関数の実装
  - [x] ariaTreeToString関数の実装
  - [x] extractAriaTree関数の実装
  - [x] テスト: ARIA主要機能

## 9. Markdown変換の実装

- [x] **9.1. Markdown変換**
  - [x] toMarkdown関数の実装
  - [x] テスト: HTML→Markdown変換

## 10. 総合テスト

- [x] **10.1. 実サイト抽出テスト**
  - [x] 元ライブラリのテストケース移植
  - [x] 実際のウェブサイトからの抽出テスト

- [x] **10.2. パフォーマンステスト**
  - [x] 処理速度の測定
  - [x] メモリ使用量の測定

## 11. ドキュメンテーション

- [x] **11.1. GoDoc用コメント**
  - [x] パッケージとエクスポートされた関数にコメント追加
  - [x] 各構造体と型定義のドキュメント

- [x] **11.2. 使用方法ドキュメント**
  - [x] README.mdの作成
  - [x] 使用例の追加

## 12. CI/CD設定

- [x] **12.1. テスト自動化**
  - [x] GitHub Actions設定
  - [x] テストワークフロー

- [x] **12.2. リリース自動化**
  - [x] リリースワークフロー
  - [x] バージョニング

## 優先順位とロードマップ

このタスクリストは先頭から順番に進めることを推奨します。各セクションの中で、依存関係によって順序が前後する場合があります。

### 最小実装（MVP）の範囲

上記タスクのうち、以下を実装することでMVP（必要最低限の機能）が完成します：

1. 環境セットアップ
2. 基本データ構造の実装
3. パーサーの実装
4. ユーティリティ関数の実装
5. コア機能の実装
6. 前処理の実装
7. メイン抽出機能の実装

### 拡張機能の範囲

MVP完成後に追加実装する機能：

8. アクセシビリティ機能の実装
9. Markdown変換の実装
10. 総合テスト
11. ドキュメンテーション
12. CI/CD設定
