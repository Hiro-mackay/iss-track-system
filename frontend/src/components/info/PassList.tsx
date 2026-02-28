"use client";

import type { PassPrediction, Brightness } from "@/lib/types";

interface Props {
  passes: PassPrediction[];
  loading?: boolean;
}

const brightnessColor: Record<Brightness, string> = {
  bright: "text-yellow-400",
  moderate: "text-blue-400",
  dim: "text-gray-500",
};

function formatTime(iso: string): string {
  return new Date(iso).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString([], { month: "short", day: "numeric" });
}

function duration(rise: string, set: string): string {
  const ms = new Date(set).getTime() - new Date(rise).getTime();
  const min = Math.floor(ms / 60000);
  const sec = Math.floor((ms % 60000) / 1000);
  return `${min}m ${sec}s`;
}

export default function PassList({ passes, loading }: Props) {
  if (loading) {
    return <p className="text-sm text-gray-500">Loading passes...</p>;
  }

  if (passes.length === 0) {
    return <p className="text-sm text-gray-500">No visible passes</p>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-xs">
        <thead>
          <tr className="text-gray-500 border-b border-gray-700">
            <th className="py-1 text-left font-medium">Date</th>
            <th className="py-1 text-left font-medium">Rise</th>
            <th className="py-1 text-right font-medium">Max El.</th>
            <th className="py-1 text-left font-medium">Set</th>
            <th className="py-1 text-right font-medium">Duration</th>
            <th className="py-1 text-right font-medium">Brightness</th>
          </tr>
        </thead>
        <tbody>
          {passes.map((p) => (
            <tr key={p.rise.time} className="border-b border-gray-800 text-gray-300">
              <td className="py-1">{formatDate(p.rise.time)}</td>
              <td className="py-1 tabular-nums">{formatTime(p.rise.time)}</td>
              <td className="py-1 text-right tabular-nums">{p.max_elevation.toFixed(1)}&deg;</td>
              <td className="py-1 tabular-nums">{formatTime(p.set.time)}</td>
              <td className="py-1 text-right tabular-nums">{duration(p.rise.time, p.set.time)}</td>
              <td className={`py-1 text-right capitalize ${brightnessColor[p.brightness] ?? "text-gray-500"}`}>
                {p.brightness}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
