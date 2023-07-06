BINARY_FILE=./bin/atomic

lint:
	golangci-lint run ./...

.PHONY: build
build: lint
	go build -o ${BINARY_FILE} ./cmd/atomic

run: build
	${BINARY_FILE}

clean:
	rm -rf ./bin

test:
	go test -v ./...

race:
	go test -race -v ./...

mock:
	go generate ./...

up:
	docker compose up -d

down:
	docker compose down
