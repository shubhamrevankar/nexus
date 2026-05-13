# Nexus System Design

This document explains what Nexus is, which pieces exist today, which pieces are planned, and how they relate.

Nexus is an AI-native operational intelligence platform. The long-term product connects company tools such as GitHub, Slack, docs, CI/CD, logs, and project management systems into one AI workspace where teams can ask questions, investigate failures, automate workflows, and understand engineering operations.

## Current Foundation

The repository currently contains the engineering foundation, not the full product.

```mermaid
flowchart LR
  Developer[Developer]

  subgraph Repo[Nexus Monorepo]
    Web[apps/web\nNext.js workspace shell]
    API[services/api\nGo API health service]
    AI[services/ai\nPython AI health service]
    Packages[packages/*\nShared config, logger, types]
    Docs[docs + architecture + adr + rfc\nEngineering knowledge base]
    Infra[Docker Compose + CI\nLocal platform tooling]
  end

  Developer --> Web
  Developer --> API
  Developer --> AI
  Web --> Packages
  API --> Infra
  AI --> Infra
  Repo --> Docs
```

Today, the system proves that the platform can run locally with separate web, API, AI, PostgreSQL, and Redis containers. Product features such as login, organizations, chat, search, integrations, and automations are intentionally not built yet.

## Runtime Architecture

This is the intended high-level runtime shape as product features are added.

```mermaid
flowchart TB
  User[User / Team Member]

  subgraph Client[Client Layer]
    WebApp[Web App\nNext.js + React]
  end

  subgraph Core[Core Platform]
    API[Core API\nGo modular monolith]
    Auth[Auth + Identity]
    Org[Organizations + Workspaces]
    Search[Knowledge Search]
    Workflow[Workflow Automation]
    Incident[Incident Management]
    Analytics[Engineering Analytics]
  end

  subgraph AIPlane[AI Plane]
    AIService[AI Service\nPython]
    RAG[RAG Pipelines]
    Agents[AI Agents]
    ModelGateway[Model Gateway]
  end

  subgraph Data[Data Layer]
    Postgres[(PostgreSQL\nPrimary relational data)]
    Redis[(Redis\nCache + realtime coordination)]
    Vector[(Vector Index\nFuture semantic search)]
    ObjectStore[(Object Storage\nFuture raw documents/artifacts)]
  end

  subgraph External[Company Tools]
    GitHub[GitHub]
    Slack[Slack]
    DocsTool[Notion / Google Drive]
    CICD[CI/CD + Deployments]
    Observability[Logs / Metrics / Alerts]
    Jira[Jira / Linear]
  end

  User --> WebApp
  WebApp --> API
  API --> Auth
  API --> Org
  API --> Search
  API --> Workflow
  API --> Incident
  API --> Analytics
  API --> Postgres
  API --> Redis
  API --> AIService
  AIService --> RAG
  AIService --> Agents
  AIService --> ModelGateway
  RAG --> Vector
  RAG --> ObjectStore
  Search --> Postgres
  Search --> Vector
  Workflow --> External
  External --> API
```

## Main Product Concept

Nexus becomes the operational layer above a company’s fragmented tools.

Instead of manually checking GitHub, Slack, docs, CI/CD, logs, deployments, and tickets, a user asks Nexus a question or creates an automation. Nexus retrieves the right context, reasons over it, and either answers the user or triggers an approved workflow.

Examples:

- “Why did deployment fail yesterday?”
- “How does authentication work?”
- “Create an incident room and draft an RCA when production deployment fails.”
- “Show engineering bottlenecks and sprint risks for this week.”

## Core Modules

Nexus starts as a modular monolith. That means the backend is one deployable Go API at first, but internally it is separated into clear product modules.

```mermaid
flowchart LR
  subgraph API[Go Core API - Modular Monolith]
    Identity[Identity Module]
    Tenancy[Organization + Workspace Module]
    Integrations[Integrations Module]
    Knowledge[Knowledge Module]
    Automation[Workflow Module]
    Incidents[Incident Module]
    Notifications[Notification Module]
    Analytics[Analytics Module]
    Platform[Platform Module\nlogging, config, health, auth middleware]
  end

  Identity --> Tenancy
  Tenancy --> Integrations
  Integrations --> Knowledge
  Knowledge --> Automation
  Automation --> Notifications
  Incidents --> Notifications
  Analytics --> Knowledge
  Platform --> Identity
  Platform --> Tenancy
```

### Why modular monolith first?

- Easier local development.
- Easier debugging.
- Lower operational complexity.
- Cleaner module boundaries before extracting services.
- Better fit for an early-stage platform.

Microservices can come later only after module boundaries are stable and real scale requires extraction.

## How A User Question Will Work

The most important user flow is asking Nexus a question about the company.

```mermaid
sequenceDiagram
  participant User
  participant Web as Web App
  participant API as Go API
  participant Auth as Auth/Tenancy
  participant AI as Python AI Service
  participant Search as Knowledge Search
  participant DB as PostgreSQL
  participant Vector as Vector Index
  participant Tools as External Tools

  User->>Web: Ask question
  Web->>API: POST /workspace/:id/ask
  API->>Auth: Verify user, org, workspace access
  Auth-->>API: Access allowed
  API->>AI: Request answer with workspace context
  AI->>Search: Retrieve relevant company knowledge
  Search->>DB: Load metadata, permissions, entities
  Search->>Vector: Semantic retrieval
  Search-->>AI: Ranked context
  AI->>AI: Reason over question + context
  AI-->>API: Answer + citations + suggested actions
  API-->>Web: Response payload
  Web-->>User: Show answer and sources
  API-->>Tools: Optional approved action in later phases
```

