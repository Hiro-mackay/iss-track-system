import type { CrewMember, ISSStatus } from "./types";
import { API_BASE } from "@/lib/api-client";

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
