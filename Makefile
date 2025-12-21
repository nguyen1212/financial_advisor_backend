.PHONY: build

SRC_PATH:= ${PWD}

prepare:
	@go install github.com/rubenv/sql-migrate/...@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

mod:
	@go mod tidy

gen:
	## Go generate codes
	@go generate ./...
	## Swagger generate
	@swag init -g app/external/framework/gin/route.go -o app/interface/api/docs --exclude migration

up:
	@go run ${SRC_PATH}/cmd/srv/...

migrate-create:
	@$(eval NAME := $(shell read -p "Enter new file name: " v && echo $$v))
	$(eval CMD:= $*)
	cd migration;\
	sql-migrate new ${NAME}

migrate-%:
	$(eval CMD:= $*)
	cd migration;\
		sql-migrate $(CMD) -env=${ENV} -config=dbconfig.yml

## should include lint and sec later
check: test

test:
	@./scripts/test.sh

up-docker:
	@docker-compose up -d
