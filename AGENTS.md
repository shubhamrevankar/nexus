# Nexus Engineering Instructions

Nexus is a production-grade AI-native operational intelligence platform.

Treat this repository as startup-grade infrastructure and a long-term scalable platform codebase.

Core principles:
- Maintainability first.
- Modular monolith first; avoid premature microservices.
- Clean architecture boundaries.
- Observability-first engineering.
- Strong typing and explicit contracts.
- Documentation-first development.
- Production-grade code only.

Rules:
- Never introduce unnecessary complexity.
- Never bypass architectural boundaries.
- Never create tight coupling between modules.
- Prefer explicit interfaces over implicit behavior.
- Prefer vertical slices over disconnected layers.
- Prefer readability over cleverness.
- Update docs, types, and tests for meaningful changes.
- Keep modules independently testable.

Before implementation:
- Understand the existing architecture.
- Understand module boundaries.
- Understand contracts and interfaces.

Do not generate tutorial-style or toy implementations.

