.PHONY: build clean test package serve update-vendor api
VERSION := $(shell git describe --always |sed -e "s/^v//")

build:
	@echo "Compiling source"
	@mkdir -p build
	go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o build/ftp-server cmd/ftp-server/main.go

clean:
	@echo "Cleaning up workspace"
	@rm -rf build
	@rm -rf dist

test:
	@echo "Running tests"
	@rm -f coverage.out
	@golint ./...
	@go vet ./...
	@go test -p 1 -v -cover ./... -coverprofile coverage.out

dist:
	goreleaser
	mkdir -p dist/upload/tar
	mkdir -p dist/upload/deb
	mkdir -p dist/upload/rpm
	mv dist/*.tar.gz dist/upload/tar
	mv dist/*.deb dist/upload/deb
	mv dist/*.rpm dist/upload/rpm

snapshot:
	@goreleaser --snapshot

dev-requirements:
	go install golang.org/x/lint/golint
	go install golang.org/x/tools/cmd/stringer
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/goreleaser/goreleaser
	go install github.com/goreleaser/nfpm

# shortcuts for development

serve: build
	@echo "Starting Ftp Server"
	./build/ftp-server

run-compose-test:
	docker-compose run --rm networkserver make test
