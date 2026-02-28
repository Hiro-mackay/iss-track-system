# Tasks: Real-Time ISS Position Tracking

## References

- [Spec](../01-realtime-tracking/spec.md)
- [Plan](../01-realtime-tracking/plan.md)

## Task Breakdown

### Phase 1: Backend Setup

- [x] Task 1.1: プロジェクト構造作成 (S)
  - go mod init, go get (akhenakh/sgp4, gorilla/websocket)

- [x] Task 1.2: Config & Models (S)
  - config.go: os.Getenv (ISS_ prefix)
  - models.go: ISSPosition, OrbitPoint, PassEvent, PassPrediction struct

### Phase 2: Backend Core

- [x] Task 2.1: TLE Fetcher (S)
  - Depends on: 1.2
  - net/http client, 3LE パース, sync.RWMutex + time.Time キャッシュ

- [x] Task 2.2: Orbit Service (M)
  - Depends on: 2.1
  - GetCurrentPosition(), sgp4.ParseTLE + FindPositionAtTime

- [x] Task 2.3: REST Handler + WebSocket (S)
  - Depends on: 2.2
  - GET /api/v1/position, GET /health, WS /ws/position (gorilla/websocket)

- [x] Task 2.4: Main & Server (S)
  - Depends on: 2.3
  - goroutine TLE refresh, CORS middleware, ServeMux, graceful shutdown

### Phase 3: Frontend Setup

- [x] Task 3.1: Next.js プロジェクト初期化 (M)
  - create-next-app, react-leaflet/leaflet/swr 追加

- [x] Task 3.2: Types & API Client (S)
  - Depends on: 3.1
  - types.ts, api.ts (fetch + WS URL)

### Phase 4: Frontend Core

- [x] Task 4.1: useISSPosition Hook (M)
  - Depends on: 3.2
  - WebSocket + exponential backoff + REST フォールバック

- [x] Task 4.2: ISSMap & ISSMarker (M)
  - Depends on: 3.2
  - dynamic import, CARTO Dark Matter, SVG divIcon

- [x] Task 4.3: App Layout (S)
  - Depends on: 4.1, 4.2
  - page.tsx, layout.tsx, globals.css

### Phase 5: Integration

- [x] Task 5.1: Backend Tests (M)
  - Depends on: 2.2
  - MOCK_TLE, 位置範囲検証 (lat/lon/alt/vel)

- [x] Task 5.2: Docker Compose (S)
  - Depends on: 2.4, 4.3
  - Dockerfile (backend/frontend), docker-compose.yml

## Estimation Guidance

| Size | Description | Typical Duration |
|------|-------------|-----------------|
| S | Single file, isolated change | < 1 hour |
| M | 2-3 files, minor integration | 1-4 hours |
