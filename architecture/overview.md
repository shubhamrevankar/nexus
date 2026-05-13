# Architecture Overview

Nexus begins as a modular monolith inside a monorepo.

## Current Boundaries

- `apps/web` contains the user-facing workspace shell.
- `services/api` contains the core backend API foundation.
- `services/ai` contains AI system foundations for future RAG, agents, and model orchestration.
- `packages/*` contains shared TypeScript contracts and utilities used by workspace apps.
- `infra/*` contains local infrastructure and future deployment assets.

## Direction

The platform should evolve through vertical slices. Product modules such as authentication, organizations, workspace, integrations, knowledge search, incident management, and workflows should expose explicit contracts and remain independently testable.

Microservices, Kubernetes, Terraform, vector databases, OpenSearch, ClickHouse, and graph databases are future options, not day-one defaults.

