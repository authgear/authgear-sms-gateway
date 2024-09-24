# GIT_NAME could be empty.
GIT_NAME ?= $(shell git describe --exact-match 2>/dev/null)
GIT_HASH ?= git-$(shell git rev-parse --short=12 HEAD)

.PHONY: vendor
vendor:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2
	go mod download
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: generate
generate:
	go generate ./pkg/... ./cmd/...
	$(MAKE) fmt

.PHONY: start
start:
	go run ./cmd/server

.PHONY: build
build:
	go build -o authgear-sms-gateway -tags "osusergo netgo static_build timetzdata" ./cmd/server

.PHONY: fmt
fmt:
	find ./pkg ./cmd -name '*.go' -not -name 'wire_gen.go' -not -name '*_mock_test.go' | sort | xargs goimports -w -format-only -local github.com/authgear/authgear-sms-gateway


.PHONY: lint
lint:
	go vet ./cmd/... ./pkg/...

.PHONY: check-dockerignore
check-dockerignore:
	./scripts/sh/check-dockerignore.sh

.PHONY: govulncheck
govulncheck:
	govulncheck -show traces,version,verbose ./...

.PHONY: test
test:
	go test ./...

.PHONY: check-tidy
check-tidy:
	$(MAKE) fmt
	go mod tidy
	git status --porcelain | grep '.*'; test $$? -eq 1

.PHONY: build-image
build-image:
	docker build --pull --file ./cmd/$(TARGET)/Dockerfile --tag $(IMAGE_NAME) --build-arg GIT_HASH=$(GIT_HASH) .

.PHONY: tag-image
tag-image: DOCKER_IMAGE = quay.io/theauthgear/$(IMAGE_NAME)
tag-image:
	docker tag $(IMAGE_NAME) $(DOCKER_IMAGE):latest
	docker tag $(IMAGE_NAME) $(DOCKER_IMAGE):$(GIT_HASH)
	if [ ! -z $(GIT_NAME) ]; then docker tag $(IMAGE_NAME) $(DOCKER_IMAGE):$(GIT_NAME); fi

.PHONY: push-image
push-image: DOCKER_IMAGE = quay.io/theauthgear/$(IMAGE_NAME)
push-image:
	docker manifest inspect $(DOCKER_IMAGE):$(GIT_HASH) > /dev/null; if [ $$? -eq 0 ]; then \
		echo "$(DOCKER_IMAGE):$(GIT_HASH) exists. Skip push"; \
	else \
		docker push $(DOCKER_IMAGE):latest ;\
		docker push $(DOCKER_IMAGE):$(GIT_HASH) ;\
		if [ ! -z $(GIT_NAME) ]; then docker push $(DOCKER_IMAGE):$(GIT_NAME); fi ;\
	fi
