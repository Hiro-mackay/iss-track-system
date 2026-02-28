# Plan: Pass Predictions

## Spec Reference

[Spec](../03-pass-predictions/spec.md)

## Architecture Decisions

- akhenakh/sgp4 GeneratePasses で直接パス予測
- Geolocation は Frontend のカスタムフックで管理、Backend は stateless
- SWR で自動更新、observe 地点変更時にキーが変わり再 fetch

## Implementation Approach

Backend に visibility service を追加し passes router で API 公開。Frontend に geolocation 対応の usePassPredictions hook と PassList テーブルを実装。

## Component Breakdown

### Component 1: Pass Prediction Service

- **Purpose**: GeneratePasses でパス予測計算
- **Location**: `backend/passes.go`
- **Changes**: GetPassPredictions(name, line1, line2, lat, lon, days)

### Component 2: Passes Handler

- **Purpose**: パス予測 REST API
- **Location**: `backend/handler_passes.go`
- **Changes**: GET /api/v1/passes, main.go にルート登録

### Component 3: usePassPredictions Hook

- **Purpose**: SWR でパス予測データ取得
- **Location**: `frontend/src/hooks/usePassPredictions.ts`
- **Changes**: useSWR, refreshInterval 300,000ms

### Component 4: PassList

- **Purpose**: パス一覧テーブル UI
- **Location**: `frontend/src/components/info/PassList.tsx`
- **Changes**: rise/culmination/set 時刻、仰角、明るさ表示

## API Changes

| Method | Path | Response |
|--------|------|----------|
| GET | /api/v1/passes?lat&lon&days=3 | PassPrediction[] |

## Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| GeneratePasses 計算精度 | Low | Low | akhenakh/sgp4 は Vallado 参照実装ベース |
| find_events の計算負荷 | Med | Low | days パラメータで範囲制限 |

## Dependencies

- 01-realtime-tracking が完了していること (TLE fetcher, models)

## Out of Scope

- プッシュ通知 (パス接近アラート)
- ISS 以外の衛星の可視予測
