APIDOC_BASE = cmd/api
APIDOC_INFO = internal/http/handler
MYSQL_URI   = mysql://$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE_NAME)

include .env

migrate:
	migrate create -ext sql -dir database/migration/ -seq $(create)

migrate_up:
	migrate -path database/migration -database '$(MYSQL_URI)' -verbose up

migrate_down:
	migrate -path database/migration -database '$(MYSQL_URI)' -verbose down

migrate_rollback:
	migrate -path database/migration -database '$(MYSQL_URI)' -verbose down $(shell echo ${step}-1 | bc)

migrate_fix: 
	migrate -path database/migration -database '$(MYSQL_URI)' force $(version)

test:
	go test -cover -coverprofile=coverage.out $$(go list ./...)

apidoc:
	swag init -d $(APIDOC_BASE),$(APIDOC_INFO) --parseInternal --pd

protob:
	protoc --go_out=proto/pb --go_opt=paths=source_relative --go-grpc_out=proto/pb --go-grpc_opt=paths=source_relative proto/*.proto

coverage:
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out

mock:
	 mockery --name=$(d) --recursive=true --output=tests/mocks