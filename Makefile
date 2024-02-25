include .envrc
#confirmation
.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# run the application
.PHONY: run/api
run/api:
	go run ./cmd/api -dsn=${moviesbase_dsn}

# connect to the postgresql database
.PHONY: db
db:
	psql ${moviesbase_dsn}

migrations/new: confirm
	@echo 'creating new migration files named: ${filename}'
	@migrate create -seq -ext=.sql -dir=./migrations ${filename}

# migrations up
.PHONY: migrations/up
migrations/up: confirm
	@echo 'running migrations up'
	@migrate -path ./migrations -database ${moviesbase_dsn} up

.PHONY: audit
audit:
	@echo 'performing audit'
	go mod tidy
	go mod verify
	go fmt ./...
	go vet ./...
	staticcheck ./...
	go test -race -vet=off ./...
