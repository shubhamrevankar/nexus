# ADR-0001: Start With a Modular Monolith

## Status

Accepted

## Context

Nexus needs strong engineering boundaries without introducing premature distributed-system complexity.

## Decision

Nexus will start as a monorepo-based modular monolith with explicit internal contracts and independently testable modules.

## Consequences

- Local development remains simple.
- Module boundaries can mature before extraction.
- Future services can be split out when operational pressure justifies it.

