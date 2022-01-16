.DEFAULT_GOAL := help

VERSION=$(shell git describe --always --long)
PROJECT_NAME := go-data-collection
CLONE_URL:=github.com/daiemna/$(PROJECT_NAME)
IDENTIFIER= $(VERSION)-$(GOOS)-$(GOARCH)
BUILD_TIME=$(shell date -u +%FT%T%z)
LDFLAGS='-extldflags "-static" -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)'

help:          ## Show available options with this Makefile
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -v grep | awk 'BEGIN { FS = ":.*?##" }; { printf "%-18s  %s\n", $$1,$$2 }'

.PHONY : test
test:	 ## Run all the tests, assumes that docker-containers are up and working
	go build . && chmod +x ./scripts/test.sh && ./scripts/test.sh
	

.PHONY : test_env
recreate_test_env:  ## Create the docker-container to run integration-tests against
	docker-compose -f devsetup/docker-compose.yml down
	docker-compose -f devsetup/docker-compose.yml up -d --build
	./devsetup/wait-for-it.sh -t 40 localhost:9042 -- echo "cassandra active, will wait for 60s to bootstrap"
	sleep 60
	docker exec -it grpc-cassandra-c cqlsh -u cassandra -p cassandra -f /etc/init.cql

.PHONY : run
run:  ## runs main file.
	go run main.go

test_clean:         ## Clean the application
	@go clean -i ./...
	@rm -rf ./$(PROJECT_NAME)
	@rm -rf build/*
	@docker-compose -f devsetup/docker-compose.yml down


build: vendor
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o '$(FLAGS)' -a -ldflags $(LDFLAGS)  .

vendor:           ## Go get vendor
	go mod vendor

crossbuild:
	mkdir -p build/${PROJECT_NAME}-$(IDENTIFIER)
	make build FLAGS="build/${PROJECT_NAME}-$(IDENTIFIER)/${PROJECT_NAME}"
	cd build \
	&& tar cvzf "${PROJECT_NAME}-$(IDENTIFIER).tgz" "${PROJECT_NAME}-$(IDENTIFIER)" \
	&& rm -rf "${PROJECT_NAME}-$(IDENTIFIER)"

release:	vendor 	test_clean ## Create a release build.
	make crossbuild GOOS=linux GOARCH=amd64
	make crossbuild GOOS=linux GOARCH=386
	make crossbuild GOOS=darwin GOARCH=amd64
	make crossbuild GOOS=windows GOARCH=amd64