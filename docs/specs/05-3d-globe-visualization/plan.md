# Plan: 3D Globe Visualization

## Spec Reference

[Spec](./spec.md)

## Architecture Decisions

- ADR-003 (要作成): Leaflet 2D から Cesium 3D Globe への移行
  - Cesium.js + resium (React bindings) を採用
  - Cesium ion token 不要構成 (OSM basemap)

## Implementation Approach

Leaflet (react-leaflet) を完全に Cesium (resium) に置換する。既存の hooks (useISSPosition, useOrbitTrack 等) はそのまま再利用し、コンポーネント層のみ入れ替える。サイドバーは HUD オーバーレイに変換する。

## Component Breakdown

### Component 1: GlobeView

- **Purpose**: Cesium Viewer のラッパー。ダークテーマ、scene lighting 有効化
- **Location**: `frontend/src/components/globe/GlobeView.tsx`
- **Changes**: 新規作成。Viewer props で UI widgets 非表示、scene.globe.enableLighting = true

### Component 2: ISSEntity

- **Purpose**: ISS の 3D マーカー。高度付き PointGraphics
- **Location**: `frontend/src/components/globe/ISSEntity.tsx`
- **Changes**: 新規作成。Cartesian3.fromDegrees(lon, lat, alt * 1000)

### Component 3: OrbitPath3D

- **Purpose**: 軌道パスの 3D ポリライン。past/future 色分け
- **Location**: `frontend/src/components/globe/OrbitPath3D.tsx`
- **Changes**: 新規作成。PolylineGraphics + 高度付き座標配列

### Component 4: HUD

- **Purpose**: 位置情報、地域名、パス予測のオーバーレイ
- **Location**: `frontend/src/components/hud/HUD.tsx`
- **Changes**: 新規作成。既存の PositionPanel/CountryInfo/PassList を CSS overlay で配置

### Component 5: CameraControls

- **Purpose**: カメラ追従モードの toggle
- **Location**: `frontend/src/components/globe/CameraControls.tsx`
- **Changes**: 新規作成。viewer.trackedEntity の切り替え

## API Changes

バックエンド API 変更なし。既存エンドポイントをそのまま利用。

## Data Model Changes

フロントエンドの types.ts に変更なし。OrbitPoint の altitude_km を 3D 座標変換に利用。

## Domain Model Changes

### Aggregate Design

MapView aggregate に CameraMode (free | follow) Value Object を追加。

### Domain Event Flow

| Event | Published By | Consumed By | Trigger |
|-------|-------------|-------------|---------|
| Tracking.PositionComputed | Backend WS | GlobeView -> ISSEntity | 1 秒間隔 |
| Visualization.OrbitPathUpdated | useOrbitTrack | GlobeView -> OrbitPath3D | 60 秒間隔 |

## Non-Functional Requirements Strategy

| NFR | Technical Approach |
|-----|--------------------|
| NFR-1: 30fps | requestRenderMode で on-demand 描画、Entity 数最小化 |
| NFR-2: 静的アセット | next.config で CopyWebpackPlugin による Cesium assets コピー |
| NFR-3: 300行制限 | コンポーネント分割 (GlobeView, ISSEntity, OrbitPath3D, HUD, CameraControls) |
| NFR-4: 1秒更新 | 既存 useISSPosition hook をそのまま利用 |

## Codebase Alignment

- dynamic import パターン (既存 ISSMap と同じ ssr: false)
- hooks 層は変更なし (useISSPosition, useOrbitTrack, usePassPredictions, useGeolocation)
- ダークテーマ CSS は globals.css を維持

## Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| resium が React 19 未対応 | High | Medium | 互換性調査。fallback: cesium 直接操作 |
| Cesium バンドルサイズ (~30MB assets) | Medium | High | 静的配信 + CDN キャッシュ |
| Next.js 16 + Cesium webpack 設定の複雑さ | Medium | Medium | next.config.ts で CopyWebpackPlugin |

## Dependencies

- cesium (npm)
- resium (npm, React bindings for Cesium)
- Leaflet / react-leaflet の削除

## Out of Scope

- ISS の 3D .glb モデル表示 (ポイントマーカーで代替)
- 地形データ (Cesium World Terrain) の表示
- VR / AR モード
