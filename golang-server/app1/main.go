package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientID     = "app"
	clientSecret = "..."
)

func main() {

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/realms/demo")
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: "yIlt1q9rYG52Mx4cOK1aQuViWgfFzokP",
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8081/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "exemplo"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "State does not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Problema ao trocar token", http.StatusInternalServerError)
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "Problema ao pegar id token", http.StatusInternalServerError)
			return
		}

		res := struct {
			OAuth2Token *oauth2.Token
			IDToken     string
		}{
			oauth2Token, rawIDToken,
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		w.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
