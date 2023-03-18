build-dir:
	mkdir -p .build

build: build-dir
	CGO_ENABLED=0 go build -o .build/squash github.com/jlewi/squash/cmd

tidy:
	gofmt -s -w .
	goimports -w .

lint:
	# golangci-lint automatically searches up the root tree for configuration files.
	golangci-lint run

test:
	go test -v ./...