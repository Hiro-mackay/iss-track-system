# Spec: ISS Information Panel

## Status

Done

## Bounded Context

Station Info (as defined in `docs/prd.md`, see [ADR-004](../../decisions/004-station-info-context.md))

## Overview

ISS に関する付加情報 -- 現在の搭乗クルー、軌道パラメータ、基本ステータス -- を左サイドバーの HUD パネルとして 3D グローブ UI に統合する。位置・軌道データだけでは伝わらない「ISS で今何が起きているか」をユーザーに提供し、トラッキング体験を豊かにする。

## Problem Statement

現在のシステムは ISS の「どこにいるか」は正確に示すが、「そこで何が起きているか」は一切表示しない。ISS には常時 6-7 名のクルーが滞在し数百の実験が進行しているが、ユーザーはこの情報を得るために別サイトを参照する必要がある。ISS のリアルタイム位置と合わせてクルー・ステータス情報を表示することで、トラッキング体験に文脈と深みを加える。

## User Stories

### Story 1: 搭乗クルーの確認

As a 宇宙ファン, I want to ISS に今誰が乗っているか知る, so that 「あの光点に人が乗っている」実感を得られる.

#### Acceptance Criteria

- **Given** ユーザーが ISS 情報パネルを開いた
  **When** クルーデータが取得される
  **Then** ISS 搭乗中のクルー名とミッション名が一覧表示される

- **Given** Open Notify API がダウンしている
  **When** ISS 情報パネルを確認する
  **Then** フロントエンドの静的フォールバックデータでクルー情報が表示される

- **Given** Open Notify API が復旧した
  **When** 次のキャッシュ更新 (1 時間後) が発生する
  **Then** 自動的に最新データに切り替わる

### Story 2: ISS ステータスの確認

As a ユーザー, I want to ISS の基本ステータスと軌道パラメータを見る, so that ミッション全体の概要と軌道特性を把握できる.

#### Acceptance Criteria

- **Given** ISS 情報パネルが表示されている
  **When** ステータスセクションを確認する
  **Then** 打ち上げ年、軌道周回数 (概算)、搭乗クルー数が表示される

- **Given** ISS 情報パネルが表示されている
  **When** 軌道パラメータセクションを確認する
  **Then** 軌道傾斜角 (deg)、周期 (min)、離心率が表示される

### Story 3: パネルの開閉

As a ユーザー, I want to ISS 情報パネルを開閉できる, so that メインの 3D ビューを邪魔せず必要な時だけ情報を見られる.

#### Acceptance Criteria

- **Given** 3D グローブが表示されている
  **When** 左サイドバーの開閉ボタンを押す
  **Then** ISS 情報パネルが表示/非表示に切り替わる

- **Given** ISS 情報パネルが非表示
  **When** ページを確認する
  **Then** 3D グローブの表示領域が最大化されている

## Functional Requirements

- [x] FR-1: バックエンドに `GET /api/v1/iss/crew` -- 現在の ISS 搭乗クルー一覧を返す
- [x] FR-2: Open Notify API (`http://api.open-notify.org/astros.json`) からクルーデータを取得し、`craft == "ISS"` でフィルタリング
- [x] FR-3: クルーデータを 1 時間メモリキャッシュ (`CrewCache` -- TLECache と同様のパターン)
- [x] FR-4: バックエンドに `GET /api/v1/iss/status` -- TLE から算出した軌道パラメータを返す (傾斜角、周期、離心率、軌道周回数)
- [x] FR-5: フロントエンドに ISS 情報パネルコンポーネント (左サイドバー、アコーディオン形式)
- [x] FR-6: クルーセクション -- 名前、ミッション名を一覧表示
- [x] FR-7: ステータスセクション -- 打ち上げ年 (静的)、軌道周回数 (概算)、搭乗クルー数
- [x] FR-8: 軌道パラメータセクション -- 傾斜角 (deg)、周期 (min)、離心率
- [x] FR-9: Open Notify API 障害時のフォールバック -- フロントエンド静的クルー JSON を表示し「(Offline data)」ラベルを付与
- [x] FR-10: 左サイドバーの開閉トグルボタン (デフォルト: 閉)

