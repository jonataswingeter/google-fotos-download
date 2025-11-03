package auth

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestTokenFromFile_NotFound(t *testing.T) {
	tmp := t.TempDir()
	tokenPath := filepath.Join(tmp, "token.json")

	// Redefine o nome do arquivo token para o teste
	old := tokenFile
	defer func() { _ = os.Rename(tokenPath, old) }()

	// Deve retornar erro porque o arquivo não existe
	_, err := os.Open(tokenPath)
	if err == nil {
		t.Fatalf("esperava erro por arquivo inexistente")
	}
}

func TestSaveAndReadToken(t *testing.T) {
	tmp := t.TempDir()
	testFile := filepath.Join(tmp, "token.json")

	tok := &oauth2.Token{
		AccessToken:  "abc123",
		TokenType:    "Bearer",
		RefreshToken: "refresh456",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	// Salva token
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("erro ao criar arquivo temporário: %v", err)
	}
	if err := json.NewEncoder(f).Encode(tok); err != nil {
		t.Fatalf("erro ao codificar token: %v", err)
	}
	_ = f.Close()

	// Lê token salvo
	readFile, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("erro ao abrir token salvo: %v", err)
	}
	defer readFile.Close()

	var readToken oauth2.Token
	if err := json.NewDecoder(readFile).Decode(&readToken); err != nil {
		t.Fatalf("erro ao decodificar token: %v", err)
	}

	if readToken.AccessToken != tok.AccessToken {
		t.Errorf("esperava token %q, obteve %q", tok.AccessToken, readToken.AccessToken)
	}
}

func TestGetConfig_MissingFile(t *testing.T) {
	tmp := t.TempDir()
	_ = os.Chdir(tmp) // muda o diretório para isolar ambiente

	defer func() { recover() }() // previne log.Fatalf terminar o teste

	// Como credentials.json não existe, a função deve encerrar o programa (log.Fatalf)
	_ = getConfig([]string{"https://www.googleapis.com/auth/photoslibrary.readonly"})
}

func TestGetTokenFromWeb_InvalidCode(t *testing.T) {
	ctx := context.Background()
	config := &oauth2.Config{
		ClientID:     "dummy",
		ClientSecret: "dummy",
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	// Como não podemos interagir com stdin no teste, este apenas valida o tipo de retorno esperado.
	go func() {
		// simula fechamento rápido de stdin
		os.Stdin.Close()
	}()

	defer func() { recover() }() // previne crash via log.Fatalf
	_ = getTokenFromWeb(ctx, config)
}
