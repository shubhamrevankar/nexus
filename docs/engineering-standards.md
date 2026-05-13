# Engineering Standards

Nexus is built as a production-grade platform, not a tutorial project.

## Principles

- Prefer maintainability over speed.
- Start as a modular monolith with clear internal contracts.
- Keep modules independently testable.
- Make observability part of the design, not an afterthought.
- Document significant architectural decisions with ADRs.
- Use RFCs for cross-cutting or high-impact changes.

## Code Expectations

- Keep implementations small, explicit, and readable.
- Avoid tight coupling between product modules.
- Add tests for behavior that could regress.
- Keep environment configuration explicit and documented.
- Prefer boring infrastructure until real scaling pressure exists.

## Review Expectations

- Validate local checks before merging.
- Update docs when behavior, workflows, or architecture changes.
- Avoid mixing unrelated refactors with feature work.

