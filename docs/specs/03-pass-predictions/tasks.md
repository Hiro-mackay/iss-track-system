# Tasks: Pass Predictions

## References

- [Spec](../03-pass-predictions/spec.md)
- [Plan](../03-pass-predictions/plan.md)

## Task Breakdown

### Phase 1: Backend

- [x] Task 1.1: Pass Prediction Service (M)
  - passes.go: GeneratePasses, 明るさ分類

- [x] Task 1.2: Passes Handler (S)
  - Depends on: 1.1
  - handler_passes.go: GET /api/v1/passes, main.go にルート登録

### Phase 2: Frontend

- [x] Task 2.1: usePassPredictions Hook (S)
  - SWR useSWR, refreshInterval 300,000ms
  - {lat, lon, days} パラメータ

- [x] Task 2.2: PassList コンポーネント (M)
  - Depends on: 2.1
  - パス一覧テーブル (rise/culmination/set, 仰角, 明るさ)

- [x] Task 2.3: Geolocation 対応 (S)
  - navigator.geolocation + フォールバック東京
  - page.tsx に統合

### Phase 3: Testing

- [x] Task 3.1: Visibility Tests (M)
  - Depends on: 1.1
  - rise < culmination < set, elevation >= 10.0, brightness 分類

## Estimation Guidance

| Size | Description | Typical Duration |
|------|-------------|-----------------|
| S | Single file, isolated change | < 1 hour |
| M | 2-3 files, minor integration | 1-4 hours |
