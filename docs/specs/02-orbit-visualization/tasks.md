# Tasks: Orbit Path Visualization

## References

- [Spec](../02-orbit-visualization/spec.md)
- [Plan](../02-orbit-visualization/plan.md)

## Task Breakdown

### Phase 1: Backend

- [x] Task 1.1: Orbit Track API (S)
  - orbit.go に GetOrbitTrack(minutes=90) 追加
  - handler_orbit.go に GET /api/v1/orbit エンドポイント追加

### Phase 2: Frontend

- [x] Task 2.1: OrbitPath コンポーネント (M)
  - Depends on: 1.1
  - splitAtDateLine 関数 (経度差 >180 で分割)
  - Polyline セグメント描画、color/opacity props

- [x] Task 2.2: GroundTrack コンポーネント (S)
  - Depends on: 2.1
  - OrbitPath を再利用 (異なる opacity)

### Phase 3: Integration

- [x] Task 3.1: page.tsx に軌道表示を統合 (S)
  - Depends on: 2.1, 2.2
  - fetchOrbit(90) を 60 秒間隔で更新
  - 前半=現在軌道、後半=将来軌道に分割

## Estimation Guidance

| Size | Description | Typical Duration |
|------|-------------|-----------------|
| S | Single file, isolated change | < 1 hour |
| M | 2-3 files, minor integration | 1-4 hours |
