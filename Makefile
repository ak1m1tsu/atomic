BINARY_FILE=./bin/ushrt

lint:
	golangci-lint run ./...

.PHONY: build
build: lint
	go build -o ${BINARY_FILE} ./cmd/ushrt

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
