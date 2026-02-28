# Tasks: ISS Information Panel

## References

- [Spec](./spec.md)
- [Plan](./plan.md)

## Task Breakdown

### Phase 1: Backend Models & Config

- [x] Task 1.1: データモデル追加 (S)
  - CrewMember, OrbitalParams, ISSStatus 構造体を models.go に追加
  - Files: backend/models.go

- [x] Task 1.2: Config 拡張 (S)
  - CrewURL, CrewCacheTTL フィールド追加
  - ISS_CREW_URL, ISS_CREW_CACHE_TTL_SECONDS 環境変数パース
  - Files: backend/config.go

### Phase 2: Backend Core (TDD -- write test first, then implement)

- [x] Task 2.1: CrewCache 実装 (M)
  - Test: crew_test.go -- ISS フィルタ、JSON パース、スティルフォールバック
  - Open Notify API から取得、craft == "ISS" フィルタ、1h TTL
  - TLECache パターン踏襲 (RWMutex, Get/Refresh, graceful degradation)
  - ACL: openNotifyResponse / openNotifyPerson 非公開構造体
  - Files: backend/crew.go, backend/crew_test.go
  - Depends on: 1.1, 1.2

- [x] Task 2.2: ISS Status 算出ロジック実装 (M)
  - Test: status_test.go -- 周期計算、apogee/perigee、周回数推定
  - sgp4.ParseTLE → Inclination, MeanMotion, Eccentricity, RevolutionNumber
  - RecoveredSemiMajorAxis() * 6371 で apogee/perigee 算出
  - orbit_count = RevolutionNumber + elapsed * MeanMotion / 1440
  - Files: backend/station.go, backend/station_test.go
  - Depends on: 1.1

### Phase 3: Backend Wiring

- [x] Task 3.1: Crew ハンドラ + ルート登録 (S)
  - HandleCrew(crew *CrewCache) http.HandlerFunc
  - GET /api/v1/iss/crew ルート登録
  - Files: backend/handler_crew.go, backend/main.go
  - Depends on: 2.1

- [x] Task 3.2: Status ハンドラ + ルート登録 (S)
  - HandleStatus(tle *TLECache, crew *CrewCache) http.HandlerFunc
  - GET /api/v1/iss/status ルート登録
  - CrewCache 初期化 + 定期リフレッシュ goroutine を main.go に追加
  - Files: backend/handler_status.go, backend/main.go
  - Depends on: 2.1, 2.2

### Phase 4: Frontend Types & Data Layer

- [x] Task 4.1: 型定義 + API クライアント (S)
  - CrewMember, OrbitalParams, ISSStatus interface 追加
  - fetchCrew(), fetchStatus() 関数追加
  - Files: frontend/src/lib/types.ts, frontend/src/lib/api.ts
  - Depends on: 3.1, 3.2

- [x] Task 4.2: SWR Hooks + フォールバックデータ (S)
  - useCrewData: SWR (1h refresh) + 静的フォールバック
  - useStationStatus: SWR (60s refresh)
  - 静的クルーフォールバックデータ定義
  - Files: frontend/src/hooks/useCrewData.ts, frontend/src/hooks/useStationStatus.ts, frontend/src/lib/fallback-crew.ts
  - Depends on: 4.1

### Phase 5: Frontend Components

- [x] Task 5.1: CrewSection コンポーネント (S)
  - クルー名 + craft の一覧表示
  - loading/error/empty 各状態のハンドリング
  - API 障害時の "(Offline)" ラベル表示
  - Files: frontend/src/components/station/CrewSection.tsx
  - Depends on: 4.2

- [x] Task 5.2: StatusSection コンポーネント (S)
  - 打ち上げ年 (1998)、軌道周回数、搭乗クルー数
  - Row パターン (PositionPanel.tsx 踏襲)
  - Files: frontend/src/components/station/StatusSection.tsx
  - Depends on: 4.2

- [x] Task 5.3: OrbitalSection コンポーネント (S)
  - 傾斜角 (deg)、周期 (min)、離心率、遠地点/近地点 (km)
  - Row パターン + tabular-nums
  - Files: frontend/src/components/station/OrbitalSection.tsx
  - Depends on: 4.2

- [x] Task 5.4: StationSidebar コンポーネント (M)
  - 3 セクションをアコーディオン形式で合成
  - 各セクション独立開閉 (chevron アイコン + CSS transition)
  - コンテナ: bg-black/70 rounded-lg p-4 backdrop-blur-sm
  - Files: frontend/src/components/station/StationSidebar.tsx
  - Depends on: 5.1, 5.2, 5.3

### Phase 6: Page Integration

- [x] Task 6.1: page.tsx に左サイドバー統合 (M)
  - stationOpen state + トグルボタン (既存左上ボタングループに追加)
  - StationSidebar を absolute top-16 left-4 z-10 w-72 で配置
  - 既存右サイドバーとの共存確認
  - Files: frontend/src/app/page.tsx
  - Depends on: 5.4

### Phase 7: Verification & Documentation

- [x] Task 7.1: 全体検証 (S)
  - task check (lint + typecheck + test + build) 通過
  - task dev で手動確認: 左サイドバー開閉、クルー表示、軌道パラメータ妥当性
  - Depends on: 6.1

- [x] Task 7.2: architecture.md 更新 (S)
  - Station Info コンテキストの追加を反映
  - 新規 API エンドポイントの記載
  - Files: docs/architecture.md
  - Depends on: 7.1

## Estimation Guidance

| Size | Description | Typical Duration |
|------|-------------|-----------------|
| S | Single file, isolated change | < 1 hour |
| M | 2-3 files, minor integration | 1-4 hours |
