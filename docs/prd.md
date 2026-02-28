# Product Requirements Document

This document defines the product domain, bounded contexts, and business rules.
It is the "what and where" counterpart to `docs/architecture.md` (the "how").

## Domain Overview

ISS (International Space Station) の軌道追跡・可視化ドメイン。TLE (Two-Line Element) データと SGP4 軌道伝播を用いて ISS のリアルタイム位置を算出し、インタラクティブな地図上に表示する。ユーザーの観測地点から ISS の可視パスを予測する機能も提供する。

## Subdomains

| Subdomain | Type | Bounded Context | Notes |
|-----------|------|-----------------|-------|
<!--
Classification guide:
- Core: competitive advantage, build in-house, highest investment
- Supporting: necessary but not differentiating, build or buy
- Generic: commodity, prefer off-the-shelf (auth, payments, email)
-->
| Orbital Computation | Core | Orbit Tracking | SGP4 propagation, position calculation |
| Map Visualization | Core | Visualization | Real-time map rendering, orbit path display |
| Pass Prediction | Core | Pass Prediction | Visibility computation for observer locations |
| TLE Data Management | Supporting | Orbit Tracking | CelesTrak fetching, 2-hour memory cache |
| Geocoding | Generic | Visualization | Nominatim reverse geocoding for location names |
| Station Information | Supporting | Station Info | Crew manifest, orbital parameters, station status |

## Ubiquitous Language

| Term | Definition | Context |
|------|------------|---------|
| TLE (Two-Line Element) | NORAD standard format for satellite orbital parameters | Orbit Tracking |
| SGP4 | Simplified General Perturbations 4; orbit propagation algorithm for near-Earth objects | Orbit Tracking |
| Epoch | The reference time at which TLE orbital elements are valid | Orbit Tracking |
| Propagation | Computing satellite position at a given time from TLE using SGP4 | Orbit Tracking |
| Ground Track | The path on Earth's surface directly below the satellite | Visualization |
| Pass | A single event where ISS rises above and sets below the observer's horizon | Pass Prediction |
| AOS (Acquisition of Signal) | The moment ISS rises above the minimum elevation angle | Pass Prediction |
| LOS (Loss of Signal) | The moment ISS sets below the minimum elevation angle | Pass Prediction |
| Culmination | The point of maximum elevation during a pass | Pass Prediction |
| Visible Pass | A pass where ISS is sunlit and the observer is in twilight or darkness (sun altitude < -6 deg) | Pass Prediction |
| Observer Location | The latitude/longitude of the user's observation point | Pass Prediction |
| Globe View | Cesium ベースの 3D 地球儀レンダリングビュー | Visualization |
| Terminator | 地球の昼夜境界線。太陽位置から計算される大円 | Visualization |
| Camera Follow Mode | カメラが ISS Entity を自動追従するモード | Visualization |
| HUD (Head-Up Display) | 3D シーン上にオーバーレイ表示する情報パネル | Visualization |
| Crew Manifest | ISS に現在搭乗中の宇宙飛行士一覧。Open Notify API から取得 | Station Info |
| Orbit Count | ISS のこれまでの軌道周回数。TLE エポックと周期から概算 | Station Info |
| Orbital Parameters | TLE から算出される軌道特性値 (傾斜角、周期、離心率) | Station Info |
| Station Status | ISS の基本ステータス情報 (打ち上げ年、クルー数、周回数) | Station Info |

<!--
Lifecycle rules:
- New terms: add with definition and owning context before using in specs or code
- Scope: a term's meaning is authoritative only within its bounded context
- Deprecation: mark as [DEPRECATED] with replacement term; remove after migration
-->

## Bounded Contexts

### Orbit Tracking

**Responsibility**: TLE data acquisition, caching, and SGP4 orbital propagation to compute ISS position (latitude, longitude, altitude, velocity).

#### Aggregates

