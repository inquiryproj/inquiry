golangciLintVersion = "v1.54.2"
gofumptVersion = "v0.5.0"
gciVersion = "v0.11.0"
oapiCodegenVersion = "v1.15.0"
mockeryVersion = "v2.19.0"

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

# Run unit tests and generate coverage report
test:
	@mkdir -p reports
	@go test -coverprofile=reports/codecoverage_all.cov ./... -cover -race -p=4
	@go tool cover -func=reports/codecoverage_all.cov > reports/functioncoverage.out
	@go tool cover -html=reports/codecoverage_all.cov -o reports/coverage.html
	@echo "View report at $(PWD)/reports/coverage.html"
	@tail -n 1 reports/functioncoverage.out

$(GOBIN)/oapi-codegen:
	@go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@${oapiCodegenVersion}

# Generate Go server interface, domain models & client from OpenAPI Spec
oapi-gen: | $(GOBIN)/oapi-codegen
	# @oapi-codegen -generate spec -package http api/api-spec.yml  > internal/app/http/spec.go
	@oapi-codegen -generate types -package http api/api-spec.yml  > internal/http/domain.go
	@oapi-codegen -generate server -package http api/api-spec.yml  > internal/http/server.go
	@oapi-codegen -generate client,types -package api api/api-spec.yml  > pkg/api/client.go

/usr/local/bin/spectral:
	@sudo npm install -g @stoplight/spectral

# Lint OpenAPI Spec
lint-oapi: | /usr/local/bin/spectral
	@spectral lint --fail-severity=warn --ruleset rule.spectral.yaml --verbose ./api/api-spec.yml 

# Mock generation
$(GOBIN)/mockery:
	@go install github.com/vektra/mockery/v2@${mockeryVersion}

mocks:
	@go generate --tags=mocks ./...
