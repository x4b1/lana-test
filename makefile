.PHONY: run
run:
	@docker-compose -f docker-compose.yml up --build --force-recreate

.PHONY: mod
mod:
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod vendor

.PHONY: test
test:
	@docker-compose -f docker-compose.test.yml run lana-test
