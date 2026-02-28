"use client";

import type { ISSPosition } from "@/lib/types";

interface Props {
  position: ISSPosition | null;
}

function formatLat(lat: number): string {
  return `${Math.abs(lat).toFixed(4)}° ${lat >= 0 ? "N" : "S"}`;
}

function formatLon(lon: number): string {
  return `${Math.abs(lon).toFixed(4)}° ${lon >= 0 ? "E" : "W"}`;
}

function formatTimestamp(iso: string): string {
  return new Date(iso).toLocaleTimeString();
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between gap-4">
      <span className="text-gray-500">{label}</span>
      <span className="text-gray-100 tabular-nums">{value}</span>
    </div>
  );
}

export default function PositionPanel({ position }: Props) {
  if (!position) {
    return (
      <div className="space-y-1 text-sm">
        <p className="text-gray-500">Waiting for data...</p>
      </div>
    );
  }

  return (
    <div className="space-y-1 text-sm">
      <Row label="Latitude" value={formatLat(position.lat)} />
      <Row label="Longitude" value={formatLon(position.lon)} />
      <Row label="Altitude" value={`${position.altitude_km.toFixed(1)} km`} />
      <Row
        label="Velocity"
        value={`${position.velocity_kmh.toFixed(0)} km/h`}
      />
      <Row label="Time" value={formatTimestamp(position.timestamp)} />
    </div>
  );
}
