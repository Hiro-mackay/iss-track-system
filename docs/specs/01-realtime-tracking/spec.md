# Spec: Real-Time ISS Position Tracking

## Status

Approved

## Overview

ISS の現在位置を WebSocket でリアルタイム配信し、Leaflet ダークモード地図上のマーカーで表示するコア機能。Backend/Frontend の基盤セットアップを含む。

## Problem Statement

ISS の現在位置をリアルタイムで地図上に表示するアプリケーションの土台が必要。TLE データの取得・キャッシュ、SGP4 軌道計算、WebSocket 配信、地図レンダリングの一連のパイプラインを構築する。

## User Stories

### Story 1: リアルタイム位置確認

As a 宇宙ファン, I want to ISS の現在位置を地図上でリアルタイムに確認する, so that 今 ISS がどこを飛んでいるか直感的にわかる.

#### Acceptance Criteria

- **Given** ユーザーがトップページを開いた
  **When** ページがロードされる
  **Then** ダークモードの世界地図上に ISS マーカーが表示され、1 秒間隔で移動する

- **Given** WebSocket 接続が確立されている
  **When** 1 秒が経過する
  **Then** ISS の緯度・経度・高度・速度が更新される

### Story 2: 接続障害からの復旧

As a ユーザー, I want to WebSocket 切断時も位置表示が続く, so that ネットワーク不安定でも体験が途切れない.

#### Acceptance Criteria

- **Given** WebSocket が切断された
  **When** 再接続を試行中
  **Then** REST API ポーリングで位置表示を継続し、WS 復帰後に自動で切り替わる

- **Given** WebSocket が切断された
  **When** 再接続を試行する
  **Then** exponential backoff (1s -> 2s -> 4s -> ... -> 30s cap) で再接続する

## Functional Requirements

- [x] FR-1: CelesTrak から TLE (3LE) を取得し 2 時間メモリキャッシュ <!-- tle.go:48-84, config.go:31 TLECacheTTL=7200s -->
- [x] FR-2: akhenakh/sgp4 で ISS の現在位置 (lat, lon, alt, velocity) を計算 <!-- orbit.go:13-42 -->
- [x] FR-3: WS /ws/position で 1 秒間隔の ISSPosition JSON 配信 <!-- handler_ws.go -->
- [x] FR-4: REST GET /api/v1/position で現在位置取得 <!-- handler_position.go -->
- [x] FR-5: REST GET /health でヘルスチェック <!-- main.go:35 -->
- [x] FR-6: Leaflet 地図に CARTO Dark Matter タイルを使用 <!-- ISSMap.tsx:18 -->
- [x] FR-7: ISS マーカーは SVG カスタムアイコン (L.divIcon) <!-- ISSMarker.tsx:7-22 L.divIcon + SVG -->
- [x] FR-8: WebSocket 切断時 exponential backoff + REST フォールバック <!-- useISSPosition.ts:27-55 backoff 1s->30s cap + REST poll -->

## Non-Functional Requirements

- [x] NFR-1: WebSocket 配信レイテンシ < 100ms <!-- in-memory SGP4 + direct WS push, runtime verification needed -->
- [x] NFR-2: SGP4 計算は 1 秒未満で完了 <!-- sgp4 C lib binding, runtime verification needed -->
- [x] NFR-3: CelesTrak 障害時は最後のキャッシュ TLE で継続 <!-- tle.go:38-47 stale cache fallback with slog.Warn -->
- [x] NFR-4: CORS は設定可能なオリジンリストで制限 <!-- cors.go + config.go:33 ISS_CORS_ORIGINS -->
- [x] NFR-5: 環境変数で全設定管理 (ISS_ prefix) <!-- config.go:14-36 全てISS_prefix -->
- [x] NFR-6: 全ソースファイル 300 行以内 <!-- max: tle.go 84行 -->
- [x] NFR-7: ダークモード固定 <!-- globals.css dark bg + CARTO Dark Matter tile -->

## Boundaries

### Always Do

- TLE は CelesTrak 公式 API から取得
- 設定値は環境変数 (ISS_ prefix) で管理
- WebSocket と REST の両方を提供

### Never Do

- シークレットをソースコードにハードコード
- CelesTrak への過剰リクエスト (2h キャッシュ厳守)

## Open Questions

(None)
