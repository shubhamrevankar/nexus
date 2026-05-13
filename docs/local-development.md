# Local Development

## Setup

```bash
corepack enable
pnpm install
```

Copy `.env.example` to `.env` for local overrides.

If Corepack cannot write shims because Node.js is installed in a protected directory, install `pnpm` directly or run package commands through `corepack pnpm`.

## Run Everything

```bash
docker compose up --build
```

## Run Workspace Tasks

```bash
pnpm dev
pnpm lint
pnpm typecheck
pnpm test
pnpm build
```

## Service-Specific Checks

```bash
cd services/api && go test ./...
cd services/ai && python -m unittest discover -s tests
```
