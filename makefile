LOCAL_BIN:=$(CURDIR)/bin
ENV_FILE:=$(CURDIR)/.env
MIGRATION_DIR:=$(CURDIR)/migration
GOOSE_BIN:=$(LOCAL_BIN)/goose

install-goose:
	GOBIN=$(LOCAL_BIN) go install -pkgdir bin github.com/pressly/goose/v3/cmd/goose@latest

make-migration:
	$(GOOSE_BIN) -dir $(MIGRATION_DIR) create "" sql

load-env:
	export $(cat $(ENV_FILE) | xargs)

migrate: load-env
	$(GOOSE_BIN) up

start:
	go run ./cmd/main.go