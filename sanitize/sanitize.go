package sanitize

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// UniqueFileName gera um nome de arquivo único com base no nome original
// e um timestamp de alta resolução (nanosegundos).
// Exemplo: "foto.jpg" → "foto_1732200000000000000.jpg"
func UniqueFileName(original string) string {
	ext := filepath.Ext(original)
	name := strings.TrimSuffix(original, ext)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

// ExtractYear retorna o ano (YYYY) extraído de um objeto time.Time.
// Exemplo: 2025-10-22 → "2025"
func ExtractYear(date time.Time) string {
	return fmt.Sprintf("%d", date.Year())
}
