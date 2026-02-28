export interface ISSPosition {
  lat: number;
  lon: number;
  altitude_km: number;
  velocity_kmh: number;
  timestamp: string;
}

export interface OrbitPoint {
  lat: number;
  lon: number;
  altitude_km: number;
}

export type Brightness = "bright" | "moderate" | "dim";

export interface PassEvent {
  time: string;
  azimuth: number;
  elevation: number;
}

export interface PassPrediction {
  rise: PassEvent;
  culmination: PassEvent;
  set: PassEvent;
  max_elevation: number;
  is_visible: boolean;
  brightness: Brightness;
}

export interface CrewMember {
  name: string;
  craft: string;
}

export interface OrbitalParams {
  inclination_deg: number;
  period_min: number;
  eccentricity: number;
  apogee_km: number;
  perigee_km: number;
}

export interface ISSStatus {
  launch_year: number;
  orbit_count: number;
  crew_count: number;
  orbital_params: OrbitalParams;
  tle_epoch: string;
}
