# ADR-001: Technology Stack Selection

## Status

Superseded (backend) by [ADR-002](002-go-orbital-library.md). Frontend decisions remain active.

## Date

2026-02-23

## Context

We need to build a real-time ISS tracking system that displays the International Space Station's current position on an interactive map. The system must compute orbital positions from TLE data using SGP4 propagation and stream updates to the browser at 1-second intervals.

Key requirements:
- Real-time position updates via WebSocket
- SGP4 orbital propagation from TLE data
- Interactive map with ISS marker
- TLE data fetching from CelesTrak with caching

## Decision

Adopt a monorepo structure with separate backend and frontend directories:

**Backend: FastAPI + Skyfield/SGP4 (Python 3.12, uv)**
- FastAPI provides native WebSocket support and async capabilities
- Skyfield is the most accurate Python library for satellite position computation
- SGP4 propagation is well-supported in the Python ecosystem
- uv for fast, reliable dependency management

**Frontend: Next.js (TypeScript, App Router, Tailwind CSS) + Leaflet (react-leaflet)**
- Next.js App Router for modern React with server components
- TypeScript for type safety across the frontend
- Tailwind CSS for rapid UI development
- Leaflet via react-leaflet for lightweight, performant map rendering

**Real-time: WebSocket at 1-second intervals**
- WebSocket for low-latency bidirectional communication
- 1-second update interval balances smoothness with resource usage

**Data: CelesTrak TLE with 2-hour memory cache**
- CelesTrak is the standard public source for satellite TLE data
- 2-hour cache avoids excessive API calls while keeping data fresh

## Consequences

### Positive

- Python ecosystem has the best satellite orbit computation libraries
- FastAPI's async support handles WebSocket connections efficiently
- Next.js + TypeScript provides a modern, type-safe frontend
- Leaflet is lightweight compared to alternatives like Mapbox GL
- Monorepo keeps related code together while allowing independent deployment

### Negative

- Two languages (Python + TypeScript) increases cognitive overhead
- No persistent database means TLE cache is lost on restart
- Leaflet's 2D projection may be less visually impressive than 3D globe alternatives

### Neutral

- uv is relatively new but rapidly maturing as a Python package manager
- CelesTrak rate limits are generous for a single-application use case

## Alternatives Considered

### Alternative 1: Full JavaScript Stack (Node.js + satellite.js)

satellite.js provides SGP4 in JavaScript, enabling a single-language stack. Not chosen because Skyfield is significantly more accurate and feature-rich for orbital mechanics, and the Python scientific computing ecosystem is unmatched.

### Alternative 2: Three.js 3D Globe

A 3D globe visualization using Three.js/CesiumJS would be more visually striking. Not chosen for initial implementation due to higher complexity; Leaflet provides a faster path to a working product. Can be added later as an enhancement.
