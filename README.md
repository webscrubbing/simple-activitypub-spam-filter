# simple-activitypub-spam-filter

[webscrubbing808/simple-activitypub-spam-filter - Docker Image | Docker Hub](https://hub.docker.com/r/webscrubbing808/simple-activitypub-spam-filter)

[🗾日本語版READMEはこちら](https://github.com/webscrubbing/simple-activitypub-spam-filter/blob/main/README-JA.md)

- A spam filter designed to address bot accounts using ActivityPub, such as those on Mastodon and Misskey.
- Operates as a reverse proxy, activating when content contains certain strings.
- Particularly effective against spam that mechanically posts specific URLs.

## Environment Values
All configurations are managed via environment variables.

**BLOCK_WORDS**
```
BLOCK_WORDS=THE_EXAMPLE_SPAM_URL.org,EXAMPLE_WORDS
```

- Specify words you want to block, separated by commas.
- Searches the entire content of ActivityPub and blocks if there's an exact match.

**LISTEN_ADDRESS**
```
LISTEN_ADDRESS=:80
LISTEN_ADDRESS=0.0.0.0:8080
```

- Designates the address to listen for incoming connections from other servers.
- Can be specified by port or address.

**WHEN_DETECT_SPAM**
```
WHEN_DETECT_SPAM=output
WHEN_DETECT_SPAM=block
```

- Defines the action to take upon detecting spam.
- `output`: If spam is detected, the content is output to standard output and the delivery continues.
- `block`: If spam is detected, sends a 400 to the originating server and stops the delivery.

**PROXY_TARGET**
```
PROXY_TARGET=http://localhost:3000
PROXY_TARGET=http://mastodon:8080
PROXY_TARGET=http://your-mastodon-apache.mastodon.svc.cluster.local
```

- Specifies the server address where the delivery will be made.
- Refer to the usage instructions below for detailed configuration methods.

## Usage
The spam filter functions as a reverse proxy. Please insert the image between the network entrance and your Mastodon or Misskey server.

## Using with Docker Compose

Here is a sample for operating with the following configuration:

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

### 1. Change the Mastodon Port
To insert the spam-filter into the communication path, change the port.

```
  mastodon-web:
    restart: always
    image: 'tootsuite/mastodon'
    command: 'bundle exec rails s -p 3050' # Change to 3050
    # All POST flows through Docker's Default network, so opening ports here is basically unnecessary.
    # ports: 
    # - "3050:3050"
```

### 2. Add the spam-filter
Add the image. Ensure the port matches the original Mastodon port, `3000`.

```
  spam-filter:
    image: webscrubbing808/simple-activitypub-spam-filter
    ports: 
    - "3000:3000"
    environment:
    - BLOCK_WORDS="THE_EXAMPLE_SPAM_URL.org,EXAMPLE_WORDS"
    - LISTEN_ADDRESS="0.0.0.0:3000"
    - WHEN_DETECT_SPAM="block"
    - PROXY_TARGET="http://mastodon-web:3050" # set service key name
```

### 3. Launch

```
> docker compose up -d
> docker compose logs -f
```

## Using with Kubernetes
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

Applying the above and swapping the target service with Ingress is recommended.

## Issue && Pull Request

Welcome. Better spam check methods and documentation improvements are needed.

## License

CC0
