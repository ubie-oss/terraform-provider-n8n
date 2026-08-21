default: test

# Pin Terraform for TF_ACC so tfenv shims do not fall back to a global
# version when tests run in a temp directory without .terraform-version.
# Write-only credential data requires Terraform 1.11+.
TF_ACC_TERRAFORM_VERSION ?= $(shell cat .terraform-version)
export TF_ACC_TERRAFORM_VERSION

.PHONY: test
test:
	unset TF_ACC && cd "internal/" && go test -count=1 -v ./...

.PHONY: testacc
testacc:
	TF_ACC=1 go test ./internal/provider/... -v $(TESTARGS) -timeout 120m

# Start pinned Community n8n (docker-compose.dev.yml), mint an API key via /rest
# (scripts/bootstrap-n8n.sh), then run TF_ACC tests against /api/v1.
# Licensed APIs (team projects, folders) skip on Community 403.
# Credential CRUD is available on Community.
# Requires Docker Compose v2, curl, and python3.
.PHONY: testacc-docker
testacc-docker:
	docker compose -f docker-compose.dev.yml up -d --wait --wait-timeout 180
	@set -e; \
	trap 'docker compose -f docker-compose.dev.yml down -v' EXIT; \
	eval "$$(./scripts/bootstrap-n8n.sh --export)"; \
	TF_ACC=1 go test ./internal/provider/... -v $(TESTARGS) -timeout 20m

.PHONY: clean
clean:
	go clean -cache -modcache -i -r

.PHONY: build
build: gen-docs go-tidy gosec deadcode
	go build -v ./

.PHONY: gosec
gosec:
	gosec ./internal/...

.PHONY: deadcode
deadcode:
	go run golang.org/x/tools/cmd/deadcode -test ./...

.PHONY: upgrade-go-mod
upgrade-go-mod:
	go get -u ./...
	go mod tidy
	go mod vendor

.PHONY: lint
lint: run-trunk-check deadcode gosec

.PHONY: run-trunk-check
run-trunk-check:
	trunk check --all -y

.PHONY: format
format: format-go format-trunk

.PHONY: format-trunk
format-trunk:
	trunk fmt --all

.PHONY: format-go
format-go:
	go fmt ./internal/...

.PHONY: install
install:
	go build -v ./ && go install .

.PHONY: gen-docs
gen-docs:
	go generate ./...

.PHONY: go-tidy
go-tidy:
	go mod tidy

.PHONY: setup-dev
setup-dev: unset-git-hooks setup-trunk

.PHONY: unset-git-hooks
unset-git-hooks:
	git config --unset-all core.hooksPath || true

.PHONY: setup-trunk
setup-trunk:
	trunk git-hooks sync

.PHONY: update
update: update-trunk

.PHONY: update-trunk
update-trunk:
	trunk upgrade
