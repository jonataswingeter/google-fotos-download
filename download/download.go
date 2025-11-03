package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jonataswingeter/google-fotos-download/sanitize"

	"google.golang.org/api/photoslibrary/v1"
)

// DownloadAll baixa todos os álbuns e suas fotos, salvando-as separadas por ano
// e garantindo nomes únicos mesmo se houver repetição de nomes.
func DownloadAll(ctx context.Context, client *http.Client) error {
	svc, err := photoslibrary.New(client)
	if err != nil {
		return fmt.Errorf("erro ao criar serviço Google Photos: %w", err)
	}

	pageToken := ""
	for {
		albumsResp, err := svc.Albums.List().PageToken(pageToken).Do()
		if err != nil {
			return fmt.Errorf("erro ao listar álbuns: %w", err)
		}

		for _, album := range albumsResp.Albums {
			if err := downloadAlbum(client, ctx, svc, album); err != nil {
				return fmt.Errorf("erro ao baixar álbum %s: %w", album.Title, err)
			}
		}

		if albumsResp.NextPageToken == "" {
			break
		}
		pageToken = albumsResp.NextPageToken
	}

	fmt.Println("✅ Download concluído.")
	return nil
}

// downloadAlbum baixa todos os itens de mídia de um álbum.
func downloadAlbum(client *http.Client, ctx context.Context, svc *photoslibrary.Service, album *photoslibrary.Album) error {
	searchReq := &photoslibrary.SearchMediaItemsRequest{
		AlbumId:  album.Id,
		PageSize: 100, // máximo permitido
	}

	pageToken := ""
	for {
		searchReq.PageToken = pageToken
		resp, err := svc.MediaItems.Search(searchReq).Do()
		if err != nil {
			return fmt.Errorf("erro ao listar fotos do álbum %s: %w", album.Title, err)
		}

		for _, item := range resp.MediaItems {
			if item.MediaMetadata == nil || item.MediaMetadata.CreationTime == "" {
				continue
			}

			date, err := time.Parse(time.RFC3339, item.MediaMetadata.CreationTime)
			if err != nil {
				date = time.Now()
			}

			year := sanitize.ExtractYear(date)
			ext := ".jpg" // default
			if item.MimeType == "image/png" {
				ext = ".png"
			} else if item.MimeType == "image/gif" {
				ext = ".gif"
			}

			// Gerar nome único usando o ID do item
			filename := sanitize.UniqueFileName(item.Id + ext)
			dir := filepath.Join("photos", year)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
			}

			path := filepath.Join(dir, filename)
			url := fmt.Sprintf("%s=d", item.BaseUrl)

			if err := downloadFile(client, url, path); err != nil {
				return fmt.Errorf("erro ao baixar %s: %w", filename, err)
			}

			fmt.Printf("📸 %s → %s\n", filename, path)
		}

		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return nil
}

// downloadFile baixa a URL de uma foto e salva localmente.
func downloadFile(client *http.Client, url, path string) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("erro no GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resposta inválida (%d) para %s", resp.StatusCode, url)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo %s: %w", path, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("erro ao salvar %s: %w", path, err)
	}

	return nil
}
