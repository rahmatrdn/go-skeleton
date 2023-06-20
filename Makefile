APIDOC_BASE		= cmd/api
APIDOC_INFO		= internal/http/handler

migrate:
	migrate create -ext sql -dir database/migration/ -seq $(create)

migrate_up:
	migrate -path database/migration -database 'mysql://mms-dev:Spe@2022#@(10.206.11.138:3306)/bni_mms_master' -verbose up

migrate_down:
	migrate -path database/migration -database 'mysql://mms-dev:Spe@2022#@(10.206.11.138:3306)/bni_mms_master' -verbose down

migrate_rollback:
	migrate -path database/migration -database 'mysql://mms-dev:Spe@2022#@(10.206.11.138:3306)/bni_mms_master' -verbose down $(shell echo ${step}-1 | bc)

migrate_fix: 
	migrate -path database/migration -database 'mysql://mms-dev:Spe@2022#@(10.206.11.138:3306)/bni_mms_master' force $(version)

test:
	go test -cover -coverprofile=coverage.out $$(go list ./...)

apidoc:
	swag init -d $(APIDOC_BASE),$(APIDOC_INFO) --parseInternal --pd
