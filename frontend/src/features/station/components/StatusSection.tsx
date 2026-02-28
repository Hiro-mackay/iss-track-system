"use client";

import type { ISSStatus } from "../types";

interface Props {
  status: ISSStatus | null;
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

export default function StatusSection({ status, isLoading }: Props) {
  if (isLoading) {
    return <p className="text-sm text-gray-500">Loading status...</p>;
  }

  if (!status) {
    return <p className="text-sm text-gray-500">Status unavailable</p>;
  }

  return (
    <div className="space-y-1 text-sm">
      <Row label="Launch Year" value={String(status.launch_year)} />
      <Row label="Orbit Count" value={status.orbit_count.toLocaleString()} />
      <Row label="Crew Count" value={String(status.crew_count)} />
    </div>
  );
}
