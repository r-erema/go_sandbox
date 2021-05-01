GOLANGCI_IMAGE=golangci/golangci-lint:latest-alpine

.PHONY: test
test:
	docker-compose up -d && docker-compose exec \
										-e MYSQL_HOST=mysql \
										-e POSTGRES_HOST=postgres \
										-e MONGODB_HOST=mongodb \
										golang go test -race -v ./...

.PHONY: lint
lint:
	docker run -v ${PWD}:/app -w /app ${GOLANGCI_IMAGE} golangci-lint run --fix --timeout 20m --sort-results
