# ==================================================================================== #
# ENVIRONMENT VARIABLES
# ==================================================================================== #

# Change these variables as necessary
main_package_path = ./cmd/app
binary_name = app

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run all quality control checks and report all issues
.PHONY: audit
audit: test
	@echo "Starting comprehensive audit..."
	@FAIL=0; \
	echo "\n=== Module Checks ==="; \
	go mod tidy -diff || FAIL=1; \
	go mod verify || FAIL=1; \
	\
	echo "\n=== Code Analysis ==="; \
	go vet ./... || FAIL=1; \
	echo "\n=== Staticcheck ==="; \
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./... || FAIL=1; \
	\
	echo "\n=== Formatting Check ==="; \
	if [ -n "$$(gofmt -l .)" ]; then \
		echo "Files needing formatting:"; \
		gofmt -l .; \
		FAIL=1; \
	fi; \
	\
	echo "\n=== Vulnerability Check ==="; \
	go run golang.org/x/vuln/cmd/govulncheck@latest ./... || FAIL=1; \
	\
	if [ $$FAIL -eq 1 ]; then \
		echo "\nAudit failed: Some checks require attention"; \
		exit 1; \
	else \
		echo "\nAudit passed: All checks OK"; \
	fi
	
## test: run user service and handler tests
.PHONY: test
test:
	go test -race -buildvcs ./tests/unit/services ./tests/unit/handlers

.PHONY: test/verbos
test/verbos:
	go test -v -race -buildvcs ./tests/unit/services ./tests/unit/handlers
## test/cover: run tests with coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./tests/unit/services ./tests/unit/handlers
	go tool cover -html=/tmp/coverage.out

## upgradeable: list upgradable dependencies
.PHONY: upgradeable
upgradeable:
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy dependencies and format code
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application locally
.PHONY: build
build:
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

## run: run the application locally
.PHONY: run
run: build
	/tmp/bin/${binary_name}

# ==================================================================================== #
# DOCKER OPERATIONS
# ==================================================================================== #

## docker/build: build docker images
.PHONY: docker/build
docker/build:
	docker compose up --build

## docker/up: start all containers in detached mode
.PHONY: docker/up
docker/up:
	docker compose up -d

## docker/down: stop and remove containers
.PHONY: docker/down
docker/down:
	docker compose down

## docker/restart: restart containers
.PHONY: docker/restart
docker/restart: docker/down docker/up

## docker/logs: view app container logs
.PHONY: docker/logs
docker/logs:
	docker compose logs -f app

## docker/ps: show running containers
.PHONY: docker/ps
docker/ps:
	docker compose ps

## docker/clean: remove containers and volumes
.PHONY: docker/clean
docker/clean:
	docker compose down -v

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to remote repository
.PHONY: push
push: confirm audit no-dirty
	git push

## production/deploy: build production binary
.PHONY: production/deploy
production/deploy: confirm audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/tmp/bin/linux_amd64/${binary_name} ${main_package_path}
	upx -5 /tmp/bin/linux_amd64/${binary_name}

# ==================================================================================== #
# MOCKS
# ==================================================================================== #

## generate/mocks: generate mock implementations
.PHONY: generate/mocks
generate/mocks:
	go run github.com/vektra/mockery/v2@v2.32.4 --config .mockery.yml