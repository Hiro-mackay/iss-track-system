# Spec: Position Info Panel

## Status

Approved

## Overview

ISS の詳細な位置情報 (緯度・経度・高度・速度) と上空の地域名をサイドバーに表示する。Nominatim Reverse Geocoding で地域名を取得し、toggle ボタンでサイドバーの表示/非表示を切り替える。

## Problem Statement

地図上のマーカーだけでは ISS の正確な座標値や速度がわからない。数値データと「今どの国の上空にいるか」を併せて表示し、より詳細な情報を提供する。

## User Stories

### Story 1: 位置情報の数値表示

As a ユーザー, I want to ISS の座標値・高度・速度を数値で見る, so that 正確な位置情報を把握できる.

#### Acceptance Criteria

- **Given** ISS の位置が更新された
  **When** サイドバーを確認する
  **Then** 緯度、経度、高度 (km)、速度 (km/h) が表示される

### Story 2: 上空地域名の表示

As a ユーザー, I want to ISS が今どの国/地域の上空にいるか知る, so that 地理的な文脈を理解できる.

#### Acceptance Criteria

- **Given** ISS の緯度経度が取得されている
  **When** 地域名を表示する
  **Then** Nominatim Reverse Geocoding で取得した国名/地域名が表示される

- **Given** ISS が海上にいる
  **When** 地域名を表示する
  **Then** 「Ocean」等の適切な表示がされる

### Story 3: サイドバーの開閉

As a ユーザー, I want to 情報パネルを非表示にできる, so that 地図を全画面で見られる.

#### Acceptance Criteria

- **Given** サイドバーが表示されている
  **When** toggle ボタンを押す
  **Then** サイドバーが非表示になり、地図が全画面になる

## Functional Requirements

- [x] FR-1: PositionPanel に緯度・経度・高度 (km)・速度 (km/h) を表示 <!-- PositionPanel.tsx:15-48 -->
- [x] FR-2: Nominatim Reverse Geocoding で ISS 上空の地域名を取得・表示 <!-- CountryInfo.tsx:29-31 -->
- [x] FR-3: サイドバーの toggle ボタンで表示/非表示を切り替え <!-- page.tsx useState + toggle button -->
- [x] FR-4: ISS が海上の場合は適切なフォールバック表示 <!-- CountryInfo.tsx:34 "Ocean" / :36 "Unknown location" -->

## Non-Functional Requirements

- [x] NFR-1: Nominatim リクエストは位置更新毎ではなく適切にデバウンス <!-- CountryInfo.tsx:10 DEBOUNCE_MS=10000 -->
- [x] NFR-2: 全ソースファイル 300 行以内 <!-- max: PassList.tsx 70行 -->

## Boundaries

### Always Do

- Nominatim の利用規約を遵守 (User-Agent 設定, 1 req/s 以下)

### Ask First

- Nominatim キャッシュ戦略の変更

### Never Do

- Nominatim への過剰リクエスト

## Open Questions

- [RESOLVED] Nominatim のレート制限 (1 req/s) に対し、10 秒デバウンスで対応済み
