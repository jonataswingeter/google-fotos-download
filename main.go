package main

import (
	"context"
	"log"

	"github.com/jonataswingeter/google-fotos-download/auth"
	"github.com/jonataswingeter/google-fotos-download/download"
)

func main() {
	ctx := context.Background()
	client := auth.GetClient(ctx, "https://www.googleapis.com/auth/photoslibrary.readonly")

	if err := download.DownloadAll(ctx, client); err != nil {
		log.Fatalf("Erro no download: %v", err)
	}
}
