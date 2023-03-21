# meet-me-bot

Discord の特定のチャンネルに予定を通知するものです。

個人的な用途のために作成しているので、時刻が毎回 21 時に固定されているなど汎用的でない部分があります。

ChatGPT の手を借りて作成しています。

## 機能一覧

- 特定のユーザーが特定のチャンネルで「次回 3/21」のように書くことで、予定として記憶しておく。
  - 予定は 21 時固定。記憶できる日時はひとつまで
  - 記憶に成功した場合は 👍 失敗した場合は 👎 のリアクションをする
- 1 時間前および予定時刻には同チャンネルでお知らせをする。
- 「次いつ？」のような文言を書きこむと次回の予定を教えてくれる。

## install

```
pip install -r requirements.txt
```

## 起動

```
python meet-me-bot.py
```

## sqlite3 について

### インストール

```
sudo apt-get install sqlite3
```

### スキーマ作成

```
sqlite3 meet-me-bot.db < schema.sql
```

### 接続コマンド（デバッグ用）

```
sqlite3 meet-me-bot.db
```

## API キー取得・サーバーへの bot 追加方法

1. [Discord 開発者ポータル](https://discord.com/developers/applications)にアクセスし、Discord アカウントでログインする。

2. 「New Application」をクリックして、新しいアプリケーションを作成する。

3. 「General Information」ページで、アプリケーションの名前を入力して、「Create」をクリックする。

4. 「Bot」タブを選択して、Bot アカウントを作成する。このときトークンを取得しておくこと。

5. Bot の設定を変更し、必要に応じて Bot のアバターや名前を設定する。このとき Privileged Gateway Intents の「MESSAGE CONTENT INTENT」を有効にすること。

6. 「Oauth2」の「URL Generator」タブから、以下を選択して GERERATED URL をクリックし、bot をサーバーに追加するか聞かれるので「はい」を選択する。

- SCOPES
  - bot
- GENERAL PERMISSIONS
  - Read Messages/View Channels
- TEXT PERMISSIONS
  - Send Messages
  - Read Message History
  - Add Reactions

## チャンネル ID 取得方法

開発者モードをオンにした上でチャンネルを右クリックし、「ID をコピー」を選択

## Discord の開発者モードをオンにする方法

1. Discord アプリを開く。
2. 左下の「ユーザー設定」アイコン（歯車）をクリックする。
3. 「詳細設定」をクリックする。
4. 「開発者モード」をオンにする。
