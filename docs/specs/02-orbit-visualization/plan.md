# Plan: Orbit Path Visualization

## Spec Reference

[Spec](../02-orbit-visualization/spec.md)

## Architecture Decisions

- 軌道トラック計算は既存の orbit.go に GetOrbitTrack() を追加
- OrbitPath は color/opacity を props で受け取り、GroundTrack が再利用
- 日付変更線分割は純粋関数 splitAtDateLine として実装

## Implementation Approach

Backend に軌道トラック API を追加し、Frontend で OrbitPath + GroundTrack コンポーネントを実装。splitAtDateLine はユーティリティ関数として OrbitPath 内に実装。

## Component Breakdown

### Component 1: Orbit Track API

- **Purpose**: 1 分間隔の軌道ポイント配列を返す
- **Location**: `backend/orbit.go` (追加), `backend/handler_orbit.go` (新規)
- **Changes**: GetOrbitTrack(minutes=90), GET /api/v1/orbit

### Component 2: OrbitPath

- **Purpose**: 軌道パスの Polyline 描画 (日付変更線分割付き)
- **Location**: `frontend/src/components/map/OrbitPath.tsx`
- **Changes**: splitAtDateLine + Polyline segments

### Component 3: GroundTrack

- **Purpose**: 将来軌道の表示 (OrbitPath 再利用)
- **Location**: `frontend/src/components/map/GroundTrack.tsx`
- **Changes**: OrbitPath を異なる opacity で呼び出し

## API Changes

| Method | Path | Response |
|--------|------|----------|
| GET | /api/v1/orbit?minutes=90 | OrbitPoint[] |

## Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| 日付変更線の描画異常 | Low | High | splitAtDateLine で確実に分割 |

## Dependencies

- 01-realtime-tracking が完了していること

## Out of Scope

- 3D 軌道表示
- 過去の軌道履歴保存