**Satellite**
- Root Entity: Satellite
- Value Objects: TLE, OrbitalPosition (lat, lon, alt, velocity), Epoch
- Invariants: TLE MUST be refreshed within 2 hours of last fetch. TLE MUST be valid (parseable 3LE format from CelesTrak).
- Domain Events: Tracking.TLERefreshed, Tracking.PositionComputed

### Visualization

**Responsibility**: Real-time map rendering, ISS marker display, orbit path drawing, position info panel with geocoded location names.

#### Aggregates

**MapView**
- Root Entity: MapView
- Value Objects: OrbitPath (array of lat/lon points), MarkerPosition, GeocodedLocation (country/region name)
- Invariants: OrbitPath MUST be split at the date line (longitude delta > 180). Nominatim requests MUST NOT exceed 1 req/s.
- Domain Events: Visualization.OrbitPathUpdated

### Pass Prediction

**Responsibility**: Computing visible ISS passes for a given observer location, including rise/culmination/set times and visibility classification.

#### Aggregates

**PassForecast**
- Root Entity: PassForecast
- Value Objects: ObserverLocation (lat, lon), Pass (AOS time, LOS time, culmination time, max elevation), VisibilityClassification (bright/moderate/dim)
- Invariants: ObserverLocation MUST have valid coordinates (-90..90 lat, -180..180 lon). Pass MUST have minimum elevation >= 10 degrees.
- Domain Events: PassPrediction.ForecastComputed

### Station Info

**Responsibility**: ISS のメタ情報 (搭乗クルー、軌道パラメータ、基本ステータス) の取得・提供。外部 API (Open Notify) と TLE データからステーション情報を集約する。

#### Aggregates

**StationInfo**
- Root Entity: StationInfo
- Value Objects: CrewMember (name, craft/mission), OrbitalParameters (inclination, period, eccentricity), StationStatus (launchYear, orbitCount, crewCount)
- Invariants: CrewMember list MUST only include `craft == "ISS"` entries. OrbitalParameters MUST be derived from current cached TLE. Crew data cache TTL is 1 hour.
- Domain Events: StationInfo.CrewUpdated

## Context Map

| Upstream | Downstream | Relationship | Downstream Pattern |
|----------|------------|-------------|-------------------|
| Orbit Tracking | Visualization | Customer-Supplier | Conformist (WebSocket + REST API) |
| Orbit Tracking | Pass Prediction | Customer-Supplier | Conformist (shared TLE data) |
| CelesTrak (External) | Orbit Tracking | Conformist | Anti-Corruption Layer (TLE parser) |
| Nominatim (External) | Visualization | Separate Ways | Anti-Corruption Layer (geocoder client) |
| Orbit Tracking | Station Info | Customer-Supplier | Conformist (shared TLE data for orbital parameters) |
| Open Notify (External) | Station Info | Conformist | Anti-Corruption Layer (crew data fetcher) |

## Domain Events

<!--
Naming convention: {BoundedContext}.{Aggregate}{PastTenseVerb}
Examples: Orders.OrderPlaced, Shipping.ShipmentDispatched

Minimal event schema:
- eventId: unique identifier
- eventType: the event name
- aggregateId: source aggregate ID
- occurredAt: ISO 8601 timestamp
- payload: event-specific data
-->

| Event | Aggregate | Published When | Consumed By |
|-------|-----------|---------------|-------------|
| Tracking.TLERefreshed | Satellite | New TLE fetched from CelesTrak (every 2 hours) | Orbit Tracking (recompute), Pass Prediction (reforecast) |
| Tracking.PositionComputed | Satellite | SGP4 propagation completes (every 1 second) | Visualization (update marker) |
| Visualization.OrbitPathUpdated | MapView | Orbit path recalculated (every 60 seconds) | Visualization (redraw polyline) |
| PassPrediction.ForecastComputed | PassForecast | Observer location set or TLE refreshed | Visualization (display pass list) |
| StationInfo.CrewUpdated | StationInfo | Crew data fetched from Open Notify API (every 1 hour) | Visualization (display crew list) |
