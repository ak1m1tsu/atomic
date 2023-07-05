BINARY_FILE=./bin/ushrt

.PHONY: build
build:
	go build -o ${BINARY_FILE} ./cmd/ushrt

run: build
	${BINARY_FILE}
