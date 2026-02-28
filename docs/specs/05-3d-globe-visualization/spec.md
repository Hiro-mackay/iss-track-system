# Spec: 3D Globe Visualization

## Status

Approved

## Bounded Context

Visualization + Orbit Tracking (as defined in `docs/prd.md`)

## Overview

Leaflet 2D 地図を Cesium.js (resium) ベースの 3D 地球儀に完全置換する。ISS の位置と軌道を実際の高度 (~420km) で 3D 空間に描画し、昼夜境界線 (terminator) やカメラ追従モードを提供する。情報パネルは HUD オーバーレイとして 3D シーン内に統合する。

## Problem Statement

2D メルカトル地図では ISS の軌道の立体的な広がりや高度感が伝わらない。3D 地球儀を用いることで、ISS が地球の周りを周回している臨場感と、昼夜境界との位置関係を直感的に把握できるようにする。

## User Stories

### Story 1: 3D 地球儀での ISS 位置確認

As a 宇宙ファン, I want to 3D 地球儀上で ISS のリアルタイム位置を確認する, so that 地球を周回している臨場感を得られる.

#### Acceptance Criteria

- **Given** ユーザーがトップページを開いた
  **When** ページがロードされる
  **Then** 3D 地球儀上の実際の高度 (~420km) に ISS マーカーが表示され、1 秒間隔で移動する

- **Given** 3D 地球儀が表示されている
  **When** マウスドラッグ/ホイールで操作する
  **Then** 地球儀を自由に回転・ズームできる

### Story 2: 高度付き軌道パスの 3D 表示

As a ユーザー, I want to ISS の軌道パスを 3D 空間で見る, so that 軌道の立体的な形状を理解できる.

#### Acceptance Criteria

- **Given** ISS の位置が表示されている
  **When** 軌道データが取得される
  **Then** 90 分間の軌道パスが実際の高度で 3D ポリラインとして描画される

- **Given** 軌道パスが描画されている
  **When** 地球儀を回転させる
  **Then** 軌道パスが地球の曲面に沿って立体的に見える

### Story 3: 昼夜境界線の表示

As a ユーザー, I want to 地球の昼夜境界を見る, so that ISS が日照側にいるか影側にいるか分かる.

#### Acceptance Criteria

- **Given** 3D 地球儀が表示されている
  **When** シーンがレンダリングされる
  **Then** Cesium のシーンライティングにより昼夜境界が表示される

### Story 4: カメラ追従モード

As a ユーザー, I want to ISS にカメラを追従させる, so that ISS の視点から地球を見られる.

#### Acceptance Criteria

- **Given** 3D 地球儀が表示されている
  **When** 追従モードを有効にする
  **Then** カメラが ISS を中心に追従し、ISS の移動に合わせてビューが更新される

- **Given** 追従モードが有効
  **When** 追従モードを無効にする
  **Then** カメラが現在位置で固定され、自由操作に戻る

### Story 5: HUD オーバーレイ

As a ユーザー, I want to 位置情報やパス予測を 3D シーン内で確認する, so that ビューを切り替えずに情報を得られる.

#### Acceptance Criteria

- **Given** 3D 地球儀が表示されている
  **When** ISS の位置が更新される
  **Then** HUD に緯度/経度/高度/速度/地域名が表示される

- **Given** 観測地点が設定されている
  **When** HUD を確認する
  **Then** パス予測一覧が表示される

## Functional Requirements

- [ ] FR-1: Cesium.js + resium で 3D 地球儀描画 (ダークテーマ)
- [ ] FR-2: ISS を 3D Entity (ポイントまたはモデル) として実際の高度に配置
- [ ] FR-3: 軌道パスを PolylineGraphics で高度付き 3D 描画 (past/future 色分け)
- [ ] FR-4: Cesium scene lighting で昼夜境界 (terminator) 表示
- [ ] FR-5: カメラ追従モード (toggle) -- ISS Entity に camera.trackEntity
- [ ] FR-6: HUD オーバーレイで位置情報 (lat/lon/alt/velocity/timestamp) 表示
- [ ] FR-7: HUD オーバーレイで上空地域名表示 (Nominatim)
- [ ] FR-8: HUD オーバーレイでパス予測テーブル表示
- [ ] FR-9: 既存 Leaflet 2D 地図コンポーネントを完全削除
- [ ] FR-10: 既存 hooks (useISSPosition, useOrbitTrack, usePassPredictions, useGeolocation) を再利用

## Non-Functional Requirements

- [ ] NFR-1: 3D シーン描画は 30fps 以上を維持
- [ ] NFR-2: Cesium assets はビルド時に静的配信設定
- [ ] NFR-3: 全ソースファイル 300 行以内
- [ ] NFR-4: WebSocket リアルタイム更新は既存の 1 秒間隔を維持

## Boundaries

### Always Do

- Cesium ion token 不要の構成 (OSM / 自然地球タイル等)
- 既存 hooks を変更せず再利用
- ダークテーマを維持

### Ask First

- Cesium ion トークン (有料機能) の導入
- 3D モデル (ISS .glb) の追加

### Never Do

- Cesium ion 有料 API への依存 (トークンなし構成)
- 既存バックエンド API の変更

## Technical Constraints

- Next.js で Cesium の静的アセット (Workers, Assets) 配信に webpack/next.config 設定が必要
- resium は React 19 対応を確認する必要あり
- Cesium はクライアントサイドのみ (SSR 不可、dynamic import 必須)

## Domain Impact

### Affected Aggregates

- **MapView**: OrbitPath -> 3D PolylineGraphics (高度付き)、MarkerPosition -> 3D Entity
- New Value Objects: CameraMode (free | follow)

### Domain Events

- Consumed: Tracking.PositionComputed (from Orbit Tracking) -- ISS マーカー更新
- Consumed: Visualization.OrbitPathUpdated (from Visualization) -- 軌道パス再描画

### Ubiquitous Language Additions

| Term | Definition | Context |
|------|------------|---------|
| Globe View | Cesium ベースの 3D 地球儀レンダリングビュー | Visualization |
| Terminator | 地球の昼夜境界線。太陽位置から計算される | Visualization |
| Camera Follow Mode | カメラが ISS Entity を自動追従するモード | Visualization |
| HUD (Head-Up Display) | 3D シーン上にオーバーレイ表示する情報パネル | Visualization |

## Open Questions

- [RESOLVED] Cesium ion token: 不要構成で進める (OSM basemap)
