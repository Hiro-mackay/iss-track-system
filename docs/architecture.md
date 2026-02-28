# Architecture

> This is the single source of truth for the current system design.
> For past decisions and rationale, see [ADRs](decisions/).

## System Overview

iss-tracker is a real-time ISS tracking system. The backend fetches TLE (Two-Line Element) data from CelesTrak, computes the ISS position using SGP4 orbital propagation via akhenakh/sgp4, and streams coordinates to the frontend over WebSocket at 1-second intervals. The frontend renders the ISS position on a Cesium.js 3D globe with orbital path visualization at true orbital altitude.

## Component Architecture

- **Frontend** -- Next.js (App Router) SPA with Cesium.js 3D globe, receives real-time position updates via WebSocket
- **Backend API** -- Go net/http server with gorilla/websocket providing REST endpoints and WebSocket connections for ISS position streaming
- **Orbit Propagator** -- akhenakh/sgp4 module that computes ISS latitude/longitude/altitude from TLE data
- **TLE Cache** -- In-memory cache (2-hour TTL) that fetches and stores TLE data from CelesTrak
- **Crew Cache** -- In-memory cache (1-hour TTL) that fetches ISS crew data from Open Notify API
- **Station Info Panel** -- Left sidebar displaying ISS crew, station status, and orbital parameters

## Data Flow

1. Backend fetches TLE data from CelesTrak on startup and caches it (2-hour refresh)
2. Backend fetches ISS crew data from Open Notify API on startup and caches it (1-hour refresh)
3. Client connects to WebSocket endpoint
4. Backend propagates ISS position using SGP4 at current timestamp
5. Position (latitude, longitude, altitude) is sent to client every 1 second
6. Frontend updates the ISS entity on the Cesium 3D globe in real-time

## Bounded Contexts

| Context | Responsibility | Key Components |
|---------|---------------|----------------|
| Orbit Tracking | TLE fetch, SGP4 propagation, real-time streaming | TLECache, orbit.go, handler_position, handler_ws |
| Visualization | 3D globe rendering, orbit path display | GlobeView, ISSEntity, OrbitPath3D |
| Pass Prediction | Satellite pass calculations, visibility | passes.go, solar.go, handler_passes |
| Station Info | ISS crew, orbital parameters, station status | CrewCache, station.go, handler_crew, handler_status |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /api/v1/position | Current ISS position |
| GET | /api/v1/orbit | Orbit track points |
| GET | /api/v1/passes | Pass predictions for observer location |
| GET | /api/v1/iss/crew | Current ISS crew members |
| GET | /api/v1/iss/status | ISS status and orbital parameters |
| WS | /ws/position | Real-time position stream |

## Technology Stack

| Layer         | Technology | Notes |
|---------------|-----------|-------|
| Language      | Go 1.22, TypeScript | Backend + Frontend |
| Framework     | net/http + gorilla/websocket, Next.js (App Router) | REST + WebSocket, React SSR |
| Database      | None (in-memory cache) | TLE + crew data cached in memory |
| Infrastructure| TBD | Local development initially |
| CI/CD         | GitHub Actions | TBD |

## Infrastructure

Local development setup initially. Backend and frontend run as separate processes in a monorepo structure (`/backend` and `/frontend`).
