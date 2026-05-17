dockerfile=test/integration/docker-compose.yaml

up-container:
	@docker compose -f "${dockerfile}" up -d
	@sleep 5

down-container:
	@docker compose -f "${dockerfile}" down

test-integrations:
	@go test -tags=integration -v ./test/integration

tests-run: up-container test-integrations down-container 