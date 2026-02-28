# iss-tracker

Real-time ISS tracking system with orbital prediction and map visualization

## Tech Stack

- Monorepo: `/backend` + `/frontend`
- Backend: Go 1.22 (net/http + gorilla/websocket + akhenakh/sgp4)
- Frontend: Next.js (TypeScript, App Router, Tailwind CSS) + Leaflet (react-leaflet)
- Real-time: WebSocket (1-second interval)
- TLE Data: CelesTrak, 2-hour memory cache

## Guiding Documents

- [Constitution](docs/constitution.md) -- project rules and principles
- [Architecture](docs/architecture.md) -- single source of truth for system design

## SDD Workflow

Start new features with `/new-spec`.

## Commands

All commands via [Taskfile](https://taskfile.dev/):

- Dev: `task dev` (backend + frontend parallel)
- Build: `task build`
- Test: `task test`
- Lint: `task lint`
- Type check: `task typecheck`
- Full CI check: `task check` (lint + typecheck + test + build)
- Docker: `task up` / `task down`
- Setup: `task setup`

Subsystem-specific: `task dev:backend`, `task test:frontend`, etc.

## Domain Context

- Domain: ISS Orbital Tracking and Visualization
- Bounded Contexts: Orbit Tracking, Visualization, Pass Prediction
- Key Invariant: TLE data MUST be refreshed within 2 hours; orbit path MUST split at the date line

## Project Rules

<!-- Add project-specific rules below -->
