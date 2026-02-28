"use client";

import { useReverseGeocode } from "../hooks/useReverseGeocode";

interface Props {
  lat: number;
  lon: number;
}

export default function CountryInfo({ lat, lon }: Props) {
  const { location, loading } = useReverseGeocode(lat, lon);

  if (!location && !loading) return null;

  return (
    <div className="text-sm">
      <span className="text-gray-500">Location</span>
      <p className="mt-0.5 text-gray-100">
        {loading && !location ? "Looking up..." : location}
      </p>
    </div>
  );
}
