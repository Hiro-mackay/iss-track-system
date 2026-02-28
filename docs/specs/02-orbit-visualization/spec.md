# Spec: Orbit Path Visualization

## Status

Approved

## Overview

ISS の軌道パス (1 周約 90 分) を地図上に Polyline で描画する。日付変更線 (経度 180/-180) をまたぐ場合の分割処理を含む。

## Problem Statement

ISS マーカーだけでは移動方向や軌道全体が把握しにくい。現在の軌道パスと将来の地上トラックを線として可視化し、ISS がどこから来てどこへ向かうか一目でわかるようにする。

## User Stories

### Story 1: 軌道パスの表示

As a ユーザー, I want to ISS の 90 分間の軌道パスを地図上に見る, so that 軌道全体を把握できる.

#### Acceptance Criteria

- **Given** ISS の位置が表示されている
  **When** 軌道データが取得される
  **Then** 90 分間の軌道パスが Polyline で描画される

- **Given** 軌道パスが日付変更線をまたぐ
  **When** 描画される
  **Then** Polyline が分割され、地図を横断する異常な線が出ない

### Story 2: 軌道データの自動更新

As a ユーザー, I want to 軌道データが定期的に更新される, so that 常に最新の軌道パスが見える.

#### Acceptance Criteria

- **Given** 軌道パスが表示されている
  **When** 60 秒が経過する
  **Then** 軌道データが自動更新される

## Functional Requirements

- [x] FR-1: REST GET /api/v1/orbit?minutes=90 で軌道ポイント配列取得 (1 分間隔) <!-- handler_orbit.go, orbit.go:44-69 -->
- [x] FR-2: OrbitPath コンポーネントで Polyline 描画 <!-- OrbitPath.tsx:38-57 -->
- [x] FR-3: splitAtDateLine 関数で経度差 >180 時に Polyline を分割 <!-- OrbitPath.tsx:9-28 -->
- [x] FR-4: GroundTrack コンポーネントで将来軌道を OrbitPath 再利用 (異なる opacity) <!-- GroundTrack.tsx:13-20 opacity 0.8/0.3 -->
- [x] FR-5: 60 秒間隔で軌道データを自動更新 <!-- useOrbitTrack.ts:9 refreshInterval: 60000 -->

## Non-Functional Requirements

- [x] NFR-1: 90 ポイントの Polyline 描画がスムーズ (フレーム落ちなし) <!-- 90 points minimal for Leaflet -->
- [x] NFR-2: 全ソースファイル 300 行以内 <!-- max: OrbitPath.tsx 57行 -->

## Boundaries

### Always Do

- 日付変更線での Polyline 分割処理
- color/opacity を props で受け取り再利用可能に

### Never Do

- 分割なしの単一 Polyline (日付変更線で壊れる)

## Open Questions

(None)
