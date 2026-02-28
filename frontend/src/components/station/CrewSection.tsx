"use client";

import type { CrewMember } from "@/lib/types";

interface Props {
  crew: CrewMember[];
  isLoading: boolean;
  isOffline: boolean;
}

export default function CrewSection({ crew, isLoading, isOffline }: Props) {
  if (isLoading) {
    return <p className="text-sm text-gray-500">Loading crew...</p>;
  }

  return (
    <div className="space-y-1 text-sm">
      {isOffline && (
        <p className="text-xs text-yellow-500">(Offline)</p>
      )}
      {crew.map((member) => (
        <div key={member.name} className="flex justify-between gap-4">
          <span className="text-gray-100">{member.name}</span>
          <span className="text-gray-500">{member.craft}</span>
        </div>
      ))}
      {crew.length === 0 && (
        <p className="text-gray-500">No crew data available</p>
      )}
    </div>
  );
}
