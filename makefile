.PHONY: run
run:
	@docker-compose -f docker-compose.yml run --rm lana-cli

.PHONY: mod
mod:
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod vendor

.PHONY: test
test:
	@docker-compose -f docker-compose.test.yml run --rm lana-test
