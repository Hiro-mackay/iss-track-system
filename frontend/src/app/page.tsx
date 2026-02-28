"use client";

import { useState } from "react";
import dynamic from "next/dynamic";
import { useISSPosition, useOrbitTrack, PositionPanel, CountryInfo } from "@/features/tracking";
import { usePassPredictions, PassList } from "@/features/passes";
import { StationSidebar } from "@/features/station";
import { useGeolocation } from "@/hooks/useGeolocation";
import LocationInput from "@/components/LocationInput";

const GlobeView = dynamic(() => import("@/features/tracking/components/GlobeView"), {
  ssr: false,
  loading: () => (
    <div className="flex h-screen w-screen items-center justify-center bg-[#0a0a0a]">
      <p className="text-gray-400">Loading globe...</p>
    </div>
  ),
});

export default function Home() {
  const { position, connected } = useISSPosition();
  const orbitPoints = useOrbitTrack();
  const { lat, lon, error: geoError, setLocation } = useGeolocation();
  const { passes, isLoading: passesLoading } = usePassPredictions(lat, lon);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [stationOpen, setStationOpen] = useState(true);
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
          onClick={() => setStationOpen(!stationOpen)}
          className="rounded-lg bg-black/70 px-3 py-2 text-sm text-gray-300 backdrop-blur-sm transition-colors hover:bg-black/90"
        >
          {stationOpen ? "Hide Station" : "Station Info"}
        </button>
      </div>

      <button
        onClick={() => setSidebarOpen(!sidebarOpen)}
        className="absolute top-4 right-4 z-10 rounded-lg bg-black/70 px-3 py-2 text-sm text-gray-300 backdrop-blur-sm transition-colors hover:bg-black/90"
      >
        {sidebarOpen ? "Hide Info" : "Show Info"}
      </button>

      {stationOpen && (
        <aside className="absolute top-16 left-4 z-10 w-72">
          <StationSidebar />
        </aside>
      )}

      {sidebarOpen && (
        <aside className="absolute top-16 right-4 z-10 w-72 space-y-4 rounded-lg bg-black/70 p-4 backdrop-blur-sm transition-all">
          <PositionPanel position={position} />
          {position && (
            <CountryInfo lat={position.lat} lon={position.lon} />
          )}
          <div>
            <h3 className="mb-2 text-sm font-medium text-gray-400">
              Upcoming Passes
            </h3>
            <PassList passes={passes} loading={passesLoading} />
          </div>
          <LocationInput
            lat={lat}
            lon={lon}
            error={geoError}
            setLocation={setLocation}
          />
        </aside>
      )}
    </main>
  );
}
