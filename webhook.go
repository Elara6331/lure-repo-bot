package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"go.arsenm.dev/lure-repo-bot/internal/queue"
	"go.arsenm.dev/lure-repo-bot/internal/types"
)

type prQueue = *queue.Queue[*types.PullRequestPayload]

func serveWebhook(ctx context.Context, addr string, jobQueue prQueue) {
	mux := http.NewServeMux()

	mux.HandleFunc("/webhook", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if req.Header.Get("X-GitHub-Event") != "pull_request" {
			http.Error(res, "Only pull_request events are accepted by this bot", http.StatusBadRequest)
			return
		}

		payload, err := secureDecode(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		jobQueue.Add(payload)
	})

	srv := http.Server{
		Addr:    addr,
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	srv.ListenAndServe()
}

func secureDecode(req *http.Request) (*types.PullRequestPayload, error) {
	sigStr := req.Header.Get("X-Hub-Signature-256")
	sig, err := hex.DecodeString(strings.TrimPrefix(sigStr, "sha256="))
	if err != nil {
		return nil, err
	}

	secretStr, ok := os.LookupEnv("LURE_BOT_SECRET")
	if !ok {
		return nil, errors.New("LURE_BOT_SECRET must be set to the secret used for setting up the github webhook")
	}
	secret := []byte(secretStr)

	h := hmac.New(sha256.New, secret)
	r := io.TeeReader(req.Body, h)

	payload := &types.PullRequestPayload{}
	err = json.NewDecoder(r).Decode(payload)
	if err != nil {
		return nil, err
	}

	if !hmac.Equal(h.Sum(nil), sig) {
		return nil, errors.New("webhook signature mismatch")
	}

	return payload, nil
}
