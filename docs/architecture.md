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
      tracking/                              # Orbit Tracking bounded context
        position.go                          # ISSPosition, OrbitPoint
        provider.go                          # TLEProvider interface
        propagator.go                        # PropagatePosition, PropagateOrbitTrack
      pass/                                  # Pass Prediction bounded context
        prediction.go                        # PassEvent, PassPrediction, Brightness consts
        pass.go                              # PredictPasses, Brightness
        solar.go                             # SunAltitude, IsPassVisible, julianDay
      station/                               # Station Info bounded context
        crew.go                              # CrewMember, CrewProvider interface
        status.go                            # OrbitalParams, ISSStatus, ComputeStatus
    repository/                              # Infrastructure (implements domain interfaces)
      tle.go                                 # TLERepository (tracking.TLEProvider)
      crew.go                                # CrewRepository (station.CrewProvider)
    service/                                 # Query services (domain-first naming)
      error.go                               # TLEError (shared across services)
      tracking/query.go                      # QueryService (position + orbit)
      pass/query.go                          # QueryService (pass predictions)
      station/query.go                       # QueryService (crew + status)
    presentation/                            # HTTP layer
      response.go                            # WriteError
      middleware/cors.go                     # CORSMiddleware
      handler/                               # Thin HTTP handlers
        position.go, orbit.go, passes.go
        crew.go, status.go, ws.go
```

**Layering rules:**
- Domain layer colocates types and behavior per bounded context (no separate model/logic split)
- Domain layer has no infrastructure dependencies (Dependency Inversion via interfaces)
- Repository implements domain interfaces
- Query services orchestrate repository + domain logic (service/{domain}/query.go)
- Handlers are thin: HTTP concerns only, delegate to query services

## Frontend Structure

```
frontend/src/
  app/
    layout.tsx                               # Root layout (Server Component)
    page.tsx                                 # Home page (Server Component)
    _components/ISSTrackerClient.tsx          # Client Component: layout state only
  features/
    tracking/                                # ISS position + orbit visualization
      components/                            # TrackingView (encapsulates hooks),
                                             # GlobeView, ISSEntity, OrbitPath3D,
                                             # CameraControls, cesium-setup,
                                             # PositionPanel, CountryInfo
      hooks/                                 # useISSPosition, useOrbitTrack, useReverseGeocode
      api.ts, types.ts, index.ts
    passes/                                  # Pass prediction feature
      components/                            # PassesPanel (encapsulates hooks), PassList
      hooks/usePassPredictions.ts
      api.ts, types.ts, index.ts
    station/                                 # Station info sidebar (self-contained)
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

**Feature organization:** Each feature colocates its components, hooks, API layer, and types with barrel exports via `index.ts`. Data hooks are encapsulated within feature components (TrackingView, PassesPanel), not called at the page level. `page.tsx` is a Server Component that delegates to `ISSTrackerClient` for client-side composition.

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
| Orbit Tracking | TLE fetch, SGP4 propagation, real-time streaming | domain/tracking, TLERepository, service/tracking, handler/position, handler/ws |
| Visualization | 3D globe rendering, orbit path display | tracking/TrackingView, GlobeView, ISSEntity, OrbitPath3D |
| Pass Prediction | Satellite pass calculations, visibility | domain/pass, service/pass, handler/passes, PassesPanel |
| Station Info | ISS crew, orbital parameters, station status | domain/station, CrewRepository, service/station, handler/crew, handler/status |

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
