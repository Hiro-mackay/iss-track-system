"use client";

import { useState } from "react";
import { useCrewData } from "../hooks/useCrewData";
import { useStationStatus } from "../hooks/useStationStatus";
import CrewSection from "./CrewSection";
import StatusSection from "./StatusSection";
import OrbitalSection from "./OrbitalSection";

function SectionToggle({
  title,
  open,
  onToggle,
  children,
}: {
  title: string;
  open: boolean;
  onToggle: () => void;
  children: React.ReactNode;
}) {
  return (
    <div>
      <button
        onClick={onToggle}
        className="flex w-full items-center justify-between text-sm font-medium text-gray-400 transition-colors hover:text-gray-200"
      >
        <span>{title}</span>
        <span
          className={`inline-block transition-transform ${open ? "rotate-90" : ""}`}
        >
          {">"}
        </span>
      </button>
      {open && <div className="mt-2">{children}</div>}
    </div>
  );
}

export default function StationSidebar() {
  const { crew, isLoading: crewLoading, isOffline } = useCrewData();
  const { status, isLoading: statusLoading } = useStationStatus();

  const [crewOpen, setCrewOpen] = useState(true);
  const [statusOpen, setStatusOpen] = useState(true);
  const [orbitalOpen, setOrbitalOpen] = useState(true);

  return (
    <div className="space-y-4 rounded-lg bg-black/70 p-4 backdrop-blur-sm">
      <SectionToggle title="Crew" open={crewOpen} onToggle={() => setCrewOpen(!crewOpen)}>
        <CrewSection crew={crew} isLoading={crewLoading} isOffline={isOffline} />
      </SectionToggle>

      <SectionToggle title="Station Status" open={statusOpen} onToggle={() => setStatusOpen(!statusOpen)}>
        <StatusSection status={status} isLoading={statusLoading} />
      </SectionToggle>

      <SectionToggle title="Orbital Parameters" open={orbitalOpen} onToggle={() => setOrbitalOpen(!orbitalOpen)}>
        <OrbitalSection params={status?.orbital_params ?? null} isLoading={statusLoading} />
      </SectionToggle>
    </div>
  );
}
