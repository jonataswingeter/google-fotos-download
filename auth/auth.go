package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const credentialsFile = "credentials.json"
const tokenFile = "token.json"

func GetClient(ctx context.Context, scopes ...string) *http.Client {
	config := getConfig(scopes)
	tok := getToken(ctx, config)
	return config.Client(ctx, tok)
}

func getConfig(scopes []string) *oauth2.Config {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Erro ao ler %s: %v", credentialsFile, err)
	}
	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		log.Fatalf("Erro ao parsear credenciais: %v", err)
	}
	return config
}

func getToken(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	tok, err := tokenFromFile()
	if err != nil {
		tok = getTokenFromWeb(ctx, config)
		saveToken(tok)
	}
	return tok
}

func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Abra o link no navegador e cole o código abaixo:\n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Erro ao ler código de autenticação: %v", err)
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Fatalf("Erro ao trocar código por token: %v", err)
	}
	return tok
}

func tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tok oauth2.Token
	err = json.NewDecoder(f).Decode(&tok)
	return &tok, err
}

func saveToken(token *oauth2.Token) {
	f, err := os.Create(tokenFile)
	if err != nil {
		log.Fatalf("Não foi possível salvar o token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	fmt.Printf("Token salvo em %s\n", tokenFile)
}
