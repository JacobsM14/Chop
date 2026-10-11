BINARY  := chop
BIN_DIR := bin
PKG     := ./...

GO      ?= go
GOFLAGS ?=
LDFLAGS ?= -s -w

.DEFAULT_GOAL := build

.PHONY: help build run test cover fmt vet lint tidy check install clean

help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: ## Compila el binario en bin/
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) .

run: ## Ejecuta el programa (pasa argumentos con ARGS="...")
	$(GO) run . $(ARGS)

test: ## Ejecuta los tests con el detector de carreras
	$(GO) test -race $(PKG)

cover: ## Genera el informe de cobertura en coverage.html
	$(GO) test -coverprofile=coverage.out $(PKG)
	$(GO) tool cover -html=coverage.out -o coverage.html

fmt: ## Formatea el código
	$(GO) fmt $(PKG)

vet: ## Análisis estático con go vet
	$(GO) vet $(PKG)

lint: ## Ejecuta golangci-lint (si está instalado)
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint no está instalado"; exit 1; }
	golangci-lint run

tidy: ## Limpia y sincroniza go.mod / go.sum
	$(GO) mod tidy

check: fmt vet test ## fmt + vet + test

install: ## Instala el binario en $GOBIN / $GOPATH/bin
	$(GO) install -ldflags '$(LDFLAGS)' .

clean: ## Borra artefactos de compilación
	rm -rf $(BIN_DIR) coverage.out coverage.html
