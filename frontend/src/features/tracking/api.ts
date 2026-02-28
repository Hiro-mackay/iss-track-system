import type { ISSPosition, OrbitPoint } from "./types";
import { API_BASE } from "@/lib/api-client";

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
