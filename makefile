IMG_NAME=toanpham123/bookmark_service
GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

IMG_TAG := $(BRANCH)
ifeq ($(BRANCH), main)
	IMG_TAG := dev
endif

ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

export IMG_TAG

.PHONY: docs
docs:
	swag init -g cmd/api/main.go 

.PHONY: run
run:
	go run ./cmd/api

.PHONY: mocks
mocks:
	go generate ./..

COVERAGE_EXCLUDE=mocks|main.go|test
COVERAGE_THRESHOLD = 80

.PHONY: test
test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

COVERAGE_FOLDER=./coverage

docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --target test -t bookmark_service:dev --output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
    if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	   echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	   exit 1; \
    else \
	   echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
   	fi	

docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

docker-release: docker-build
	docker push $(IMG_NAME):$(IMG_TAG)

DOCKER_PASSWORD?=
DOCKER_USERNAME?=
docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	
.PHONY: build

PRIVATE_KEY ?= ./private.pem
PUBLIC_KEY ?= ./public.pem

generate-rsa-key:
	openssl genpkey -algorithm RSA -out $(PRIVATE_KEY) -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in $(PRIVATE_KEY) -out $(PUBLIC_KEY)
