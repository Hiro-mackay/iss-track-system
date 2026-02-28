import type { PassPrediction, Brightness } from "./types";
import { API_BASE } from "@/lib/api-client";

const BRIGHTNESS_VALUES: ReadonlySet<string> = new Set<Brightness>([
  "bright",
  "moderate",
  "dim",
]);

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
