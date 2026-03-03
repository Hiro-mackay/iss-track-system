"use client";

import { useState } from "react";
import TrackingView from "@/features/tracking/components/TrackingView";
import { PassesPanel } from "@/features/passes";
import { StationSidebar } from "@/features/station";

export default function ISSTrackerClient() {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [stationOpen, setStationOpen] = useState(true);

  return (
    <TrackingView
      sidebarOpen={sidebarOpen}
      onSidebarToggle={() => setSidebarOpen(!sidebarOpen)}
      stationOpen={stationOpen}
      onStationToggle={() => setStationOpen(!stationOpen)}
      passesSlot={<PassesPanel />}
      stationSlot={<StationSidebar />}
    />
  );
}
