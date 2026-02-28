# Plan: Real-Time ISS Position Tracking

## Spec Reference

[Spec](../01-realtime-tracking/spec.md)

## Architecture Decisions

- ADR-001 -> ADR-002: Go (net/http + akhenakh/sgp4) backend, Next.js + Leaflet frontend
- TLE キャッシュはインメモリ (sync.RWMutex + time.Time)
- WebSocket を主通信、REST をフォールバック

## Implementation Approach

Backend と Frontend を並行開発。Backend (Go) は TLE 取得 → SGP4 計算 → WebSocket/REST 配信。Frontend は Leaflet 地図 + ISS マーカー + WebSocket Hook。Docker Compose で統合。

## Component Breakdown

### Component 1: Backend Config & Models

- **Purpose**: 環境変数管理とデータモデル定義
- **Location**: `backend/config.go`, `backend/models.go`
- **Changes**: os.Getenv (ISS_ prefix), ISSPosition/OrbitPoint/PassPrediction struct

### Component 2: TLE Fetcher

- **Purpose**: CelesTrak TLE 取得 + メモリキャッシュ
- **Location**: `backend/tle.go`
- **Changes**: net/http client, sync.RWMutex cache (2h TTL)

### Component 3: Orbit Service

- **Purpose**: SGP4 で現在位置計算
- **Location**: `backend/orbit.go`
- **Changes**: sgp4.ParseTLE + FindPositionAtTime + ToGeodetic, velocity 計算

### Component 4: WebSocket & REST

- **Purpose**: リアルタイム配信 + REST エンドポイント
- **Location**: `backend/handler_ws.go`, `backend/handler_position.go`
- **Changes**: gorilla/websocket 1s ticker loop, GET /position, GET /health

### Component 5: App Entry

- **Purpose**: Go サーバー起動、TLE 定期更新
- **Location**: `backend/main.go`
- **Changes**: goroutine (TLE fetch + periodic refresh), CORS middleware, ServeMux routing

### Component 6: Frontend Types & API

- **Purpose**: 型定義、API クライアント
- **Location**: `frontend/src/lib/types.ts`, `frontend/src/lib/api.ts`
- **Changes**: ISSPosition 型, fetch + WebSocket URL ヘルパー

### Component 7: WebSocket Hook

- **Purpose**: リアルタイム位置取得
- **Location**: `frontend/src/hooks/useISSPosition.ts`
- **Changes**: WS 接続, exponential backoff, REST フォールバック

### Component 8: Map & Marker

- **Purpose**: Leaflet 地図 + ISS 表示
- **Location**: `frontend/src/components/map/ISSMap.tsx`, `ISSMarker.tsx`
- **Changes**: dynamic import (SSR off), CARTO Dark Matter, SVG icon

### Component 9: App Shell

- **Purpose**: レイアウト、スタイル
- **Location**: `frontend/src/app/page.tsx`, `layout.tsx`, `globals.css`
- **Changes**: 地図フルスクリーン, ダークモード

## API Changes

| Method | Path | Response |
|--------|------|----------|
| GET | /api/v1/position | ISSPosition |
| GET | /health | {"status": "ok"} |
| WS | /ws/position | ISSPosition (1s) |

## Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| CelesTrak ダウン | Med | Low | キャッシュ TLE で継続 |
| akhenakh/sgp4 コミュニティ規模 | Low | Low | SGP4 は安定アルゴリズム、joshuaferrara/go-satellite をフォールバック |

## Out of Scope

- 軌道パス描画 (02-orbit-visualization)
- パス予測 (03-pass-predictions)
- 情報パネル (04-position-info-panel)
