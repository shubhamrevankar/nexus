# Auth, Organizations, and Workspace Slice

This is the first product vertical slice for Nexus.

## What It Adds

- User registration and login.
- PostgreSQL-backed users and sessions.
- Organization creation during registration.
- Initial workspace creation during registration.
- Authenticated workspace listing.
- Web screens for onboarding and workspace overview.

## API Endpoints

### Register

```http
POST /v1/auth/register
Content-Type: application/json
```

```json
{
  "email": "you@company.com",
  "name": "Ada Lovelace",
  "password": "password123",
  "organizationName": "Acme Labs",
  "workspaceName": "Engineering"
}
```

### Login

```http
POST /v1/auth/login
Content-Type: application/json
```

```json
{
  "email": "you@company.com",
  "password": "password123"
}
```

### Current User

```http
GET /v1/me
Authorization: Bearer <session-token>
```

### Workspaces

```http
GET /v1/workspaces
Authorization: Bearer <session-token>
```

## Local Flow

1. Run `docker compose up --build`.
2. Open `http://localhost:3000`.
3. Click `Create workspace`.
4. Register a user, organization, and workspace.
5. Land on `/workspace` and see persisted organization/workspace data.

## Current Security Scope

This is a development-grade auth foundation, not a complete production auth system yet.

Implemented now:

- Passwords are hashed with bcrypt.
- Session tokens are random and stored hashed in PostgreSQL.
- Authenticated endpoints require bearer tokens.

Still needed before production:

- Secure HTTP-only cookies or a hardened token storage strategy.
- Email verification and password reset.
- Rate limiting and brute-force protection.
- Session revocation UI.
- Audit logging for identity and membership changes.

