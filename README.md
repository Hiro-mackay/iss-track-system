# iss-tracker

Real-time ISS tracking system with orbital prediction and map visualization.
Displays the International Space Station's current position on an interactive map using TLE orbital data and SGP4 propagation.

## Tech Stack

- **Backend**: Go 1.22 (net/http + gorilla/websocket + akhenakh/sgp4)
- **Frontend**: Next.js (TypeScript, App Router, Tailwind CSS) + Leaflet (react-leaflet)
- **Real-time**: WebSocket (1-second interval)
- **TLE Data**: CelesTrak, 2-hour memory cache

## Project Structure

```
.
├── backend/                         # Go API server
├── frontend/                        # Next.js application
├── CLAUDE.md                        # AI assistant project config
├── docs/
│   ├── constitution.md              # Project rules and principles
│   ├── architecture.md              # System design (source of truth)
│   ├── decisions/                   # Architecture Decision Records
│   └── specs/                       # Feature specifications
└── .claude/                         # Claude Code config
```

## Development

```bash
# Backend
cd backend
go run .

# Frontend
cd frontend
npm install
npm run dev
```

## Guiding Documents

- [Constitution](docs/constitution.md) -- project rules and principles
- [Architecture](docs/architecture.md) -- single source of truth for system design

## SDD Workflow

Start new features with `/new-spec` in Claude Code.
