#confirmation
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# run the application
run:
	@go run ./cmd/api

# connect to the postgresql database
db:
	psql ${moviesbase_dsn}

# migrations up
up: confirm
	@echo 'running migrations up'
	@migrate -path ./migrations -database ${moviesbase_dsn} up