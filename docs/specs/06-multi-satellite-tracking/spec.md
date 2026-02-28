# Spec: Multi-Satellite Tracking

## Status

Draft

## Bounded Context

Orbit Tracking + Visualization (as defined in `docs/prd.md`)

## Overview

ISS 単体のトラッキングを拡張し、CelesTrak から複数衛星の TLE を取得・管理して同時追跡を可能にする。ユーザーはカタログから衛星を選択し、3D グローブ上に複数衛星を色分け表示できる。バックエンドの Satellite Aggregate を汎用化し、NORAD Catalog Number ベースで任意の衛星を扱えるようにする。

## Problem Statement

現在のシステムは ISS (NORAD ID: 25544) 専用で設計されている。Hubble (20580)、Tiangong (54216)、Starlink 衛星群など他の人工衛星にも関心を持つユーザーが多いが、衛星ごとに別システムを構築するのは非効率である。Satellite Aggregate を汎用化することで、単一システムで複数衛星の同時追跡を実現する。

## User Stories

### Story 1: 衛星カタログからの選択

As a 宇宙ファン, I want to カタログから追跡したい衛星を選択する, so that ISS 以外の衛星もリアルタイムで追跡できる.

#### Acceptance Criteria

- **Given** ユーザーが衛星セレクタ UI を開いた
  **When** カタログ一覧が表示される
  **Then** プリセット衛星 (ISS, Hubble, Tiangong) と NORAD ID 入力フィールドが表示される

- **Given** ユーザーがプリセットから Hubble を選択した
  **When** 選択が確定される
  **Then** Hubble の TLE が取得され、3D グローブ上にマーカーが追加される

### Story 2: 複数衛星の同時表示

As a ユーザー, I want to 複数の衛星を同時に 3D グローブ上で見る, so that 衛星間の位置関係を把握できる.

#### Acceptance Criteria

- **Given** ISS と Hubble が選択されている
  **When** 3D グローブが表示される
  **Then** 両方の衛星がそれぞれ異なる色のマーカーで表示される

- **Given** 複数衛星が表示されている
  **When** ある衛星をクリックする
  **Then** その衛星が「アクティブ衛星」になり、HUD に当該衛星の位置情報が表示される

### Story 3: 衛星ごとの軌道パス表示

As a ユーザー, I want to 各衛星の軌道パスを個別に表示する, so that 軌道の違いを比較できる.

#### Acceptance Criteria

- **Given** 複数衛星が選択されている
  **When** 軌道パスが描画される
  **Then** 各衛星の軌道が衛星カラーに合わせた色で描画される

- **Given** 3 衛星以上が選択されている
  **When** パフォーマンスが低下する恐れがある
  **Then** 同時追跡数の上限 (5 衛星) で警告が表示される

## Functional Requirements

- [ ] FR-1: バックエンドの TLECache を複数衛星対応 (NORAD ID をキーとするマップ) に拡張
- [ ] FR-2: `GET /api/v1/satellites` -- プリセット衛星カタログを返す REST エンドポイント
- [ ] FR-3: `GET /api/v1/position?norad_id=25544` -- NORAD ID 指定で任意の衛星位置を返す
- [ ] FR-4: `GET /api/v1/orbit?norad_id=25544` -- NORAD ID 指定で軌道パスを返す
- [ ] FR-5: `GET /ws/position?satellites=25544,20580` -- 複数衛星の位置を同時ストリーミング
- [ ] FR-6: フロントエンドに衛星セレクタ UI (プリセット + NORAD ID 手動入力)
- [ ] FR-7: 3D グローブ上に複数衛星を色分けマーカーで描画
- [ ] FR-8: アクティブ衛星の切り替え -- HUD 表示が選択衛星に連動
- [ ] FR-9: 各衛星の軌道パスを衛星カラーで個別描画
- [ ] FR-10: 同時追跡数の上限を 5 衛星とし、上限到達時に UI で警告

## Non-Functional Requirements

- [ ] NFR-1: Performance -- 5 衛星同時表示でも 30fps 以上を維持
- [ ] NFR-2: Performance -- TLE フェッチは衛星あたり個別に失敗しても他衛星に影響しない
- [ ] NFR-3: Performance -- CelesTrak への同時リクエストは最大 3 並列に制限
- [ ] NFR-4: Compatibility -- 既存の ISS 単体 API (`/api/v1/position`, `/ws/position`) は後方互換を維持

## Boundaries

### Always Do

- ISS はデフォルトで常に選択済み (削除不可)
- 各衛星の TLE キャッシュは独立管理 (2 時間 TTL)
- CelesTrak の利用規約を遵守 (過度なリクエスト禁止)
- 既存 API の後方互換性を維持

### Ask First

- CelesTrak 以外の TLE ソース (Space-Track.org 等) の追加
- 衛星グループ一括取得 (Starlink 全衛星等)
- 衛星の 3D モデル (.glb) 表示

### Never Do

- CelesTrak への 1 秒間隔未満のリクエスト
- ユーザー認証なしでの衛星データ永続化
- 同時追跡 5 衛星超

## Technical Constraints

- CelesTrak 3LE フォーマット: `https://celestrak.org/NORAD/elements/gp.php?CATNR={id}&FORMAT=3LE`
- 既存の akhenakh/sgp4 ライブラリをそのまま利用可能 (TLE 文字列を渡すだけ)
- WebSocket メッセージフォーマットを拡張する際は `satellite_id` フィールドを追加
- フロントエンドのグローブコンポーネントは Cesium Entity を衛星数分動的に追加

## Domain Impact

### Affected Aggregates

- **Satellite**: 単一インスタンスからマルチインスタンス対応に拡張。NORAD ID を識別子に追加
- **MapView**: 複数 Entity / PolylineGraphics の動的管理。アクティブ衛星の概念を追加

### Domain Events

- Published: Tracking.TLERefreshed -- 衛星ごとに個別発行
- Published: Tracking.PositionComputed -- `satellite_id` 付きで発行
- Published: Tracking.SatelliteAdded -- 新しい衛星が追跡対象に追加されたとき
- Published: Tracking.SatelliteRemoved -- 衛星が追跡対象から削除されたとき

### Ubiquitous Language Additions

| Term | Definition | Context |
|------|------------|---------|
| NORAD Catalog Number | NORAD が割り当てる衛星の一意識別番号 (例: ISS = 25544) | Orbit Tracking |
| Satellite Catalog | 追跡可能な衛星のプリセット一覧 | Orbit Tracking |
| Active Satellite | HUD 情報表示の対象として選択されている衛星 | Visualization |
| Satellite Color | 各衛星に割り当てられた識別色 (マーカー・軌道パスに使用) | Visualization |

## Open Questions

- [NEEDS CLARIFICATION] パス予測も衛星ごとに対応するか、ISS 限定のままとするか
- [NEEDS CLARIFICATION] 衛星プリセットの初期ラインナップ (ISS, Hubble, Tiangong 以外に何を含めるか)
