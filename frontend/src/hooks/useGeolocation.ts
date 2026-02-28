"use client";

import { useCallback, useEffect, useState } from "react";

interface GeolocationState {
  lat: number;
  lon: number;
  loading: boolean;
  error: string | null;
}

export interface UseGeolocationReturn extends GeolocationState {
  setLocation: (lat: number, lon: number) => void;
}

const DEFAULT_LAT = 35.6762;
const DEFAULT_LON = 139.6503;

const geolocationSupported =
  typeof navigator !== "undefined" && !!navigator.geolocation;

export function useGeolocation(): UseGeolocationReturn {
  const [state, setState] = useState<GeolocationState>({
    lat: DEFAULT_LAT,
    lon: DEFAULT_LON,
    loading: geolocationSupported,
    error: geolocationSupported ? null : "Geolocation not supported",
  });

  useEffect(() => {
    if (!geolocationSupported) return;

    navigator.geolocation.getCurrentPosition(
      (position) => {
        setState({
          lat: position.coords.latitude,
          lon: position.coords.longitude,
          loading: false,
          error: null,
        });
      },
      () => {
        setState({ lat: DEFAULT_LAT, lon: DEFAULT_LON, loading: false, error: "Permission denied" });
      },
      { timeout: 10000 }
    );
  }, []);

  const setLocation = useCallback((lat: number, lon: number) => {
    setState((prev) => ({ ...prev, lat, lon }));
  }, []);

  return { ...state, setLocation };
}
