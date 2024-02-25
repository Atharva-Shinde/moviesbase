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
audit: vendor
	@echo 'performing audit'
	go fmt ./...
	go vet ./...
	staticcheck ./...
	go test -race -vet=off ./...

.PHONY: vendor
vendor:
	@echo 'vendoring dependencies'
	@go mod tidy
	@go mod verify
	go mod vendor