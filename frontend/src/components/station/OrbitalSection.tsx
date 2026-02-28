"use client";

import type { OrbitalParams } from "@/lib/types";

interface Props {
  params: OrbitalParams | null;
  isLoading: boolean;
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between gap-4">
      <span className="text-gray-500">{label}</span>
      <span className="text-gray-100 tabular-nums">{value}</span>
    </div>
  );
}

export default function OrbitalSection({ params, isLoading }: Props) {
  if (isLoading) {
    return <p className="text-sm text-gray-500">Loading orbital data...</p>;
  }

  if (!params) {
    return <p className="text-sm text-gray-500">Orbital data unavailable</p>;
  }

  return (
    <div className="space-y-1 text-sm">
      <Row label="Inclination" value={`${params.inclination_deg.toFixed(2)} deg`} />
      <Row label="Period" value={`${params.period_min.toFixed(2)} min`} />
      <Row label="Eccentricity" value={params.eccentricity.toFixed(7)} />
      <Row label="Apogee" value={`${params.apogee_km.toFixed(1)} km`} />
      <Row label="Perigee" value={`${params.perigee_km.toFixed(1)} km`} />
    </div>
  );
}
