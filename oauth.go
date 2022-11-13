package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/oauth2"
	"io"
)

const state = "FIXME"

var config = oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://twitter.com/i/oauth2/authorize",
		TokenURL:  "https://api.twitter.com/2/oauth2/token",
		AuthStyle: oauth2.AuthStyleInHeader,
	},
	RedirectURL: "http://localhost:8080/auth",
	Scopes:      []string{"tweet.write", "tweet.read", "users.read"},
}

var codeVerifier = generateBase64Encoded32byteRandomString() // TODO: 本番ではブラウザごとのセッションに保存してください

func buildAuthorizationURL(config oauth2.Config) string {
	// PKCE 対応 https://datatracker.ietf.org/doc/html/rfc7636
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	hashed := h.Sum(nil)
	codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hashed)

	url := config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	return url
}

func generateBase64Encoded32byteRandomString() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}
