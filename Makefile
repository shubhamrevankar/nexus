.PHONY: bootstrap build dev format lint test typecheck docker-up docker-down api-test ai-test

bootstrap:
	pnpm install

dev:
	pnpm dev

build:
	pnpm build

lint:
	pnpm lint

typecheck:
	pnpm typecheck

test: api-test ai-test
	pnpm test

format:
	pnpm format

docker-up:
	docker compose up --build

docker-down:
	docker compose down --remove-orphans

api-test:
	cd services/api && go test ./...

ai-test:
	cd services/ai && python -m unittest discover -s tests

