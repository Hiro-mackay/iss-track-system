# Tasks: 3D Globe Visualization

## References

- [Spec](./spec.md)
- [Plan](./plan.md)

## Task Breakdown

### Phase 1: Setup

- [ ] Task 1.1: Cesium + resium インストール & Next.js 設定 (M)
  - npm install cesium resium
  - next.config.ts に Cesium assets コピー設定 (CopyWebpackPlugin or cesium-webpack-plugin)
  - React 19 + resium 互換性検証
  - Files: package.json, next.config.ts

- [ ] Task 1.2: Leaflet 削除 (S)
  - Depends on: 1.1
  - npm uninstall leaflet react-leaflet @types/leaflet
  - 旧 map コンポーネント削除 (ISSMap, ISSMarker, OrbitPath, GroundTrack)
  - Files: package.json, src/components/map/ (削除)

### Phase 2: Core 3D Components

- [ ] Task 2.1: GlobeView コンポーネント (M)
  - Depends on: 1.1
  - Cesium Viewer ラッパー (ダークテーマ、scene lighting 有効)
  - dynamic import with ssr: false
  - Files: src/components/globe/GlobeView.tsx

- [ ] Task 2.2: ISSEntity コンポーネント (S)
  - Depends on: 2.1
  - ISSPosition を受け取り、高度付き 3D ポイントマーカーを配置
  - Cartesian3.fromDegrees(lon, lat, alt_km * 1000)
  - Files: src/components/globe/ISSEntity.tsx

- [ ] Task 2.3: OrbitPath3D コンポーネント (M)
  - Depends on: 2.1
  - OrbitPoint[] を高度付き 3D PolylineGraphics で描画
  - past (bright) / future (dim) 色分け
  - Files: src/components/globe/OrbitPath3D.tsx

- [ ] Task 2.4: CameraControls コンポーネント (S)
  - Depends on: 2.2
  - カメラ追従モード toggle (viewer.trackedEntity)
  - Files: src/components/globe/CameraControls.tsx

### Phase 3: HUD & Integration

- [ ] Task 3.1: HUD オーバーレイ (M)
  - Depends on: 2.1
  - 位置情報、地域名、パス予測を CSS overlay で表示
  - 既存 PositionPanel/CountryInfo/PassList ロジックを再利用
  - Files: src/components/hud/HUD.tsx

- [ ] Task 3.2: page.tsx 統合 (M)
  - Depends on: 2.1, 2.2, 2.3, 2.4, 3.1
  - GlobeView + ISSEntity + OrbitPath3D + CameraControls + HUD を配置
  - 既存 hooks を接続
  - Files: src/app/page.tsx

### Phase 4: Cleanup & Docs

- [ ] Task 4.1: 不要コード削除 (S)
  - Depends on: 3.2
  - Leaflet 関連 CSS/imports の残骸を削除
  - globals.css から Leaflet 固有スタイルを削除
  - Files: src/app/globals.css, 他

- [ ] Task 4.2: architecture.md 更新 (S)
  - Depends on: 3.2
  - Leaflet -> Cesium への技術スタック変更を反映
  - Files: docs/architecture.md

## Estimation Guidance

| Size | Description | Typical Duration |
|------|-------------|-----------------|
| S | Single file, isolated change | < 1 hour |
| M | 2-3 files, minor integration | 1-4 hours |
