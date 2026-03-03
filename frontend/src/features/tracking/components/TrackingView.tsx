"use client";

import { useState } from "react";
import dynamic from "next/dynamic";
import { useISSPosition } from "../hooks/useISSPosition";
import { useOrbitTrack } from "../hooks/useOrbitTrack";
import PositionPanel from "./PositionPanel";
import CountryInfo from "./CountryInfo";

const GlobeView = dynamic(() => import("./GlobeView"), {
  ssr: false,
  loading: () => (
    <div className="flex h-screen w-screen items-center justify-center bg-[#0a0a0a]">
      <p className="text-gray-400">Loading globe...</p>
    </div>
  ),
});

interface Props {
  sidebarOpen: boolean;
  onSidebarToggle: () => void;
  stationOpen: boolean;
  onStationToggle: () => void;
  passesSlot: React.ReactNode;
  stationSlot: React.ReactNode;
}

export default function TrackingView({
  sidebarOpen,
  onSidebarToggle,
  stationOpen,
  onStationToggle,
  passesSlot,
  stationSlot,
}: Props) {
  const { position, connected } = useISSPosition();
  const orbitPoints = useOrbitTrack();
  const [followISS, setFollowISS] = useState(false);

  return (
    <main className="relative h-screen w-screen">
      <GlobeView
        position={position}
        orbitPoints={orbitPoints}
        followISS={followISS}
        onFollowChange={setFollowISS}
      />

      <div className="absolute top-4 left-4 z-10 flex items-center gap-2">
        <div className="rounded-lg bg-black/70 px-3 py-2 text-sm backdrop-blur-sm">
          <div className="flex items-center gap-2">
            <span
              className={`inline-block h-2 w-2 rounded-full ${
                connected ? "bg-green-500" : "bg-red-500"
              }`}
            />
            <span className="text-gray-300">
              {connected ? "Live" : "Reconnecting..."}
            </span>
          </div>
        </div>

        <button
          onClick={() => setFollowISS(!followISS)}
          className="rounded-lg bg-black/70 px-3 py-2 text-sm text-gray-300 backdrop-blur-sm transition-colors hover:bg-black/90"
        >
          {followISS ? "Free Camera" : "Follow ISS"}
        </button>

        <button
          onClick={onStationToggle}
          className="rounded-lg bg-black/70 px-3 py-2 text-sm text-gray-300 backdrop-blur-sm transition-colors hover:bg-black/90"
        >
          {stationOpen ? "Hide Station" : "Station Info"}
        </button>
      </div>

      <button
        onClick={onSidebarToggle}
        className="absolute top-4 right-4 z-10 rounded-lg bg-black/70 px-3 py-2 text-sm text-gray-300 backdrop-blur-sm transition-colors hover:bg-black/90"
      >
        {sidebarOpen ? "Hide Info" : "Show Info"}
      </button>

      {stationOpen && (
        <aside className="absolute top-16 left-4 z-10 w-72">
          {stationSlot}
        </aside>
      )}

      {sidebarOpen && (
        <aside className="absolute top-16 right-4 z-10 w-72 space-y-4 rounded-lg bg-black/70 p-4 backdrop-blur-sm transition-all">
          <PositionPanel position={position} />
          {position && (
            <CountryInfo lat={position.lat} lon={position.lon} />
          )}
          {passesSlot}
        </aside>
      )}
    </main>
  );
}
