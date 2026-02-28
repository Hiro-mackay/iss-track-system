# ADR-002: Go Orbital Mechanics Library Selection

## Status

Accepted

## Date

2026-02-23

## Context

We are planning to migrate the backend from Python (FastAPI + Skyfield) to Go for improved performance, lower memory footprint, and simplified deployment (single binary). This requires selecting a Go library for SGP4 orbital propagation to replace Skyfield.

Skyfield provides high-accuracy satellite tracking through two layers:
1. **SGP4 propagation** -- the core orbit prediction algorithm (shared by all implementations)
2. **Coordinate transformation pipeline** -- IAU 2000A precession-nutation, polar motion, GAST, proper time scales

For our ISS tracker (map visualization, 1-second WebSocket updates), the accuracy analysis shows:

| Error Source | Magnitude | Notes |
|---|---|---|
| TLE data precision at epoch | ~1 km | Dominant error; unavoidable |
| SGP4 propagation drift | 1-3 km/day | Algorithm limitation |
| Nutation omission (basic SGP4) | ~0.3 km | Skyfield corrects this |
| Polar motion omission | ~0.02 km | Skyfield corrects this |

The TLE ~1 km error dominates. Skyfield's coordinate transformation improvements (~0.3 km total) are smaller than the baseline error and invisible on map display. At zoom level 6, 1 pixel ~ 2.4 km -- the entire SGP4 error is sub-pixel.

The most impactful accuracy lever is **TLE freshness** (our 2-hour cache), not the propagation library.

### Go Library Landscape

| Library | Stars | Last Active | Pass Prediction | SDP4 | Notes |
|---|---|---|---|---|---|
| joshuaferrara/go-satellite | 94 | 2022 | No | No | Most adopted, but dormant. 10 open issues |
| akhenakh/sgp4 | 2 | 2025 | Yes | No | Newest, cleanest API, active development |
| morphism/sgp4go | 4 | 2020 | No | No | C transpile, non-idiomatic Go |
| 3rdof2/sgp4 | 2 | 2024 | No | No | Learning project, no accuracy guarantees |
| SharkEzz/go-sgp4 | 0 | 2022 | No | Unknown | CGo/SWIG, adds build complexity |
| hebl/gofa | 14 | 2023 | N/A | N/A | IAU SOFA in Go (coordinate transforms only) |

ISS orbital period is ~92 minutes, well within SGP4's 225-minute limit. SDP4 (deep-space) support is not needed.

### Self-Implementation Assessment

SGP4 is ~2,000-3,000 lines of C in Vallado's reference implementation. The algorithm is well-documented (Spacetrack Report #3). However, self-implementation is not recommended because:

- Subtle floating-point evaluation order differences cause divergence from reference
- Existing Go ports are already verified against Vallado's test vectors (position error < 0.007 mm between implementations)
- Edge case handling (circular orbits, equatorial orbits, singularities) is error-prone
- Maintenance burden with no accuracy benefit

## Decision

Adopt **akhenakh/sgp4** as the primary orbital propagation library.

Rationale:
- **Built-in pass prediction** -- directly needed for Spec 03 (pass predictions), avoids reimplementing AOS/LOS/max elevation calculation
- **Clean, idiomatic Go API** -- TLE parsing, SGP4 propagation, geodetic conversion, look angles in one package
- **Active development** -- created May 2025, still receiving updates
- **ISS-compatible** -- SGP4 handles ISS orbit (92 min period); SDP4 absence is irrelevant
- **Verified accuracy** -- port of Daniel Warner's C++ SGP4 implementation, tested against reference vectors

Mitigation for low adoption risk:
- The library is a straightforward port of a well-known C++ implementation
- SGP4 is a stable, frozen algorithm (no spec changes expected)
- If abandoned, the codebase is small enough to vendor and maintain
- go-satellite serves as a proven fallback

### What We Lose from Skyfield

| Skyfield Feature | Impact on ISS Tracker | Mitigation |
|---|---|---|
| IAU 2000A precession-nutation | ~0.3 km (sub-pixel on map) | Not needed for map visualization |
| Polar motion (IERS EOP) | ~0.02 km (negligible) | Not needed |
| JPL ephemeris (DE421/440) | Sun/Moon positions for eclipse calc | Can add hebl/gofa later if needed |
| Precise time scales (UTC/TAI/TT) | Meters-level difference | Basic UTC handling is sufficient |

None of these losses are user-visible in our application.

## Consequences

### Positive

- Single binary deployment (no Python runtime, no virtualenv)
- Lower memory footprint and faster startup
- Pass prediction built into the library (Spec 03)
- WebSocket handling is native in Go standard library
- Consistent language if other Go services are added

### Negative

- Smaller community than Skyfield (2 stars vs 1,400+)
- No SDP4 -- limits future expansion to LEO satellites only
- Library is new (2025) with limited production track record
- If precise coordinate transforms are ever needed, hebl/gofa must be integrated separately

### Neutral

- SGP4 accuracy is identical across all implementations -- no regression from Skyfield for our use case
- TLE freshness remains the dominant accuracy factor regardless of library choice
- The library's lack of SDP4 is irrelevant unless we track GEO/HEO satellites in the future

## Alternatives Considered

### Alternative 1: joshuaferrara/go-satellite

The most adopted Go satellite library (94 stars). Provides SGP4 propagation, coordinate conversions, and SpaceTrack API integration.

Not chosen because:
- Dormant since 2022, 10 open issues with no maintainer response
- No built-in pass prediction (would need separate implementation for Spec 03)
- Multiple forks suggest users have had to work around bugs
- Remains a viable fallback if akhenakh/sgp4 proves unsuitable

### Alternative 2: Self-Implementation of SGP4

Port Vallado's reference C implementation to idiomatic Go.

Not chosen because:
- Floating-point evaluation order differences make verification difficult
- Existing ports already achieve < 0.007 mm agreement with reference
- Maintenance burden with zero accuracy benefit over existing libraries
- Engineering time better spent on application features

### Alternative 3: CGo Bindings (SharkEzz/go-sgp4 or custom)

Wrap an established C/C++ SGP4 library via CGo/SWIG.

Not chosen because:
- Adds C++ compiler and SWIG as build dependencies
- Breaks cross-compilation simplicity (a key Go advantage)
- CGo overhead for frequent calls (1 Hz propagation)
- Existing pure Go implementations are sufficient

### Alternative 4: Keep Python Backend (No Migration)

Continue with FastAPI + Skyfield.

Not chosen because this ADR assumes the Go migration decision has been made separately. If the migration is reconsidered, Skyfield remains an excellent choice for this domain.
