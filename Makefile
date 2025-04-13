include ./configs/local.env

LOCAL_BIN:=$(CURDIR)/infra/bin
LOCAL_MIGRATION_DIR=${MIGRATION_DIR}
LOCAL_MIGRATION_DSN="host=localhost port=${PG_PORT} dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD}"

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config ./.golangci.yml

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.

generate_user_api:
	mkdir -p pkg/user_v1
	protoc --proto_path proto/user_v1 --proto_path vendor.protogen \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/user_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	--grpc-gateway_out=pkg/user_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	proto/user_v1/user.proto

generate_auth_api:
	mkdir -p pkg/auth_v1
	protoc --proto_path proto/auth_v1 --proto_path vendor.protogen \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	proto/auth_v1/auth.proto

generate_access_api:
	mkdir -p pkg/access_v1
	protoc --proto_path proto/access_v1/ --proto_path vendor.protogen \
	--go_out=pkg/access_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/access_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	proto/access_v1/access.proto

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi


local-migration-status:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_DSN} up -v
	
local-migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_DSN} down -v

.PHONY: test
test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/bovinxx/auth-service/internal/services/user/user -count 5

.PHONY: test-coverage
test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/bovinxx/auth-service/internal/service/user/...,github.com/bovinxx/auth-service/internal/api/user... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

grpc-load-test:
	ghz \
		--proto proto/auth_v1/auth.proto \
		--call auth.AuthService.Login \
		--data '{ "username": "bovinxx", "password": "bovinxx" }' \
		--rps 50 \
		--total 50 \
		--insecure \
		localhost:50051