Key rule: AI should not blindly act. For risky operations such as rollback, ticket creation, or incident escalation, Nexus should require explicit permissions, audit logs, and approval flows.

## How Company Data Will Enter Nexus

Nexus needs reliable ingestion before it can answer useful questions.

```mermaid
flowchart TB
  subgraph Sources[External Sources]
    GitHub[GitHub repos, PRs, commits]
    Slack[Slack messages]
    Docs[Docs and wikis]
    CI[CI/CD runs]
    Logs[Logs, alerts, incidents]
    PM[Linear/Jira tickets]
  end

  subgraph Ingestion[Ingestion Platform]
    Connectors[Connectors]
    Jobs[Ingestion Jobs]
    Normalize[Normalize + Deduplicate]
    Permissions[Apply workspace permissions]
  end

  subgraph Storage[Storage]
    Raw[Raw source records]
    Relational[(PostgreSQL entities)]
    Embeddings[(Vector embeddings)]
  end

  subgraph UseCases[Product Use Cases]
    Ask[AI question answering]
    Search[Semantic search]
    Insights[Engineering insights]
    Automation[Workflow triggers]
  end

  Sources --> Connectors
  Connectors --> Jobs
  Jobs --> Normalize
  Normalize --> Permissions
  Permissions --> Raw
  Permissions --> Relational
  Permissions --> Embeddings
  Relational --> Ask
  Embeddings --> Ask
  Relational --> Search
  Embeddings --> Search
  Relational --> Insights
  Relational --> Automation
```

## Local Development Topology

This is what `docker compose up --build` runs locally.

```mermaid
flowchart LR
  Browser[Browser\nlocalhost:3000]

  subgraph Docker[Docker Compose]
    Web[web container\nNext.js]
    API[api container\nGo]
    AI[ai container\nPython]
    Postgres[(postgres container\nlocalhost:5432)]
    Redis[(redis container\nlocalhost:6379)]
  end

  Browser --> Web
  Web --> API
  Web --> AI
  API --> Postgres
  API --> Redis
  AI --> Postgres
  AI --> Redis
```

Current health endpoints:

- Web: `http://localhost:3000`
- API: `http://localhost:8080/healthz`
- AI: `http://localhost:8090/healthz`

## Repository Map

```text
apps/web
  User interface and workspace experience.

services/api
  Core backend API. Starts as a Go modular monolith.

services/ai
  Python AI plane for future RAG, agents, embeddings, and model orchestration.

packages/config
  Shared TypeScript config helpers.

packages/logger
  Shared TypeScript structured logging helper.

packages/types
  Shared TypeScript contracts.

infra/docker
  Dockerfiles for local service builds.

docs, architecture, adr, rfc
  Engineering knowledge base and decision history.
```

## Build Phases

```mermaid
flowchart LR
  P0[Phase 0\nEngineering foundation] --> P1[Phase 1\nAuth + orgs + workspace]
  P1 --> P2[Phase 2\nIntegrations + ingestion]
  P2 --> P3[Phase 3\nKnowledge search + RAG]
  P3 --> P4[Phase 4\nAI workspace + citations]
  P4 --> P5[Phase 5\nWorkflow automation]
  P5 --> P6[Phase 6\nIncidents + analytics]
```

### Phase 0: Foundation

Already started.

- Monorepo setup.
- Web/API/AI service stubs.
- Docker Compose.
- CI workflow.
- Engineering docs.

### Phase 1: Auth, Organizations, Workspace

The first real product slice.

- Users can sign in.
- Users belong to organizations.
- Organizations contain workspaces.
- API has real persistence in PostgreSQL.
- Web app has basic workspace navigation.

Implementation notes for this slice live in `docs/auth-workspace-slice.md`.

### Phase 2: Integrations and Ingestion

- Connect GitHub first.
- Store repos, commits, pull requests, and metadata.
- Add ingestion jobs and sync status.

### Phase 3: Knowledge Search and RAG

- Convert indexed company data into searchable knowledge.
- Add embeddings and semantic retrieval.
- Return answers with sources.

### Phase 4: AI Workspace

- User-facing chat/workspace experience.
- Conversation history.
- Citations.
- Workspace-aware answers.

### Phase 5: Workflow Automation

- Trigger-based automations.
- Approval gates.
- Audit logs.
- Notifications.

### Phase 6: Incidents and Analytics

- Incident timelines.
- Deployment failure analysis.
- Engineering bottleneck dashboards.

## Key Design Rules

- Build vertical slices, not isolated layers.
- Keep the Go API modular before considering microservices.
- Keep the Python AI service focused on AI orchestration, retrieval, and model workflows.
- Store source-of-truth business data in PostgreSQL first.
- Use Redis for cache, coordination, and future realtime support.
- Add vector search only when ingestion and knowledge models exist.
- Require permissions, audit logs, and approval flows before AI performs external actions.
