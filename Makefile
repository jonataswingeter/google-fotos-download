# Makefile para o projeto google-fotos-download

APP_NAME = google-fotos-download

# Targets principais
.PHONY: all run test tidy clean

all: tidy build

build:
	@echo "🛠️  Compilando o aplicativo..."
	go build -o $(APP_NAME) ./...

run:
	@echo "🚀 Executando o aplicativo..."
	go run .

test:
	@echo "🧪 Executando testes..."
	go test ./... -v

tidy:
	@echo "🔄 Atualizando dependências..."
	go mod tidy

clean:
	@echo "🧹 Limpando arquivos gerados..."
	rm -f $(APP_NAME)
	rm -rf photos
	rm -f token.json
