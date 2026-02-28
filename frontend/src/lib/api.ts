import type {
  ISSPosition,
  OrbitPoint,
  PassPrediction,
  CrewMember,
  ISSStatus,
  Brightness,
} from "./types";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

const BRIGHTNESS_VALUES: ReadonlySet<string> = new Set<Brightness>([
  "bright",
  "moderate",
  "dim",
]);

export function getWebSocketUrl(): string {
  const base = API_BASE.replace(/^http/, "ws");
  return `${base}/ws/position`;
}

export function isISSPosition(data: unknown): data is ISSPosition {
  return (
    typeof data === "object" &&
    data !== null &&
    typeof (data as ISSPosition).lat === "number" &&
    typeof (data as ISSPosition).lon === "number" &&
    typeof (data as ISSPosition).altitude_km === "number" &&
    typeof (data as ISSPosition).velocity_kmh === "number" &&
    typeof (data as ISSPosition).timestamp === "string"
  );
}

function isOrbitPoint(data: unknown): data is OrbitPoint {
  return (
    typeof data === "object" &&
    data !== null &&
    typeof (data as OrbitPoint).lat === "number" &&
    typeof (data as OrbitPoint).lon === "number" &&
    typeof (data as OrbitPoint).altitude_km === "number"
  );
}

function isPassEvent(data: unknown): boolean {
  return (
    typeof data === "object" &&
    data !== null &&
    typeof (data as { time: unknown }).time === "string" &&
    typeof (data as { azimuth: unknown }).azimuth === "number" &&
    typeof (data as { elevation: unknown }).elevation === "number"
  );
}

function isPassPrediction(data: unknown): data is PassPrediction {
  if (typeof data !== "object" || data === null) return false;
  const p = data as PassPrediction;
  return (
    isPassEvent(p.rise) &&
    isPassEvent(p.culmination) &&
    isPassEvent(p.set) &&
    typeof p.max_elevation === "number" &&
    typeof p.is_visible === "boolean" &&
    BRIGHTNESS_VALUES.has(p.brightness)
  );
}

function isCrewMember(data: unknown): data is CrewMember {
  return (
    typeof data === "object" &&
    data !== null &&
    typeof (data as CrewMember).name === "string" &&
    typeof (data as CrewMember).craft === "string"
  );
}

function isOrbitalParams(data: unknown): boolean {
  if (typeof data !== "object" || data === null) return false;
  const o = data as Record<string, unknown>;
  return (
    typeof o.inclination_deg === "number" &&
    typeof o.period_min === "number" &&
    typeof o.eccentricity === "number" &&
    typeof o.apogee_km === "number" &&
    typeof o.perigee_km === "number"
  );
}

function isISSStatus(data: unknown): data is ISSStatus {
  if (typeof data !== "object" || data === null) return false;
  const s = data as ISSStatus;
  return (
    typeof s.launch_year === "number" &&
    typeof s.orbit_count === "number" &&
    typeof s.crew_count === "number" &&
    isOrbitalParams(s.orbital_params) &&
    typeof s.tle_epoch === "string"
  );
}

export async function fetchPosition(): Promise<ISSPosition> {
  const res = await fetch(`${API_BASE}/api/v1/position`);
  if (!res.ok) throw new Error(`Failed to fetch position: ${res.status}`);
  const data: unknown = await res.json();
  if (!isISSPosition(data)) throw new Error("Invalid position response");
  return data;
}

export async function fetchOrbit(minutes = 90): Promise<OrbitPoint[]> {
  const res = await fetch(`${API_BASE}/api/v1/orbit?minutes=${minutes}`);
  if (!res.ok) throw new Error(`Failed to fetch orbit: ${res.status}`);
  const data: unknown = await res.json();
  if (!Array.isArray(data) || !data.every(isOrbitPoint))
    throw new Error("Invalid orbit response");
  return data;
}

export async function fetchPasses(
  lat: number,
  lon: number,
  days = 3
): Promise<PassPrediction[]> {
  const res = await fetch(
    `${API_BASE}/api/v1/passes?lat=${lat}&lon=${lon}&days=${days}`
  );
  if (!res.ok) throw new Error(`Failed to fetch passes: ${res.status}`);
  const data: unknown = await res.json();
  if (!Array.isArray(data) || !data.every(isPassPrediction))
    throw new Error("Invalid passes response");
  return data;
}

export async function fetchCrew(): Promise<CrewMember[]> {
  const res = await fetch(`${API_BASE}/api/v1/iss/crew`);
  if (!res.ok) throw new Error(`Failed to fetch crew: ${res.status}`);
  const data: unknown = await res.json();
  if (!Array.isArray(data) || !data.every(isCrewMember))
    throw new Error("Invalid crew response");
  return data;
}

export async function fetchStatus(): Promise<ISSStatus> {
  const res = await fetch(`${API_BASE}/api/v1/iss/status`);
  if (!res.ok) throw new Error(`Failed to fetch status: ${res.status}`);
  const data: unknown = await res.json();
  if (!isISSStatus(data)) throw new Error("Invalid status response");
  return data;
}
