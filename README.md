# simple-activitypub-spam-filter

- Mastodon/Misskey等ActivityPubを利用するBotアカウントに対応するためのスパムフィルターです
- リバースプロキシとして動作し、コンテンツに任意の文字列が含まれている場合に動作します
- URLを機械的に投稿するスパムに対し特に有効です

## Environment Values
すべての設定は環境変数により行います。

**BLOCK_WORDS**
```
BLOCK_WORDS=THE_EXAMPLE_SPAM_URL.org,EXAMPLE_WORDS
```

- ブロックしたいワードをカンマ区切りで指定します。
- ActivityPubのコンテンツ内を全文検索し、完全一致するワードが含まれていた場合にブロック対象とします。

**LISTEN_ADDRESS**
```
LISTEN_ADDRESS=:80
LISTEN_ADDRESS=0.0.0.0:8080
```

- 他のサーバーからの通信を待ち受けるアドレスを指定します。
- ポート単位でもアドレス単位でも利用できます。

**WHEN_DETECT_SPAM**
```
WHEN_DETECT_SPAM=output
WHEN_DETECT_SPAM=block
```

- スパムを検知したときにどのような動作を行うかを定義します。
- `output`: スパムを検知した場合、標準出力にContentの内容を出力し配送は続けます。
- `block`: スパムを検知した場合、配送元サーバーに400を送信し配送を取りやめます。

**PROXY_TARGET**
```
PROXY_TARGET=http://localhost
PROXY_TARGET=http://mastodon:8080
PROXY_TARGET=http://your-mastodon-apache.mastodon.svc.cluster.local
```

- 配送を行うサーバーのアドレスを指定します。
- 詳しい設定方法は下記の利用方法を参照してください。

## 利用方法
スパムフィルターはリバースプロキシとして動作します。ネットワークの入り口とMastodonサーバー・Misskeyサーバーの間にイメージを追加してください。

## Docker Composeで利用する場合

## Kubernetesで利用する場合

## オンプレミスで利用する場合

