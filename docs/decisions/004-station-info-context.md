# ADR-004: Station Info Bounded Context and Open Notify API

## Status

Proposed

## Date

2026-02-26

## Context

ISS Information Panel 機能を追加するにあたり、2 つの設計判断が必要になった。

1. **バウンデッドコンテキストの配置**: クルー情報や軌道パラメータは既存の Orbit Tracking (TLE/SGP4 計算) や Visualization (地図描画) とは責務が異なる。既存コンテキストに押し込むとドメインの境界が曖昧になる。

2. **クルーデータの取得元**: ISS に現在搭乗中の宇宙飛行士一覧を取得する外部 API の選定が必要。データは低頻度更新 (Expedition 交代は数ヶ月単位) のため、高可用性よりもシンプルさを優先する。

## Decision

### Station Info バウンデッドコンテキストの新設

ISS のメタ情報 (クルー、軌道パラメータ、基本ステータス) を管理する新しいバウンデッドコンテキスト「Station Info」を導入する。

- Aggregate: StationInfo (CrewMember, OrbitalParameters, StationStatus)
- Upstream: Orbit Tracking (TLE データを参照して軌道パラメータを算出)
- Downstream: Visualization (HUD パネルでの表示)

### Open Notify API の採用

クルーデータの取得に Open Notify API (`http://api.open-notify.org/astros.json`) を採用する。

- 無料、認証不要
- レスポンス形式がシンプル (`{"people": [{"name": "...", "craft": "ISS"}]}`)
- 1 時間メモリキャッシュ + フロントエンド静的データフォールバックで可用性を確保

## Consequences

### Positive

- 既存コンテキスト (Orbit Tracking, Visualization) の責務が肥大化しない
- Station Info は Orbit Tracking の TLE キャッシュを参照するだけで追加の外部通信不要 (軌道パラメータ)
- Open Notify API は認証不要で導入コストが最小
- フォールバック戦略により API ダウン時も最低限の情報を提供可能

### Negative

- Open Notify API はコミュニティメンテナンスで SLA なし (ダウンタイムのリスク)
- 新しいバウンデッドコンテキストの追加はドメインモデルの複雑度を増す
- 静的フォールバックデータは Expedition 交代時に手動更新が必要

### Neutral

- Open Notify API は HTTP のみ (HTTPS 非対応) -- バックエンド経由でアクセスするため問題なし
- クルーデータの更新頻度が低い (数ヶ月単位) ため、キャッシュ戦略は単純で済む

## Alternatives Considered

### Alternative 1: Visualization コンテキストに統合

クルーやステータス情報を Visualization の MapView Aggregate に追加する案。責務が「地図描画」から逸脱するため不採用。MapView が肥大化し、将来のコンテキスト分割が困難になるリスクがある。

### Alternative 2: NASA Open Data API

NASA の公式 API は豊富なデータを提供するが、API キー登録が必要で、レスポンス構造が複雑。クルー一覧という単純なユースケースに対してオーバースペック。将来的に実験データ等を追加する際の候補として保留。

### Alternative 3: 静的データのみ (API 不使用)

外部 API に依存せず、フロントエンドに JSON でクルーデータを埋め込む案。Expedition 交代時の手動更新が運用負荷となるため、API を主系として採用し静的データはフォールバックに限定した。
