GIT_HASH ?= git-$(shell git rev-parse --short=12 HEAD)
IMAGE ?= quay.io/theauthgear/authgear-sms-gateway:$(GIT_HASH)

.PHONY: vendor
vendor:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2
	go mod download
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

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
	docker build --pull --file ./cmd/server/Dockerfile --tag $(IMAGE) .

.PHONY: push-image
push-image:
	docker push $(IMAGE)
