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
