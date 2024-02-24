package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Activity struct {
	Actor  string `json:"actor"`
	Object struct {
		Content string `json:"content"`
	} `json:"object"`
}

func IsSpam(content string) bool {
	blockWordEnv := os.Getenv("BLOCK_WORDS")
	blockWordList := strings.Split(blockWordEnv, ",")

	if blockWordEnv == "" || len(blockWordList) == 0 {
		fmt.Println("Environment variable BLOCK_WORDS is not set. Please set it to block words. ex: \"example_spam_word_1,example_spam_word_2\"")
		os.Exit(1)
	}

	for _, block := range blockWordList {
		if strings.Contains(content, block) {
			return true
		}
	}
	return false
}

func main() {
	target := os.Getenv("PROXY_TARGET")
	if target == "" {
		fmt.Println("Environment variable PROXY_TARGET is not set. Please set it to the target URL. ex: \"http://mastodon:8080\"")
		os.Exit(1)
	}
	proxyUrl, _ := url.Parse(target)

	whenDetectSpam := os.Getenv("WHEN_DETECT_SPAM")
	if whenDetectSpam == "" || (whenDetectSpam != "block" && whenDetectSpam != "output") {
		fmt.Println("Environment variable WHEN_DETECT_SPAM is not set. Please set it to block or report")
		os.Exit(1)
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.Director = func(req *http.Request) {
		originalPath := req.URL.Path
		req.URL.Scheme = proxyUrl.Scheme
		req.URL.Host = proxyUrl.Host
		req.URL.Path = originalPath
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Body == nil {
			proxy.ServeHTTP(w, req)
			return
		}

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println("Error reading body")
			os.Exit(1)
		}

		if req.Body.Close() != nil {
			fmt.Println("Error closing body")
			os.Exit(1)
		}

		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var activity Activity
		if err := json.Unmarshal(bodyBytes, &activity); err != nil {
			proxy.ServeHTTP(w, req)
			return
		}

		if IsSpam(activity.Object.Content) {
			fmt.Println("Spam detected: ", activity.Object.Content)
			if whenDetectSpam == "block" {
				http.Error(w, "Spam detected", http.StatusForbidden)
			} else {
				proxy.ServeHTTP(w, req)
			}
			return
		}

		// Not spam
		proxy.ServeHTTP(w, req)
		return
	})

	addr := os.Getenv("LISTEN_ADDRESS")
	if addr == "" {
		fmt.Println("Environment variable LISTEN_ADDRESS is not set. Using default 0.0.0.0:80")
		addr = "0.0.0.0:80"
	}

	log.Println("Starting reverse proxy server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
