# Plan: Position Info Panel

## Spec Reference

[Spec](../04-position-info-panel/spec.md)

## Architecture Decisions

- Nominatim は Frontend から直接呼び出し (バックエンドプロキシ不要)
- レート制限対策はクライアント側デバウンス (5-10 秒間隔)
- サイドバー toggle 状態は useState で管理

## Implementation Approach

PositionPanel と CountryInfo の 2 つのコンポーネントを実装し、page.tsx のサイドバーに配置。toggle ボタンで開閉制御。

## Component Breakdown

### Component 1: PositionPanel

- **Purpose**: 緯度・経度・高度・速度の数値表示
- **Location**: `frontend/src/components/info/PositionPanel.tsx`
- **Changes**: ISSPosition を props で受け取り、フォーマットして表示

### Component 2: CountryInfo

- **Purpose**: ISS 上空の地域名表示
- **Location**: `frontend/src/components/info/CountryInfo.tsx`
- **Changes**: Nominatim fetch + デバウンス、海上フォールバック

### Component 3: Sidebar Toggle

- **Purpose**: サイドバーの表示/非表示
- **Location**: `frontend/src/app/page.tsx` (更新)
- **Changes**: useState toggle, CSS transition

## API Changes

なし (Nominatim は外部 API を直接呼び出し)

## Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| Nominatim レート制限 | Low | Med | デバウンス (5-10s 間隔) |
| Nominatim ダウン | Low | Low | エラー時は地域名を非表示 |

## Dependencies

- 01-realtime-tracking が完了していること (ISSPosition, page.tsx)

## Out of Scope

- 地域名のバックエンドキャッシュ
- 多言語対応の地域名表示
