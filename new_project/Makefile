GOLANG_IMAGE=golang:1.15.3
GOLANG_CI_LINT_IMAGE=golangci/golangci-lint:latest-alpine
GOLANG_MIGRATE_IMAGE=migrate/migrate:v4.13.0x

.PHONY: build
build:
	docker run \
		--rm \
		-e GOOS=linux \
		-e CGO_ENABLED=0 \
		-e GOARCH=amd64 \
		-v ${PWD}:/app \
		-w /app \
		$(GOLANG_IMAGE) go build -ldflags="-s -w" -o deployment/aws/build/bin/users cmd/main.go

.PHONY: create-migration
create-migration:
	docker run \
		-v ${PWD}/internal/migration:/migration \
		${GOLANG_MIGRATE_IMAGE} \
			create \
			-ext sql \
			-dir /migration \
			-seq ${name}

.PHONY: migration-up
migration-up:
	docker run \
			--env DB_DSN=${DB_DSN} \
			-v ${PWD}:/app \
			-w /app \
			$(GOLANG_IMAGE) go run ./cmd/migration/main.go

.PHONY: proto-gen
proto-gen:
	protoc  --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			pkg/infrastructure/grpc/server.proto

.PHONY: mock-gen
mock-gen:
	docker run \
		--rm \
		-v ${PWD}:/app \
		-w /app $(GOLANG_IMAGE) \
			go get ./...
			mockgen -destination=./test/mock/request_mock.go -package=mock \
					github.com/twizar/users/internal/infrastructure/request Request
			mockgen -destination=./test/mock/logger_mock.go -package=mock \
					github.com/twizar/users/internal/infrastructure/log Logger
			mockgen -destination=./test/mock/hashing_mock.go -package=mock \
					github.com/twizar/users/internal/infrastructure/hashing Hashing
			mockgen -destination=./test/mock/repository_mock.go -package=mock -mock_names=User=MockUserRepo \
					github.com/twizar/users/internal/infrastructure/repository/user User
			mockgen -destination=./test/mock/uuid_mock.go -package=mock -mock_names=Generator=MockUUIDGenerator \
					github.com/twizar/users/internal/infrastructure/uuid Generator
			mockgen -destination=./test/mock/user_service_mock.go -package=mock \
					github.com/twizar/users/internal/application/service UserService
			mockgen -destination=./test/mock/email_queue_mock.go -package=mock -mock_names=Deliverer=MockDeliverer \
					github.com/twizar/users/internal/infrastructure/message Deliverer
			mockgen -destination=./test/mock/user_buffer_mock.go -package=mock -mock_names=User=MockUserBuffer \
					github.com/twizar/users/internal/infrastructure/sign_up_buffer User

.PHONY: go-lint
go-lint:
	docker run -v ${PWD}:/app -w /app $(GOLANG_CI_LINT_IMAGE) golangci-lint run -v --fix --timeout 20m --sort-results

.PHONY: go-test
go-test:
	docker run \
		--env-file=test.env \
		-v ${PWD}:/app \
		-w /app $(GOLANG_IMAGE) \
		go test -race -cover -v -coverpkg=./... -coverprofile=cover.out ./...
		go tool cover -html=cover.out

.PHONY: deploy
deploy: go-lint go-test build

.PHONY: build-dev-debug
build-dev-debug:
	GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o ./deployment/dev/build/users ./cmd/main.go

.PHONY: aws-sam-dev-debug-api
aws-sam-dev-debug-api: build-dev-debug
	sam local start-api \
		-d 2345 \
		-p ${PORT} \
		--template ./deployment/dev/template.yaml \
		--debugger-path ./deployment/dev/linux_debugger \
		--debug-args="-delveAPI=2" \
		--env-vars ./deployment/dev/test.json

.PHONY: aws-sam-dev-api
aws-sam-dev-api: build-dev-debug
	sam local start-api \
		-p ${PORT} \
		--template ./deployment/dev/template.yaml \
		--env-vars ./deployment/dev/test.json

.PHONY: dev-redis
dev-redis:
	docker run --rm -d -p ${PORT}:6379 $(REDIS_IMAGE)