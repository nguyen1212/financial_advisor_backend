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
	@swag init -g app/external/framework/gin/route.go -o app/delivery/rest/docs --exclude migration

up:
	@go run ${SRC_PATH}/cmd/srv/...

build:
	# jsoniter will be used to switch the json library instead of default one.
	# this is for performance consideration.
	@go build -tags=jsoniter -o financial_advisor_srv ${SRC_PATH}/cmd/srv/...

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
