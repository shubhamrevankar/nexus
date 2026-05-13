# Nexus

Nexus is an AI-native operational intelligence platform for engineering teams and startups.

The repository is initialized as a production-grade engineering foundation. It intentionally avoids product features until the platform basics are reliable: monorepo tooling, local development, CI, Docker, observability conventions, and documentation workflows.

## Workspace

```text
apps/
  web/               Next.js workspace shell
services/
  api/               Go modular-monolith API foundation
  ai/                Python AI systems foundation
packages/
  config/            Shared environment/config helpers
  logger/            Shared TypeScript logging helper
  types/             Shared TypeScript contracts
infra/
  docker/            Local Docker assets
docs/                Engineering documentation
architecture/        Architecture overviews and diagrams
adr/                 Architecture Decision Records
rfc/                 Request for Comments documents
observability/       Observability configuration and notes
testing/             Testing strategy and shared fixtures
scripts/             Developer automation scripts
```

## Prerequisites

- Node.js 20+
- pnpm 9+
- Go 1.22+
- Python 3.12+
- Docker and Docker Compose

## Local Development

```bash
corepack enable
pnpm install
pnpm dev
```

If Corepack cannot write package-manager shims on your machine, install `pnpm` directly or run commands through `corepack pnpm`.

Useful commands:

```bash
pnpm lint
pnpm typecheck
pnpm test
pnpm build
docker compose up --build
```

The Go and Python services can also be checked without installing Node dependencies:

```bash
cd services/api && go test ./...
cd services/ai && python -m unittest discover -s tests
```

## Engineering Standards

- Start as a modular monolith with explicit internal boundaries.
- Prefer simple, testable components over premature distributed systems.
- Record meaningful technical decisions as ADRs.
- Propose large or cross-cutting changes through RFCs.
- Treat logs, metrics, tracing, request IDs, and health checks as platform requirements.

See `docs/engineering-standards.md`, `architecture/overview.md`, `adr/0001-modular-monolith.md`, and `rfc/0000-template.md`.
