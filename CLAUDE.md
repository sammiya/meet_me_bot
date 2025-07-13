# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

Discord の特定のチャンネルに予定を通知する Go 製の Bot です。特定のユーザーが「次回 3/21」のような書式で予定を設定すると、1時間前と予定時刻に通知を送信します。

## 開発コマンド

### 実行・ビルド
```bash
go run .                    # 直接実行（.env ファイルが必要）
go build -v .               # バイナリをビルド（配布用）
./meet_me_bot              # ビルド済みバイナリの実行
```

### テスト・リント
```bash
go test -v ./...            # 全テストの実行
staticcheck ./...           # リンター（ST1000 チェックを除く）
```

### データベース操作
```bash
sqlite3 meet-me-bot.db      # SQLite データベースに接続（デバッグ用）
```

## 環境設定

実行には以下の環境変数が必要です（.env ファイルで設定）：
- `DISCORD_TOKEN`: Discord Bot のトークン
- `DISCORD_CHANNEL`: 監視・投稿対象のチャンネル ID
- `DISCORD_MEETING_ANNOUNCER`: 予定設定権限を持つユーザー ID

## アーキテクチャ

### メインロジック（main.go）
- Discord API（discordgo）を使用したメッセージハンドリング
- 5秒間隔のティッカーで定期的な通知チェック（本番は1分間隔を想定）
- 特定ユーザーからの「次回 M/D」形式のメッセージを予定として記録
- 「次いつ？」形式のメッセージに予定情報で応答

### データ管理（utils/db.go）
- SQLite を使用した単一テーブル（meetings）での予定管理
- 同時に存在できる予定は1つのみ（新規追加時に既存削除）
- 1時間前通知の送信状態も管理

### 日付処理（utils/date.go）
- 正規表現による日付文字列のパース（M/D 形式）
- JST タイムゾーンでの時刻処理（21時固定）
- 日付フォーマット（M/D（曜日）形式）

### テスト構成
- 各ユーティリティ関数のユニットテスト（utils/*_test.go）
- 日付パースとデータベース操作のテスト

### CI/CD
- GitHub Actions で自動テスト・リント・リリース
- 静的解析は staticcheck を使用（ST1000 除く）
- プッシュ時に自動バイナリリリース作成