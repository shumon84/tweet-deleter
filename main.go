package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

var (
	bearerToken  = os.Getenv("TWITTER_BEARER_TOKEN")
	clientSecret = os.Getenv("TWITTER_CLIENT_SECRET")
	clientID     = os.Getenv("TWITTER_CLIENT_ID")
)

var t *TweetDeleter

func main() {
	log.Default().SetFlags(log.LstdFlags | log.Lshortfile)
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", func(resp http.ResponseWriter, req *http.Request) {
		log.Print("[Info] /auth")
		queryCode := req.URL.Query().Get("code")
		if queryCode == "" {
			log.Println("code not found")
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		queryState := req.URL.Query().Get("state")
		if queryState == "" {
			log.Println("state not found")
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		if queryState != state {
			log.Println("invalid state")
			resp.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(context.Background(), queryCode,
			oauth2.SetAuthURLParam("code_verifier", codeVerifier))
		if err != nil {
			log.Printf("failed to exchange token: %v\n", err)
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("token scope: %v\n", token.Extra("scope"))

		oAuthClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
		tweetDeleter, err := NewTweetDeleter("list.txt", oAuthClient)
		if err != nil {
			log.Println("[Error]", err)
			resp.WriteHeader(http.StatusInternalServerError)
		}
		t = tweetDeleter
	})
	mux.HandleFunc("/login", func(resp http.ResponseWriter, req *http.Request) {
		log.Print("[Info] /login")
		url := buildAuthorizationURL(config)
		log.Println(url)
		resp.Header().Set("Location", url)
		resp.WriteHeader(http.StatusFound)
	})
	mux.HandleFunc("/delete", func(resp http.ResponseWriter, req *http.Request) {
		log.Print("[Info] /delete")
		t.delete(context.Background())
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("start tweet delete"))
	})
	mux.HandleFunc("/progress", func(resp http.ResponseWriter, req *http.Request) {
		log.Print("[Info] /progress")
		num, denom := t.getProgress()
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, "curret progress is %d / %d", num, denom)
	})

	s := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	log.Print("[Info] Server start !")
	log.Fatal(s.ListenAndServe())
}
