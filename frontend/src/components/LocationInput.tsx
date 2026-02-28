"use client";

import { useState } from "react";

interface Props {
  lat: number;
  lon: number;
  error: string | null;
  setLocation: (lat: number, lon: number) => void;
}

export default function LocationInput({ lat, lon, error, setLocation }: Props) {
  const [latInput, setLatInput] = useState(String(lat));
  const [lonInput, setLonInput] = useState(String(lon));
  const [validationError, setValidationError] = useState<string | null>(null);
  const [prevLat, setPrevLat] = useState(lat);
  const [prevLon, setPrevLon] = useState(lon);

  if (prevLat !== lat || prevLon !== lon) {
    setPrevLat(lat);
    setPrevLon(lon);
    setLatInput(String(lat));
    setLonInput(String(lon));
  }

  if (!error) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const parsedLat = parseFloat(latInput);
    const parsedLon = parseFloat(lonInput);

    if (isNaN(parsedLat) || isNaN(parsedLon)) {
      setValidationError("Enter valid numbers");
      return;
    }
    if (parsedLat < -90 || parsedLat > 90) {
      setValidationError("Latitude must be between -90 and 90");
      return;
    }
    if (parsedLon < -180 || parsedLon > 180) {
      setValidationError("Longitude must be between -180 and 180");
      return;
    }

    setValidationError(null);
    setLocation(parsedLat, parsedLon);
  };

  return (
    <div>
      <h3 className="mb-2 text-sm font-medium text-gray-400">
        Observer Location
      </h3>
      <p className="mb-2 text-xs text-yellow-500">{error} - using default</p>
      <form onSubmit={handleSubmit} className="space-y-2">
        <div className="flex gap-2">
          <label className="flex-1">
            <span className="text-xs text-gray-500">Lat</span>
            <input
              type="number"
              step="any"
              value={latInput}
              onChange={(e) => setLatInput(e.target.value)}
              className="w-full rounded border border-gray-600 bg-gray-800 px-2 py-1 text-xs text-gray-300 focus:border-blue-500 focus:outline-none"
            />
          </label>
          <label className="flex-1">
            <span className="text-xs text-gray-500">Lon</span>
            <input
              type="number"
              step="any"
              value={lonInput}
              onChange={(e) => setLonInput(e.target.value)}
              className="w-full rounded border border-gray-600 bg-gray-800 px-2 py-1 text-xs text-gray-300 focus:border-blue-500 focus:outline-none"
            />
          </label>
        </div>
        {validationError && (
          <p className="text-xs text-red-400">{validationError}</p>
        )}
        <button
          type="submit"
          className="w-full rounded bg-blue-600 px-2 py-1 text-xs text-white transition-colors hover:bg-blue-700"
        >
          Update
        </button>
      </form>
    </div>
  );
}
