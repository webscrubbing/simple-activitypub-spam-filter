# simple-activitypub-spam-filter

[webscrubbing808/simple-activitypub-spam-filter - Docker Image | Docker Hub](https://hub.docker.com/r/webscrubbing808/simple-activitypub-spam-filter)

- Mastodon/Misskey等ActivityPubを利用するBotアカウントに対応するためのスパムフィルターです
- リバースプロキシとして動作し、コンテンツに任意の文字列が含まれている場合に動作します
- 特定のURLを機械的に投稿するスパムに対し特に有効です

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
WHEN_DETECT_SPAM=soft
```

- スパムを検知したときにどのような動作を行うかを定義します。
- `output`: スパムを検知した場合、標準出力にContentの内容を出力し配送は続けます。
- `block`: スパムを検知した場合、配送元サーバーに400を送信し配送を取りやめます。送信元サーバーにスパムを送信していることを通知できます。
- `soft`: スパムを検知した場合、配送元サーバーに200を送信し配送を取りやめます。再送を防ぐのに有効です。

**PROXY_TARGET**
```
PROXY_TARGET=http://localhost:3000
PROXY_TARGET=http://mastodon:8080
PROXY_TARGET=http://your-mastodon-apache.mastodon.svc.cluster.local
```

- 配送を行うサーバーのアドレスを指定します。
- 詳しい設定方法は下記の利用方法を参照してください。

## 利用方法
スパムフィルターはリバースプロキシとして動作します。ネットワークの入り口とMastodonサーバー・Misskeyサーバーの間にイメージを追加してください。

## Docker Composeで利用する場合

下記のような構成で運用している場合のサンプルです。

```
services:
  mastodon-db:
    restart: always
    image: 'postgres:alpine'

  mastodon-redis:
    restart: always
    image: 'redis:alpine'

  mastodon-web:
    restart: always
    image: 'tootsuite/mastodon'
    command: 'bundle exec rails s -p 3000'
    ports: 
    - "3000:3000"

  mastodon-sidekiq:
    restart: always
    image: 'tootsuite/mastodon'
    command: 'bundle exec sidekiq'
```

### 1. MastodonのPortを変更する
spam-filterを通信経路に差し込むために、ポートを変更します。

```
  mastodon-web:
    restart: always
    image: 'tootsuite/mastodon'
    command: 'bundle exec rails s -p 3050' # 3050に変更
    # すべての通信がDockerのDefaultネットワーク経由で流れるのでここでのポート開放は基本的に不要になる
    # ports: 
    # - "3050:3050"
```

### 2. spam-filterを追加する
イメージを追加します。ポートが元のMastodonのものである `3000` と一致するようにします。

```
  spam-filter:
    image: webscrubbing808/simple-activitypub-spam-filter:v0.1.0
    ports: 
    - "3000:3000"
    environment:
    - BLOCK_WORDS="THE_EXAMPLE_SPAM_URL.org,EXAMPLE_WORDS"
    - LISTEN_ADDRESS="0.0.0.0:3000"
    - WHEN_DETECT_SPAM="block"
    - PROXY_TARGET="http://mastodon-web:3050" # serviceのキー名を指定
```

### 3. 起動する
```
> docker compose up -d
> docker compose logs -f
```

## Kubernetesで利用する場合
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: spam-filter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: spam-filter
  template:
    metadata:
      labels:
        app: spam-filter
    spec:
      containers:
        - name: spam-filter
          image: docker.io/webscrubbing808/simple-activitypub-spam-filter:v0.1.0
          env:
            - name: BLOCK_WORDS
              value: "THE_EXAMPLE_SPAM_URL.org,EXAMPLE_WORDS"
            - name: LISTEN_ADDRESS
              value: "0.0.0.0:80"
            - name: WHEN_DETECT_SPAM
              value: "block"
            - name: PROXY_TARGET
              value: "http://your_mastodon_web.namespace.svc.cluster.local"
---
apiVersion: v1
kind: Service
metadata:
  name: spam-filter
spec:
    selector:
        app: spam-filter
    ports:
      - protocol: TCP
        port: 80
        targetPort: 80
    type: ClusterIP
```

上記を適応し、Ingressで対象のサービスを差し替えるのがよいと思われます。

## Issue && Pull Request

Welcome よりよいスパムチェック方法や、ドキュメントの改善などが必要です。

## License

CC0