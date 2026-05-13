# Testing Strategy

Nexus should use layered testing with fast feedback first.

- Unit tests validate module behavior.
- Contract tests validate boundaries between apps, services, and packages.
- Integration tests validate database, cache, and external-service adapters.
- End-to-end tests should be reserved for critical user journeys.

The foundation currently includes smoke tests for the web shell, Go API health checks, and Python AI health checks.