## Non-Functional Requirements

- [x] NFR-1: Performance -- クルー API は 1 時間キャッシュ。リアルタイム性は不要
- [x] NFR-2: Resilience -- Open Notify API 障害時はフロントエンド静的データにフォールバック。各セクション (クルー/ステータス/軌道) は独立して動作し、1 つの障害でパネル全体が壊れない
- [x] NFR-3: Performance -- 軌道パラメータ API は TLE キャッシュを参照するのみ。追加の外部リクエストなし
- [x] NFR-4: Consistency -- 全ソースファイル 300 行以内

## Boundaries

### Always Do

- 各セクション (クルー / ステータス / 軌道パラメータ) は独立してフォールバック
- パネルはデフォルト折りたたみ (メインの 3D ビューを邪魔しない)
- ダークテーマを維持 (`bg-black/70 backdrop-blur-sm`)
- 既存の右サイドバー (PositionPanel, PassList) には影響しない

### Ask First

- クルーの顔写真表示 (NASA API からの画像取得)
- ISS 実験データベースとの連携 (NASA OSDR 等)
- ライブストリーム埋め込み (現在はスコープ外)

### Never Do

- クルーの個人情報 (SNS アカウント等) の表示
- Open Notify API への 1 分間隔未満のリクエスト
- 静的フォールバックデータの自動更新 (手動メンテナンス)

## Technical Constraints

- Open Notify API: `http://api.open-notify.org/astros.json` -- 無料、認証不要
  - レスポンス形式: `{"number": 10, "people": [{"name": "...", "craft": "ISS"}, ...], "message": "success"}`
  - `craft == "ISS"` でフィルタし、ISS 搭乗者のみ返す
- 軌道パラメータは TLE Line 2 から直接パース可能:
  - 傾斜角 (Inclination): col 9-16
  - 離心率 (Eccentricity): col 27-33 (先頭に `0.` を付加)
  - 平均運動 (Mean Motion): col 53-63 → 周期 = 1440 / mean_motion (min)
- 軌道周回数 = TLE の Revolution Number at Epoch (col 64-68) + (現在時刻 - TLE エポック) / 周期
- ISS 基本情報 (打ち上げ年: 1998、初乗組員: 2000) はフロントエンド静的データとして保持
- 静的フォールバッククルーデータは `frontend/src/lib/fallback-crew.ts` に JSON で管理。Expedition 交代時に手動更新・デプロイ
- 左サイドバーは右サイドバーと同様の z-index・ポジショニングパターン (`absolute top-16 left-4 z-10`)

## Domain Impact

### Affected Aggregates

- New Aggregate: **StationInfo** (Root Entity: StationInfo)
  - Value Objects: CrewMember (name, craft/mission), OrbitalParameters (inclination, period, eccentricity), StationStatus (launchYear, orbitCount, crewCount)
  - Invariants: CrewMember list MUST only include `craft == "ISS"` entries. OrbitalParameters MUST be derived from current cached TLE. CrewCache TTL is 1 hour.

### Domain Events

- Published: StationInfo.CrewUpdated -- クルーデータが Open Notify API から更新されたとき
- Consumed: Tracking.TLERefreshed (from Orbit Tracking) -- 軌道パラメータの再計算トリガー

### Ubiquitous Language Additions

| Term | Definition | Context |
|------|------------|---------|
| Crew Manifest | ISS に現在搭乗中の宇宙飛行士一覧。Open Notify API から取得、`craft == "ISS"` でフィルタ | Station Info |
| Orbit Count | ISS のこれまでの軌道周回数。TLE の Revolution Number + エポックからの経過周回数で算出 | Station Info |
| Orbital Parameters | TLE から算出される軌道特性値 (傾斜角、周期、離心率) | Station Info |
| Station Status | ISS の基本ステータス情報 (打ち上げ年、クルー数、周回数) | Station Info |

## Open Questions

None -- all questions resolved.
