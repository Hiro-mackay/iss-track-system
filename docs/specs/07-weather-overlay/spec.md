# Spec: Weather Overlay

## Status

Draft

## Bounded Context

Visualization + Pass Prediction (as defined in `docs/prd.md`)

## Overview

観測地点の天候情報 (雲量・天気) を取得し、パス予測に「観測条件」を付加する。3D グローブ上にはクラウドカバー (雲量) タイルレイヤーをオーバーレイ表示する。天候が悪い場合は ISS が見えない可能性が高いため、パス予測の実用性を大幅に向上させる。

## Problem Statement

現在のパス予測は天文学的条件 (仰角・太陽位置・衛星照度) のみで可視判定を行っている。しかし現実の観測では雲量が最大の障害となる。晴天時にしか ISS を肉眼で視認できないにもかかわらず、曇天のパスも「Visible」として表示されるため、ユーザーが無駄な待機をする問題がある。

## User Stories

### Story 1: パス予測に天候条件を表示

As a ISS 観測者, I want to パス予測に天候情報が付加されている, so that 実際に見えるかどうか事前に判断できる.

#### Acceptance Criteria

- **Given** 観測地点が設定されている
  **When** パス予測一覧が表示される
  **Then** 各パスに予測時間帯の天候条件 (晴れ/曇り/雨) と雲量 (%) が表示される

- **Given** パス予測一覧が表示されている
  **When** 雲量が 70% 以上のパスがある
  **Then** 該当パスに「観測困難」のラベルが表示される

### Story 2: 現在の天候表示

As a ユーザー, I want to 観測地点の現在の天候を確認する, so that 今すぐ外に出て ISS を見る価値があるか判断できる.

#### Acceptance Criteria

- **Given** 観測地点が設定されている
  **When** 天候情報が取得される
  **Then** HUD に現在の天気アイコン、気温、雲量が表示される

- **Given** 天候 API がエラーを返した
  **When** HUD を確認する
  **Then** 天候情報は非表示 (フォールバック) となり、パス予測は天候なしで表示される

### Story 3: グローブ上の雲レイヤー

As a ユーザー, I want to 3D グローブ上で雲の分布を見る, so that 地球全体の天候パターンを視覚的に把握できる.

#### Acceptance Criteria

- **Given** 3D グローブが表示されている
  **When** 雲レイヤーが有効
  **Then** 半透明の雲タイルが地球表面にオーバーレイされる

- **Given** 雲レイヤーが有効
  **When** トグルボタンで無効にする
  **Then** 雲レイヤーが非表示になる

## Functional Requirements

- [ ] FR-1: バックエンドに天候取得エンドポイント `GET /api/v1/weather?lat=X&lon=Y` を追加
- [ ] FR-2: Open-Meteo API から観測地点の現在天候と 3 日間予報を取得
- [ ] FR-3: 天候データを 30 分間メモリキャッシュ (座標の 0.1 度丸めをキーに)
- [ ] FR-4: パス予測レスポンスに `weather_condition` (clear/cloudy/overcast/rain) と `cloud_cover_percent` を付加
- [ ] FR-5: HUD に現在の天候情報 (天気・気温・雲量) を表示するコンポーネント
- [ ] FR-6: パス予測テーブルに雲量カラムと「観測困難」ラベルを追加
- [ ] FR-7: Cesium ImageryLayer で雲タイルを半透明オーバーレイ表示
- [ ] FR-8: 雲レイヤーの表示/非表示トグルボタン

## Non-Functional Requirements

- [ ] NFR-1: Performance -- 天候 API 呼び出しは 30 分キャッシュで 1 日あたり最大 48 回/地点に制限
- [ ] NFR-2: Resilience -- 天候 API 障害時はグレースフルデグラデーション (天候情報なしで動作継続)
- [ ] NFR-3: Performance -- 雲タイルレイヤーの追加で描画が 25fps を下回らない

## Boundaries

### Always Do

- 天候 API には無料・トークン不要の Open-Meteo を使用
- 天候データ取得失敗時は天候なしで動作 (既存機能に影響しない)
- 雲量しきい値はハードコードせず環境変数で設定可能に

### Ask First

- 有料天候 API (OpenWeatherMap 等) への切り替え
- 天候に基づくプッシュ通知 (「今夜は晴れてパスがあります」)

### Never Do

- 天候 API の API キーをソースコードにハードコード
- 天候情報の永続化 (メモリキャッシュのみ)
- 天候 API への 1 分間隔未満のリクエスト

## Technical Constraints

- Open-Meteo Forecast API: `https://api.open-meteo.com/v1/forecast` (無料、トークン不要、10,000 req/day)
- パラメータ: `latitude`, `longitude`, `hourly=cloud_cover,weather_code,temperature_2m`, `forecast_days=3`
- 雲タイルは OpenWeatherMap Tile API または RainViewer API を使用 (どちらも無料枠あり)
- パス予測への天候統合はバックエンドで行い、フロントエンドは表示のみ

## Domain Impact

### Affected Aggregates

- **PassForecast**: WeatherCondition Value Object を追加。PassPrediction に天候フィールドを付加
- **MapView**: CloudLayer の表示状態を管理

### Domain Events

- Published: Visualization.WeatherUpdated -- 観測地点の天候データが更新されたとき
- Consumed: PassPrediction.ForecastComputed -- 天候データを付加してレスポンスを拡張

### Ubiquitous Language Additions

| Term | Definition | Context |
|------|------------|---------|
| Cloud Cover | 雲量。0-100% で表現し、観測条件の主要指標 | Pass Prediction |
| Weather Condition | 天候分類 (clear/cloudy/overcast/rain)。WMO Weather Code から変換 | Pass Prediction |
| Observation Difficulty | 雲量 >= 70% のパスに付与される「観測困難」ラベル | Pass Prediction |
| Cloud Layer | 3D グローブ上に半透明オーバーレイされる雲分布タイル | Visualization |

## Open Questions

- [NEEDS CLARIFICATION] 雲タイルの提供元 (RainViewer API vs OpenWeatherMap Tile API) のどちらが適切か
- [NEEDS CLARIFICATION] 「観測困難」のしきい値を雲量 70% とするか、ユーザー設定可能にするか
