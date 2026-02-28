# Architecture

> This is the single source of truth for the current system design.
> For past decisions and rationale, see [ADRs](decisions/).

## System Overview

iss-tracker is a real-time ISS tracking system. The backend fetches TLE (Two-Line Element) data from CelesTrak, computes the ISS position using SGP4 orbital propagation via akhenakh/sgp4, and streams coordinates to the frontend over WebSocket at 1-second intervals. The frontend renders the ISS position on a Cesium.js 3D globe with orbital path visualization at true orbital altitude.

## Component Architecture

- **Frontend** -- Next.js (App Router) SPA with Cesium.js 3D globe, receives real-time position updates via WebSocket
- **Backend API** -- Go net/http server with gorilla/websocket providing REST endpoints and WebSocket connections for ISS position streaming
- **Orbit Propagator** -- akhenakh/sgp4 module that computes ISS latitude/longitude/altitude from TLE data
- **TLE Repository** -- In-memory repository (2-hour TTL) that fetches and stores TLE data from CelesTrak
- **Crew Repository** -- In-memory repository (1-hour TTL) that fetches ISS crew data from Open Notify API
- **Station Info Panel** -- Left sidebar displaying ISS crew, station status, and orbital parameters

## Backend Structure

```
backend/
  cmd/server/main.go                        # Entrypoint, DI wiring
  internal/
    config/config.go                         # LoadConfig from env vars
    domain/
      model/                                 # Domain types + repository interfaces
        position.go                          # ISSPosition, OrbitPoint
        pass.go                              # PassEvent, PassPrediction, Brightness
        crew.go                              # CrewMember
        station.go                           # OrbitalParams, ISSStatus
        repository.go                        # TLEProvider, CrewProvider interfaces
      orbit/                                 # Pure domain logic (no I/O)
        propagator.go                        # PropagatePosition, PropagateOrbitTrack
        pass.go                              # PredictPasses, Brightness
        solar.go                             # SunAltitude, IsPassVisible
        status.go                            # ComputeStatus
    repository/                              # Infrastructure (implements domain interfaces)
      tle.go                                 # TLERepository (TLEProvider)
      crew.go                                # CrewRepository (CrewProvider)
    service/query/                           # CQRS query services (orchestrate repo + domain)
      position.go                            # PositionQueryService
      pass.go                                # PassQueryService
      crew.go                                # CrewQueryService
      station.go                             # StationQueryService
    presentation/                            # HTTP layer
      response.go                            # WriteError
      middleware/cors.go                     # CORSMiddleware
      handler/                               # Thin HTTP handlers
        position.go, orbit.go, passes.go
        crew.go, status.go, ws.go
```

**Layering rules:**
- Domain layer has no infrastructure dependencies (Dependency Inversion via interfaces)
- Repository implements domain interfaces
- Query services orchestrate repository + domain logic
- Handlers are thin: HTTP concerns only, delegate to query services

## Frontend Structure

```
frontend/src/
  app/
    layout.tsx, page.tsx, globals.css
  features/
    tracking/                                # ISS position + orbit visualization
      components/                            # GlobeView, ISSEntity, OrbitPath3D,
                                             # CameraControls, cesium-setup,
                                             # PositionPanel, CountryInfo
      hooks/                                 # useISSPosition, useOrbitTrack, useReverseGeocode
      api.ts, types.ts, index.ts
    passes/                                  # Pass prediction feature
      components/PassList.tsx
      hooks/usePassPredictions.ts
      api.ts, types.ts, index.ts
    station/                                 # Station info sidebar
      components/                            # StationSidebar, CrewSection,
                                             # StatusSection, OrbitalSection
      hooks/                                 # useCrewData, useStationStatus
      api.ts, types.ts, index.ts
  components/LocationInput.tsx               # Generic reusable
  hooks/useGeolocation.ts                    # Generic reusable
  lib/
    api-client.ts                            # API_BASE constant
    fallback-crew.ts                         # Offline fallback
```

**Feature organization:** Each feature colocates its components, hooks, API layer, and types with barrel exports via `index.ts`.

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
| Orbit Tracking | TLE fetch, SGP4 propagation, real-time streaming | TLERepository, orbit/propagator, handler/position, handler/ws |
| Visualization | 3D globe rendering, orbit path display | tracking/GlobeView, ISSEntity, OrbitPath3D |
| Pass Prediction | Satellite pass calculations, visibility | orbit/pass, orbit/solar, handler/passes |
| Station Info | ISS crew, orbital parameters, station status | CrewRepository, orbit/status, handler/crew, handler/status |

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
| Infrastructure| Docker Compose | Local + CI |
| CI/CD         | GitHub Actions | Lint, typecheck, test, build |

## Infrastructure

Monorepo with `/backend` and `/frontend`. Development via `task dev` (parallel). Full CI via `task check` (lint + typecheck + test + build). Docker Compose for containerized deployment.
