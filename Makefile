APIDOC_BASE = cmd/api
APIDOC_INFO = internal/http/handler
MYSQL_CONNECTION = mysql://$(MYSQL_URI)

include .env

migrate:
	migrate create -ext sql -dir database/migration/ -seq $(create)

migrate_up:
	migrate -path database/migration -database '$(MYSQL_CONNECTION)' -verbose up

migrate_down:
	migrate -path database/migration -database '$(MYSQL_CONNECTION)' -verbose down

migrate_rollback:
	migrate -path database/migration -database '$(MYSQL_CONNECTION)' -verbose down $(shell echo ${step}-1 | bc)

migrate_fix: 
	migrate -path database/migration -database '$(MYSQL_CONNECTION)' force $(version)

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
	mockery

jwt-keygen:
	openssl genrsa -out private_key.pem 4096
	openssl rsa -in private_key.pem -pubout -out public_key.pem

api-run:
	go run cmd/api/main.go

worker-run:
	go run cmd/worker/main.go $(topic)