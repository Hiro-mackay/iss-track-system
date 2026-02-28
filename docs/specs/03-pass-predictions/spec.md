# Spec: Pass Predictions

## Status

Approved

## Overview

ユーザーの観測地点から ISS がいつ可視になるかを予測し、一覧表示する。akhenakh/sgp4 の GeneratePasses でパス予測を行い、ブラウザ Geolocation で観測地点を自動取得する。

## Problem Statement

ISS を肉眼で観測するには「衛星が日照で観測地が薄暮以下」の条件を満たすタイミングを事前に知る必要がある。軌道計算から可視パスを予測し、rise/culmination/set の時刻と仰角を提供する。

## User Stories

### Story 1: 可視パスの確認

As a 天体観測者, I want to 自分の場所から ISS がいつ見えるかを知る, so that 観測の準備ができる.

#### Acceptance Criteria

- **Given** 観測地点が設定されている
  **When** パス予測を表示する
  **Then** 3 日分の可視パス一覧が表示される (rise/culmination/set 時刻、最大仰角、明るさ)

- **Given** パスが可視条件を満たす (衛星日照 + 観測地薄暮以下)
  **When** 一覧に表示される
  **Then** is_visible=true、brightness が bright/moderate/dim で分類されている

### Story 2: 観測地点の自動取得

As a ユーザー, I want to ブラウザから自分の位置を自動取得したい, so that 手動入力なしでパス予測が見られる.

#### Acceptance Criteria

- **Given** ブラウザの位置情報を許可した
  **When** ページがロードされる
  **Then** navigator.geolocation の座標が観測地点として使用される

- **Given** 位置情報を拒否した
  **When** ページがロードされる
  **Then** デフォルト (東京: 35.68, 139.69) が使用され、手動入力も可能

## Functional Requirements

- [x] FR-1: akhenakh/sgp4 GeneratePasses で rise/culmination/set を検出 <!-- passes.go:21 GeneratePasses -->
- [x] FR-2: is_sunlit + 太陽高度 < -6 度で可視判定 <!-- solar.go:IsPassVisible + passes.go:39 -->
- [x] FR-3: 明るさ分類 (max_elevation >= 60: bright, >= 30: moderate, else: dim) <!-- passes.go:57-66 brightness() -->
- [x] FR-4: REST GET /api/v1/passes?lat=...&lon=...&days=3 でパス予測取得 <!-- handler_passes.go -->
- [x] FR-5: navigator.geolocation で観測地点を自動取得 <!-- useGeolocation.ts:29 -->
- [x] FR-6: Geolocation 拒否時はデフォルト東京 + 手動入力 UI <!-- useGeolocation.ts + LocationInput.tsx -->
- [x] FR-7: SWR で 5 分間隔の自動更新 (usePassPredictions hook) <!-- usePassPredictions.ts:9 refreshInterval: 300000 -->
- [x] FR-8: PassList コンポーネントでパス一覧テーブル表示 <!-- PassList.tsx -->

## Non-Functional Requirements

- [x] NFR-1: パス予測計算はオンデマンド実行 <!-- handler_passes.go でリクエスト毎に計算 -->
- [x] NFR-2: パス予測計算は 3 日分で 5 秒以内 <!-- sgp4 C binding, runtime verification needed -->
- [x] NFR-3: 全ソースファイル 300 行以内 <!-- max: PassList.tsx 70行 -->

## Boundaries

### Always Do

- sgp4.Location で observer 位置を作成
- GeneratePasses で可視パスを計算

### Ask First

- パス予測の仰角しきい値の変更 (現在 10 度)
- 予測日数の変更 (現在 3 日)

### Never Do

- サーバーサイドで個人の位置情報を保存

## Open Questions

(None)
