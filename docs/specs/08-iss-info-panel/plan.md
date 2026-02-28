# Plan: ISS Information Panel

## Spec Reference

[Spec](./spec.md)

## Architecture Decisions

- [ADR-004](../../decisions/004-station-info-context.md): Station Info バウンデッドコンテキスト新設 + Open Notify API 採用

## Implementation Approach

バックエンド (Go) に 2 つの REST エンドポイント (`/api/v1/iss/crew`, `/api/v1/iss/status`) を追加し、フロントエンド (Next.js) に左サイドバーの ISS 情報パネルを構築する。バックエンドは既存の TLECache パターンを踏襲した CrewCache を新設し、フロントエンドは既存の SWR hook + HUD パネルパターンを踏襲する。

## Component Breakdown

### Component 1: CrewCache

- **Purpose**: Open Notify API からの ISS クルーデータ取得・キャッシュ (1h TTL)
- **Location**: `backend/crew.go`
- **Changes**: 新規作成。TLECache (tle.go) と同一パターン: RWMutex, Get/Refresh, graceful degradation

### Component 2: HandleCrew

- **Purpose**: `GET /api/v1/iss/crew` ハンドラ
- **Location**: `backend/handler_crew.go`
- **Changes**: 新規作成。CrewCache.Get() → JSON レスポンス

### Component 3: HandleStatus

- **Purpose**: `GET /api/v1/iss/status` ハンドラ。TLE から軌道パラメータを算出
- **Location**: `backend/handler_status.go`
- **Changes**: 新規作成。sgp4.ParseTLE → tle.Inclination, MeanMotion, Eccentricity, RevolutionNumber, RecoveredSemiMajorAxis

### Component 4: StationSidebar

- **Purpose**: 左サイドバーに配置する ISS 情報パネル (3 セクション合成)
- **Location**: `frontend/src/components/station/StationSidebar.tsx`
- **Changes**: 新規作成。CrewSection + StatusSection + OrbitalSection をアコーディオンで合成

### Component 5: SWR Hooks

- **Purpose**: クルーデータ・ステータスデータの取得 hook
- **Location**: `frontend/src/hooks/useCrewData.ts`, `useStationStatus.ts`
- **Changes**: 新規作成。usePassPredictions パターン踏襲

## API Changes

### `GET /api/v1/iss/crew` (新規)

```json
[
  { "name": "Oleg Kononenko", "craft": "ISS" }
]
```

- Open Notify API から取得、`craft == "ISS"` フィルタ
- 1 時間キャッシュ、失敗時はスティルデータ返却

### `GET /api/v1/iss/status` (新規)

```json
{
  "launch_year": 1998,
  "orbit_count": 157234,
  "crew_count": 7,
  "orbital_params": {
    "inclination_deg": 51.64,
    "period_min": 92.87,
    "eccentricity": 0.0001234,
    "apogee_km": 422.1,
    "perigee_km": 420.5
  },
  "tle_epoch": "2026-02-25T12:00:00Z"
}
```

- TLE から sgp4 ライブラリ経由で算出
- `orbit_count = RevolutionNumber + elapsed_min * MeanMotion / 1440`
- `period_min = 1440 / MeanMotion`
- `apogee_km = sma * (1+e) - 6371`, `perigee_km = sma * (1-e) - 6371`

## Data Model Changes

### Backend (models.go に追加)

```go
type CrewMember struct {
    Name  string `json:"name"`
    Craft string `json:"craft"`
}

type OrbitalParams struct {
    InclinationDeg float64 `json:"inclination_deg"`
    PeriodMin      float64 `json:"period_min"`
    Eccentricity   float64 `json:"eccentricity"`
    ApogeeKm       float64 `json:"apogee_km"`
    PerigeeKm      float64 `json:"perigee_km"`
}

type ISSStatus struct {
    LaunchYear    int           `json:"launch_year"`
    OrbitCount    int           `json:"orbit_count"`
    CrewCount     int           `json:"crew_count"`
    OrbitalParams OrbitalParams `json:"orbital_params"`
    TLEEpoch      string        `json:"tle_epoch"`
}
```

### Frontend (types.ts に追加)

```typescript
interface CrewMember { name: string; craft: string }
interface OrbitalParams { inclination_deg: number; period_min: number; eccentricity: number; apogee_km: number; perigee_km: number }
interface ISSStatus { launch_year: number; orbit_count: number; crew_count: number; orbital_params: OrbitalParams; tle_epoch: string }
```

## Domain Model Changes

### Aggregate Design

新規 Aggregate: **StationInfo**
- Root Entity: StationInfo
- Value Objects: CrewMember, OrbitalParams, StationStatus
- Invariants: CrewMember は `craft == "ISS"` のみ。OrbitalParams は現在の TLE キャッシュから算出

### Domain Event Flow

| Event | Published By | Consumed By | Trigger |
|-------|-------------|-------------|---------|
| StationInfo.CrewUpdated | CrewCache.Refresh | Visualization (crew panel) | 1 時間ごと |
| Tracking.TLERefreshed | TLECache.Refresh | HandleStatus (軌道パラメータ再算出) | 2 時間ごと |

## Non-Functional Requirements Strategy

| NFR | Technical Approach |
|-----|--------------------|
| NFR-1: クルー 1h キャッシュ | CrewCache に 1h TTL。Config で `ISS_CREW_CACHE_TTL_SECONDS` 設定可能 |
| NFR-2: 独立フォールバック | 各セクション (Crew/Status/Orbital) が独立した hook を使用。1 つの API 障害で他は影響なし |
| NFR-3: 追加外部リクエストなし | HandleStatus は既存 TLECache を参照するのみ。sgp4.ParseTLE はインメモリ計算 |
| NFR-4: 300 行制限 | 全新規ファイルは 80 行以内。StationSidebar は 3 セクションコンポーネントに分割 |

## Codebase Alignment

- **Cache**: `crew.go` は `tle.go` と同一パターン (RWMutex, TTL, graceful degradation)
- **Handler**: `handler_crew.go` / `handler_status.go` は `handler_position.go` と同一パターン (factory, slog, JSON encode)
- **Config**: `ISS_CREW_*` prefix で既存 `ISS_*` 命名規則に準拠
- **SWR Hook**: `useCrewData.ts` / `useStationStatus.ts` は `usePassPredictions.ts` パターン踏襲
- **Panel**: Row layout は `PositionPanel.tsx` の `flex justify-between gap-4` パターン踏襲
- **Sidebar**: `bg-black/70 rounded-lg p-4 backdrop-blur-sm` で既存右サイドバーに統一

## Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| Open Notify API ダウン | Medium | Medium | 3 層フォールバック: バックエンドスティル → SWR キャッシュ → フロントエンド静的データ |
| Open Notify レスポンス形式変更 | Medium | Low | ACL パターン (非公開構造体) で分離。パース失敗時はスティル返却 |
| RecoveredSemiMajorAxis() の単位誤り | Medium | Medium | テストで既知 TLE から apogee/perigee を検証 (ISS: 400-430km 範囲) |
| 左サイドバーが既存 UI と衝突 | Low | Medium | top-16 left-4 配置。既存左上ボタン群の下に配置 |

## Dependencies

- 新規外部依存なし (Open Notify API はバックエンドの HTTP クライアントで呼び出し)
- 既存依存: akhenakh/sgp4 (軌道パラメータ算出に利用)

## Out of Scope

- クルーの顔写真表示
- NASA ISS ライブストリーム埋め込み
- ISS 実験データベース連携
- クルーデータの永続化 (メモリキャッシュのみ)
