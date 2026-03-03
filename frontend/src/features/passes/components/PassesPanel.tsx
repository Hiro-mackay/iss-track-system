"use client";

import { useGeolocation } from "@/hooks/useGeolocation";
import LocationInput from "@/components/LocationInput";
import { usePassPredictions } from "../hooks/usePassPredictions";
import PassList from "./PassList";

export default function PassesPanel() {
  const { lat, lon, error: geoError, setLocation } = useGeolocation();
  const { passes, isLoading } = usePassPredictions(lat, lon);

  return (
    <>
      <div>
        <h3 className="mb-2 text-sm font-medium text-gray-400">
          Upcoming Passes
        </h3>
        <PassList passes={passes} loading={isLoading} />
      </div>
      <LocationInput
        lat={lat}
        lon={lon}
        error={geoError}
        setLocation={setLocation}
      />
    </>
  );
}
