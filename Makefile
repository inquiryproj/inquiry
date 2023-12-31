golangciLintVersion = "v1.55.0"
gofumptVersion = "v0.5.0"
gciVersion = "v0.11.0"
oapiCodegenVersion = "v1.15.0"
mockeryVersion = "v2.36.0"
dockerRepo = "wimsp/inquiry"

$(GOBIN)/golangci-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@${golangciLintVersion}

# Lint Go Code
lint: | $(GOBIN)/golangci-lint
	@echo Linting...
	@golangci-lint  -v --concurrency=3 --config=.golangci.yml --issues-exit-code=1 run \
	--out-format=colored-line-number 

$(GOBIN)/gofumpt:
	@go install mvdan.cc/gofumpt@${gofumptVersion}
	@go mod tidy

# Format Go Code
gofumpt: | $(GOBIN)/gofumpt
	@gofumpt -w $(shell ls  -d $(PWD)/*/)

$(GOBIN)/gci:
	@go install github.com/daixiang0/gci@${gciVersion}
	@go mod tidy

# Format imports
gci: | $(GOBIN)/gci
	@gci write --section Standard --section Default --section "Prefix(github.com/inquiryproj/inquiry)" $(shell ls  -d $(PWD)/*)

format: gci gofumpt

coverage-report:
	@open ./reports/coverage.html
	
# Run unit tests and generate coverage report
test:
	@mkdir -p reports
	@go test -coverprofile=reports/codecoverage_all.cov ./... -cover -race -p=4
	@go tool cover -func=reports/codecoverage_all.cov > reports/functioncoverage.out
	@go tool cover -html=reports/codecoverage_all.cov -o reports/coverage.html
	@echo "View report at $(PWD)/reports/coverage.html"
	@tail -n 1 reports/functioncoverage.out

# Run unit tests + integration tests and generate coverage report
test-integration:
	@mkdir -p reports
	@go test -coverprofile=reports/codecoverage_all.cov --tags=integration ./... -cover -race -p=4
	@go tool cover -func=reports/codecoverage_all.cov > reports/functioncoverage.out
	@go tool cover -html=reports/codecoverage_all.cov -o reports/coverage.html
	@echo "View report at $(PWD)/reports/coverage.html"
	@tail -n 1 reports/functioncoverage.out

$(GOBIN)/oapi-codegen:
	@go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@${oapiCodegenVersion}

# Generate Go server interface, domain models & client from OpenAPI Spec
oapi-gen: | $(GOBIN)/oapi-codegen
	# @oapi-codegen -generate spec -package http api/api-spec.yml  > internal/app/http/spec.go
	@oapi-codegen -generate types -package api api/api-spec.yml  > internal/http/api/domain.go
	@oapi-codegen -generate server -package api api/api-spec.yml  > internal/http/api/server.go
	@oapi-codegen -generate client,types -package api api/api-spec.yml  > pkg/api/client.go

/usr/local/bin/spectral:
	@sudo npm install -g @stoplight/spectral

# Lint OpenAPI Spec
lint-oapi: | /usr/local/bin/spectral
	@spectral lint --fail-severity=warn --ruleset rule.spectral.yaml --verbose ./api/api-spec.yml 

# Mock generation
$(GOBIN)/mockery:
	@go install github.com/vektra/mockery/v2@${mockeryVersion}

mocks: | $(GOBIN)/mockery
	@go generate --tags=mocks ./...
	@${MAKE} format

docker-build-and-push:
	@docker build --push --tag ${dockerRepo}:${VERSION} --tag wimsp/inquiry:latest --platform=linux/amd64 .

docker-build-and-push-arm:
	@docker buildx build --push --tag ${dockerRepo}:${VERSION} --tag wimsp/inquiry:latest --platform=linux/amd64 .

