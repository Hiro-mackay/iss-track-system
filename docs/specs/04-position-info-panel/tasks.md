# Tasks: Position Info Panel

## References

- [Spec](../04-position-info-panel/spec.md)
- [Plan](../04-position-info-panel/plan.md)

## Task Breakdown

### Phase 1: Components

- [x] Task 1.1: PositionPanel (S)
  - 緯度/経度/高度/速度の数値フォーマット表示
  - ISSPosition を props で受け取り

- [x] Task 1.2: CountryInfo (M)
  - Nominatim Reverse Geocoding fetch
  - デバウンス (5-10 秒間隔)
  - 海上/エラー時のフォールバック表示

### Phase 2: Integration

- [x] Task 2.1: page.tsx にサイドバー統合 (S)
  - Depends on: 1.1, 1.2
  - サイドバー toggle ボタン (useState)
  - PositionPanel + CountryInfo + PassList 配置

## Estimation Guidance

| Size | Description | Typical Duration |
|------|-------------|-----------------|
| S | Single file, isolated change | < 1 hour |
| M | 2-3 files, minor integration | 1-4 hours |
