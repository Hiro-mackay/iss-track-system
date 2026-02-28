# ADR-003: Leaflet 2D から Cesium 3D Globe への移行

## Status

Accepted

## Date

2026-02-24

## Context

ADR-001 で「Leaflet の 2D 地図は 3D globe より視覚的インパクトが弱い」ことを認識しつつ、初期実装の速度を優先して 2D を採用した。Spec 01-04 の実装が完了し基盤が安定した今、ISS の軌道を 3D 地球儀で高度付きに可視化し、昼夜境界やカメラ追従モードを提供する要件が生まれた (Spec 05)。

2D メルカトル図法では以下の問題がある:
- ISS の高度 (~420km) を表現できない
- 軌道の立体的な曲線が直線に見える
- 昼夜境界と衛星日照の関係が分かりにくい
- 日付変更線での Polyline 分割が必要 (3D 球体なら不要)

## Decision

Leaflet (react-leaflet) を **Cesium.js + resium** (React bindings) に完全置換する。

選定理由:
- **3D 地球儀のデファクト標準**: 地理空間 3D ビジュアライゼーションで最も成熟したライブラリ
- **Scene lighting 内蔵**: `globe.enableLighting = true` だけで昼夜境界を描画
- **Entity API**: PointGraphics, PolylineGraphics で高度付きオブジェクト配置が容易
- **Camera tracking**: `viewer.trackedEntity` で ISS 追従モードを 1 行で実現
- **Cesium ion 不要構成**: OSM basemap を使用し、有料トークンへの依存を回避

実装方針:
- 既存 hooks (useISSPosition, useOrbitTrack 等) はそのまま再利用
- コンポーネント層のみ入れ替え: `components/map/` -> `components/globe/`
- サイドバーを HUD オーバーレイに変換: `components/hud/`
- バックエンド API 変更なし

## Consequences

### Positive

- ISS の高度と軌道の立体感を直感的に表現できる
- 昼夜境界が scene lighting で自動描画される
- 日付変更線の Polyline 分割ロジックが不要になる (球体座標系)
- カメラ追従モードで ISS 視点の体験を提供できる
- 既存の hooks/API 層はそのまま再利用可能

### Negative

- Cesium の静的アセットが ~30MB (Workers, Assets) と大きい
- Next.js との統合に webpack 設定 (CopyWebpackPlugin) が必要
- resium の React 19 互換性リスク (未検証)
- Leaflet 比でランタイムメモリ使用量が増加
- 2D 地図の操作性に慣れたユーザーへの学習コスト

### Neutral

- SSR 不可は Leaflet と同様 (dynamic import, ssr: false)
- コンポーネント数は同程度 (4 map -> 5 globe + 1 hud)
- Leaflet 固有の OrbitPath.splitAtDateLine が不要になり、OrbitPath3D は単純化される

## Alternatives Considered

### Alternative 1: Three.js + react-three-fiber

低レベル 3D ライブラリ。バンドルサイズは小さいが、地球儀テクスチャ、地理座標変換、scene lighting を全て自前実装する必要がある。Cesium が提供する GIS 機能 (座標系、Entity API、カメラ制御) を再発明することになり、工数が大幅に増加する。

### Alternative 2: react-globe.gl

Three.js ベースの 3D 地球儀ラッパー。導入は簡単だが、カスタマイズの幅が狭い。高度付き 3D ポリラインやカメラ追従モードの実装が困難。Cesium の Entity API に比べて表現力が不足する。

### Alternative 3: Leaflet を維持し 3D は別ページで提供

既存の 2D 地図を残しつつ、3D ビューを別ルートで提供する。コードの二重管理が発生し、メンテナンスコストが高い。ユーザーの回答「3D に完全置換」に反する。
