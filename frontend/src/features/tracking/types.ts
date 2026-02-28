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
